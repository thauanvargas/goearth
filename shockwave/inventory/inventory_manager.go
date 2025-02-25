package inventory

import (
	"context"
	"fmt"
	"sync"
	"time"

	g "xabbo.b7c.io/goearth"
	"xabbo.b7c.io/goearth/internal/debug"
	"xabbo.b7c.io/goearth/shockwave/in"
	"xabbo.b7c.io/goearth/shockwave/out"
)

var dbg = debug.NewLogger("[inventory]")

var ErrScanSuccess = fmt.Errorf("scan completed successfully")

// Manager tracks the state of the inventory.
type Manager struct {
	ix          g.Interceptor
	updated     g.VoidEvent
	itemRemoved g.Event[ItemArgs]

	scanCtx   context.Context
	scanDone  func(error)
	scanPage  int
	scanItems map[int]struct{}
	scanCh    chan []Item

	mtx   *sync.RWMutex
	items map[int]Item
}

// NewManager creates a new inventory Manager using the provided extension.
func NewManager(ix g.Interceptor) *Manager {
	mgr := &Manager{
		ix:    ix,
		mtx:   &sync.RWMutex{},
		items: map[int]Item{},
	}
	ix.Intercept(out.GETSTRIP).With(mgr.handleGetStrip)
	ix.Intercept(in.STRIPINFO_2).With(mgr.handleStripInfo2)
	ix.Intercept(in.REMOVESTRIPITEM).With(mgr.handleRemoveStripItem)
	return mgr
}

// Item gets the item with the specified ID.
func (mgr *Manager) Item(id int) *Item {
	mgr.mtx.RLock()
	defer mgr.mtx.RUnlock()
	if item, ok := mgr.items[id]; ok {
		return &item
	} else {
		return nil
	}
}

// Items iterates over all inventory items.
func (mgr *Manager) Items(yield func(item Item) bool) {
	mgr.mtx.RLock()
	for _, item := range mgr.items {
		mgr.mtx.RUnlock()
		if !yield(item) {
			return
		}
		mgr.mtx.RLock()
	}
	mgr.mtx.RUnlock()
}

// ItemCount returns the number of items in the inventory.
func (mgr *Manager) ItemCount() int {
	return len(mgr.items)
}

// Scan performs a full load of the inventory by requesting each inventory page.
// The returned context is canceled once the scan has finished successfully or unsuccessfully.
// Calling context.Cause on the scan context will return ErrScanSuccess if the scan completed successfully.
// Otherwise it will return context.DeadlineExceeded if the operation timed out,
// or context.Canceled if it was explicitly canceled.
// Multiple calls to Scan while a scan is in progress will return the same context.
func (mgr *Manager) Scan() context.Context {
	mgr.mtx.Lock()
	defer mgr.mtx.Unlock()

	if mgr.scanCtx != nil {
		return mgr.scanCtx
	}

	dbg.Printf("beginning scan")

	mgr.scanPage = 0
	mgr.scanItems = map[int]struct{}{}
	mgr.scanCtx, mgr.scanDone = context.WithCancelCause(mgr.ix.Context())
	mgr.scanCh = make(chan []Item)

	go mgr.performScan()

	return mgr.scanCtx
}

func (mgr *Manager) performScan() {
	defer func() {
		mgr.mtx.Lock()
		defer mgr.mtx.Unlock()
		mgr.scanDone(context.DeadlineExceeded)
		mgr.scanCtx = nil
	}()

	attempt := 1
	mgr.ix.Send(out.GETSTRIP, []byte("new"))
scan:
	for {
		select {
		case items := <-mgr.scanCh:
			mgr.scanPage++
			var last, wrapped bool
			if len(items) < 9 {
				last = true
			} else {
				for _, item := range items {
					if _, wrapped = mgr.scanItems[item.ItemId]; wrapped {
						break
					}
					mgr.scanItems[item.ItemId] = struct{}{}
				}
			}
			if !wrapped {
				dbg.Printf("scanned page %d (%d items)", mgr.scanPage, len(items))
			}
			if last || wrapped {
				dbg.Printf("completing scan")
				mgr.scanDone(ErrScanSuccess)
			} else {
				// continue scan
				select {
				case <-time.After(550 * time.Millisecond):
					dbg.Printf("continuing scan")
					mgr.ix.Send(out.GETSTRIP, []byte("next"))
				case <-mgr.scanCtx.Done():
					break scan
				}
			}
		case <-time.After(time.Second):
			// timed out
			if attempt < 3 {
				attempt++
				// retry scan
				dbg.Printf("timed out, retrying (attempt %d)", attempt)
				mgr.ix.Send(out.GETSTRIP, []byte("next"))
			} else {
				dbg.Printf("timed out, aborting (attempt %d)", attempt)
				break scan
			}
		case <-mgr.scanCtx.Done():
			// canceled
			break scan
		}
	}
}

func (mgr *Manager) CancelScan() bool {
	mgr.mtx.Lock()
	defer mgr.mtx.Unlock()

	if mgr.scanCtx == nil {
		return false
	}

	dbg.Printf("cancelling scan")

	mgr.scanDone(context.Canceled)
	mgr.scanCtx = nil
	return true
}

func (mgr *Manager) loadItems(items []Item) (added []Item) {
	mgr.mtx.Lock()
	defer mgr.mtx.Unlock()

	n := 0
	for _, item := range items {
		if _, exists := mgr.items[item.ItemId]; !exists {
			n++
			added = append(added, item)
		}
		mgr.items[item.ItemId] = item
	}

	if n > 0 {
		dbg.Printf("added %d item(s)", n)
	}
	return added
}

func (mgr *Manager) removeItem(id int) (item Item, ok bool) {
	mgr.mtx.Lock()
	defer mgr.mtx.Unlock()

	if item, ok = mgr.items[id]; ok {
		delete(mgr.items, id)
		dbg.Printf("removed item (ID: %d)", id)
	} else {
		dbg.Printf("failed to find item to remove (ID: %d)", id)
	}

	return
}

// handlers

func (mgr *Manager) handleGetStrip(e *g.Intercept) {
	mgr.mtx.RLock()
	defer mgr.mtx.RUnlock()
	if mgr.scanCtx != nil {
		dbg.Printf("blocking getstrip (scan in progress)")
		e.Block()
	}
}

func (mgr *Manager) handleStripInfo2(e *g.Intercept) {
	var inv Inventory
	e.Packet.Read(&inv)

	mgr.loadItems(inv.Items)
	mgr.updated.Dispatch()

	mgr.mtx.Lock()
	defer mgr.mtx.Unlock()

	if mgr.scanCtx != nil {
		e.Block()
		select {
		case mgr.scanCh <- inv.Items:
		default:
			dbg.Println("WARNING: failed to send items on inventory scan channel")
		}
	}
}

func (mgr *Manager) handleRemoveStripItem(e *g.Intercept) {
	itemId := e.Packet.ReadInt()
	if item, ok := mgr.removeItem(itemId); ok {
		mgr.itemRemoved.Dispatch(ItemArgs{item})
	}
}
