package trade

import (
	g "xabbo.b7c.io/goearth"
)

// Offers is an array that holds the offers of the trader and tradee, respectively.
type Offers [2]Offer

// Trader returns the offer of the user who initiated the trade.
func (offers Offers) Trader() Offer {
	return offers[0]
}

// Tradee returns the offer of the user who received the trade request.
func (offers Offers) Tradee() Offer {
	return offers[1]
}

// Offer represents a user's offer in a trade.
type Offer struct {
	Accepted  bool
	UserId    int
	ItemCount int
	Items     []TradeItem
}

type TradeItem struct {
	ItemId     int
	Type       string
	Id         int
	ClassId    int
	Class      string
	Colors     string
	DimX, DimY int
	Category   string
	Groupable  int
	Data       string
	Day        int
	Month      int
	Year       int
	SongId     int
	SongName   string
	SongDesc   string
	PosterName string
}

func (item *TradeItem) Parse(p *g.Packet) {
	*item = TradeItem{}
	item.ItemId = p.ReadInt()
	item.Type = p.ReadString()
	item.Id = p.ReadInt()
	item.ClassId = p.ReadInt()
	item.Class = p.ReadString()
	if item.Type == "s" {
		item.Colors = p.ReadString()
		item.DimX = p.ReadInt()
		item.DimY = p.ReadInt()
	} else {
		item.DimX = 1
		item.DimY = 1
	}
	item.Category = p.ReadString()
	item.Groupable = p.ReadInt()
	item.Data = p.ReadString()
	item.Day = p.ReadInt()
	item.Month = p.ReadInt()
	item.Year = p.ReadInt()
	if item.Type == "s" {
		item.SongId = p.ReadInt()
		item.SongName = "furni_" + item.Class + "_name"
		item.SongDesc = "furni_" + item.Class + "_desc"
	} else {
		if item.Class == "poster" {
			item.PosterName = "poster_" + item.Data + "_name"
		} else {
			item.PosterName = "wallitem_" + item.Class + "_name"
		}
	}
}
