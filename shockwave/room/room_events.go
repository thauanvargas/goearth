package room

import g "xabbo.b7c.io/goearth"

// Args hold the arguments for room events.
type Args struct {
	Id   int
	Info *Info
}

// ObjectArgs holds the arguments for floor item events involving a single item.
type ObjectArgs struct {
	Object Object
}

// ObjectUpdateArgs holds the arguments for floor item update events.
type ObjectUpdateArgs struct {
	Pre    Object // Prev is the previous state of the object before the update.
	Object Object // Cur is the current state of the object after the update.
}

// ObjectArgs holds the arguments for floor item events involving a list of items.
type ObjectsArgs struct {
	Objects []Object
}

// SlideArgs holds the arguments for floor item and entity slide events.
type SlideArgs struct {
	From, To     Point
	ObjectSlides []SlideObjectArgs
	// Source contains the object that caused the slide, if it is available.
	Source        *Object
	SlideMoveType SlideMoveType
	// EntitySlide contains arguments for an entity slide event, if an entity is involved in this event.
	EntitySlide *SlideEntityArgs
}

// SlideObjectArgs holds the arguments for floor item slide events.
type SlideObjectArgs struct {
	// Object contains the state of the object after the slide update.
	Object   Object
	From, To Tile
}

type SlideEntityArgs struct {
	// Entity contains the state of the entity after the slide update.
	Entity   Entity
	From, To Tile
}

// ItemArgs holds the arguments for wall item events involing a single item.
type ItemArgs struct {
	Item Item
}

// ItemUpdateArgs holds the arguments for wall item update events.
type ItemUpdateArgs struct {
	Pre  Item // Pre is the previous state of the item before the update.
	Item Item // Item is the current state of the item after the update.
}

// ItemsArgs holds the arguments for wall item events involing a list of items.
type ItemsArgs struct {
	Items []Item
}

// EntityArgs holds the arguments for events involving a single entity.
type EntityArgs struct {
	Entity Entity
}

// EntityUpdateArgs holds the arguments for entity update events.
type EntityUpdateArgs struct {
	Pre    Entity // Prev is the previous state of the entity before the update.
	Entity Entity // Cur is the current state of the entity after the update.
}

// EntitiesArgs holds the arguments for events involving a list of entities.
type EntitiesArgs struct {
	Entered  bool
	Entities []Entity
}

// EntityChat holds the arguments for chat events.
type EntityChatArgs struct {
	EntityArgs
	Type    ChatType
	Message string
}

// Entered registers an event handler that is invoked when the user enters a room.
func (mgr *Manager) Entered(handler g.EventHandler[Args]) {
	mgr.entered.Register(handler)
}

func (mgr *Manager) RightsUpdated(handler g.VoidHandler) {
	mgr.rightsUpdated.Register(handler)
}

// ObjectsLoaded registers an event handler that is invoked when floor items are loaded.
func (mgr *Manager) ObjectsLoaded(handler g.EventHandler[ObjectsArgs]) {
	mgr.objectsLoaded.Register(handler)
}

// ObjectAdded registers an event handler that is invoked when a floor item is added to the room.
func (mgr *Manager) ObjectAdded(handler g.EventHandler[ObjectArgs]) {
	mgr.objectAdded.Register(handler)
}

// ObjectUpdated registers an event handler that is invoked when a floor item is updated in the room.
func (mgr *Manager) ObjectUpdated(handler g.EventHandler[ObjectUpdateArgs]) {
	mgr.objectUpdated.Register(handler)
}

// ObjectRemoved registers an event handler that is invoked when a floor item is removed from the room.
func (mgr *Manager) ObjectRemoved(handler g.EventHandler[ObjectArgs]) {
	mgr.objectRemoved.Register(handler)
}

// Slide registers an event handler that is invoked when floor items or an entity slides, e.g. along a roller.
func (mgr *Manager) Slide(handler g.EventHandler[SlideArgs]) {
	mgr.slide.Register(handler)
}

// ItemsLoaded registers an event handler that is invoked when wall items are loaded.
func (mgr *Manager) ItemsLoaded(handler g.EventHandler[ItemsArgs]) {
	mgr.itemsLoaded.Register(handler)
}

// ItemAdded registers an event handler that is invoked when a wall item is added to the room.
func (mgr *Manager) ItemAdded(handler g.EventHandler[ItemArgs]) {
	mgr.itemAdded.Register(handler)
}

// ItemUpdated registers an event handler that is invoked when a wall item is updated in the room.
func (mgr *Manager) ItemUpdated(handler g.EventHandler[ItemUpdateArgs]) {
	mgr.itemUpdated.Register(handler)
}

// ItemRemoved registers an event handler that is invoked when an item is removed from the room.
func (mgr *Manager) ItemRemoved(handler g.EventHandler[ItemArgs]) {
	mgr.itemRemoved.Register(handler)
}

// EntitiesAdded registers an event handler that is invoked when entities are loaded or enter the room.
// The Entered flag on the EntitiesArgs indicates whether the entity entered the room.
// If not, the entities were already in the room and are being loaded.
func (mgr *Manager) EntitiesAdded(handler g.EventHandler[EntitiesArgs]) {
	mgr.entitiesAdded.Register(handler)
}

func (mgr *Manager) EntityUpdated(handler g.EventHandler[EntityUpdateArgs]) {
	mgr.entityUpdated.Register(handler)
}

// EntityChat registers an event handler that is invoked when an entity sends a chat message.
func (mgr *Manager) EntityChat(handler g.EventHandler[EntityChatArgs]) {
	mgr.entityChat.Register(handler)
}

// EntityLeft registers an event handler that is invoked when an entity leaves the room.
func (mgr *Manager) EntityLeft(handler g.EventHandler[EntityArgs]) {
	mgr.entityLeft.Register(handler)
}

// Left registers an event handler that is invoked when the user leaves the room.
func (mgr *Manager) Left(handler g.EventHandler[Args]) {
	mgr.left.Register(handler)
}
