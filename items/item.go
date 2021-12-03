package item

import (
	"errors"
	"github.com/michaelvmata/path/modifiers"
)

const (
	Armor  = "Armor"
	Weapon = "Weapon"
	Tablet = "Tablet"
)

type Item struct {
	Name      string
	Type      string
	Slot      string
	Modifiers []modifiers.Modifier
}

func (i *Item) AddModifier(modifierType modifiers.Type, value int) {
	modifier := modifiers.Modifier{
		Type:  modifierType,
		Value: value,
	}
	i.Modifiers = append(i.Modifiers, modifier)
}

func NewItem(name string, itemType string, slot string) *Item {
	return &Item{
		Name:      name,
		Type:      itemType,
		Slot:      slot,
		Modifiers: make([]modifiers.Modifier, 0),
	}
}

const (
	Empty    = ""
	Head     = "Head"
	Neck     = "Neck"
	Body     = "Body"
	Arms     = "Arms"
	Hands    = "Hands"
	Waist    = "Waist"
	Legs     = "Legs"
	Feet     = "Feet"
	Wrist    = "Wrist"
	Fingers  = "Fingers"
	MainHand = "MainHand"
)

type Gear struct {
	Head     *Item
	Neck     *Item
	Body     *Item
	Arms     *Item
	Hands    *Item
	Waist    *Item
	Legs     *Item
	Feet     *Item
	Wrist    *Item
	Fingers  *Item
	OffHand  *Item
	MainHand *Item
}

func NewGear() *Gear {
	return &Gear{}
}

func (g *Gear) Equip(item *Item) (*Item, error) {
	if item.Type != Weapon && item.Type != Armor {
		return nil, errors.New("can't equip item")
	}
	var previous *Item
	switch item.Slot {
	case Head:
		previous = g.Head
		g.Head = item
	case Neck:
		previous = g.Neck
		g.Neck = item
	case Body:
		previous = g.Body
		g.Body = item
	case Arms:
		previous = g.Arms
		g.Arms = item
	case Waist:
		previous = g.Waist
		g.Waist = item
	case Legs:
		previous = g.Legs
		g.Legs = item
	case Feet:
		previous = g.Feet
		g.Feet = item
	case Wrist:
		previous = g.Wrist
		g.Wrist = item
	case Fingers:
		previous = g.Fingers
		g.Fingers = item
	case MainHand:
		previous = g.MainHand
		g.MainHand = item
	default:
		return nil, errors.New("bad item slot")
	}
	return previous, nil
}
