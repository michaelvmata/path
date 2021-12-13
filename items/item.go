package item

import (
	"errors"
	"github.com/michaelvmata/path/modifiers"
)

type Item struct {
	UUID      string
	Name      string
	Keywords  []string
	Slot      string
	Modifiers []modifiers.Modifier
}

func (i *Item) AddModifier(modifierType string, value int) {
	modifier := modifiers.Modifier{
		Type:  modifierType,
		Value: value,
	}
	i.Modifiers = append(i.Modifiers, modifier)
}

func (i Item) HasKeyword(keyword string) bool {
	for _, candidate := range i.Keywords {
		if candidate == keyword {
			return true
		}
	}
	return false
}

func NewArmor(UUID string, name string, slot string, keywords []string) *Item {
	return &Item{
		UUID:      UUID,
		Name:      name,
		Slot:      slot,
		Keywords:  keywords,
		Modifiers: make([]modifiers.Modifier, 0),
	}
}

type Container struct {
	Capacity int
	Items    []*Item
}

func NewContainer(capacity int) Container {
	return Container{
		Capacity: capacity,
		Items:    make([]*Item, 0),
	}
}

func (c *Container) AddItem(item *Item) error {
	if len(c.Items) == c.Capacity {
		return errors.New("container full")
	}
	c.Items = append(c.Items, item)
	return nil
}

func (c Container) IndexOfItem(keyword string) int {
	for i, item := range c.Items {
		if item.HasKeyword(keyword) {
			return i
		}
	}
	return -1
}

func (c *Container) RemItemAtIndex(index int) *Item {
	item := c.Items[index]
	c.Items = append(c.Items[:index], c.Items[index+1:]...)
	return item
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

func (g *Gear) Remove(keyword string) *Item {
	if g.Head != nil && g.Head.HasKeyword(keyword) {
		item := g.Head
		g.Head = nil
		return item
	}
	if g.Neck != nil && g.Neck.HasKeyword(keyword) {
		item := g.Neck
		g.Neck = nil
		return item
	}
	if g.Body != nil && g.Body.HasKeyword(keyword) {
		item := g.Body
		g.Body = nil
		return item
	}
	if g.Arms != nil && g.Arms.HasKeyword(keyword) {
		item := g.Arms
		g.Arms = nil
		return item
	}
	if g.Hands != nil && g.Hands.HasKeyword(keyword) {
		item := g.Hands
		g.Hands = nil
		return item
	}
	if g.Waist != nil && g.Waist.HasKeyword(keyword) {
		item := g.Waist
		g.Waist = nil
		return item
	}
	if g.Legs != nil && g.Legs.HasKeyword(keyword) {
		item := g.Legs
		g.Legs = nil
		return item
	}
	if g.Feet != nil && g.Feet.HasKeyword(keyword) {
		item := g.Feet
		g.Feet = nil
		return item
	}
	if g.Wrist != nil && g.Wrist.HasKeyword(keyword) {
		item := g.Wrist
		g.Wrist = nil
		return item
	}
	if g.Fingers != nil && g.Fingers.HasKeyword(keyword) {
		item := g.Fingers
		g.Fingers = nil
		return item
	}
	if g.OffHand != nil && g.OffHand.HasKeyword(keyword) {
		item := g.OffHand
		g.OffHand = nil
		return item
	}
	if g.MainHand != nil && g.MainHand.HasKeyword(keyword) {
		item := g.MainHand
		g.MainHand = nil
		return item
	}
	return nil
}

func (g *Gear) Equip(item *Item) (*Item, error) {
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
	case Hands:
		previous = g.Hands
		g.Hands = item
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
