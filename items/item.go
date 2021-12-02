package item

import (
	"errors"
	"github.com/michaelvmata/path/modifiers"
)

type Type int

const (
	Armor Type = iota
	Weapon
	Tablet
)

type Item struct {
	Name      string
	Type      Type
	Slot      Slot
	Modifiers []modifiers.Modifier
}

func (i *Item) AddModifier(modifierType modifiers.Type, value int) {
	modifier := modifiers.Modifier{
		Type:  modifierType,
		Value: value,
	}
	i.Modifiers = append(i.Modifiers, modifier)
}

func NewItem(name string, itemType Type, slot Slot) *Item {
	return &Item{
		Name:      name,
		Type:      itemType,
		Slot:      slot,
		Modifiers: make([]modifiers.Modifier, 0),
	}
}

type Slot int

const (
	Empty Slot = iota
	Head
	Neck
	Body
	Arms
	Hands
	Waist
	Legs
	Feet
	Wrist
	Fingers
	OffHand
	MainHand
)

type Gear []*Item

func NewGear() Gear {
	e := make([]*Item, MainHand-1)
	return e
}

func (g Gear) Equip(item *Item) (*Item, error) {
	var previous *Item
	if item.Type == Armor {
		if g[item.Slot] != nil {
			previous = g[item.Slot]
		}
		g[item.Slot] = item
	} else if item.Type == Weapon {
		if g[MainHand] != nil {
			previous = g[MainHand]
		}
		g[MainHand] = item
	} else {
		return nil, errors.New("can't equip item")
	}
	return previous, nil
}
