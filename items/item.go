package item

import (
	"errors"
	"fmt"
	"github.com/michaelvmata/path/modifiers"
	"strings"
)

const (
	PortalType = "Portal"
	WeaponType = "Weapon"
	ArmorType  = "Armor"
)

type item struct {
	uuid        string
	name        string
	keywords    []string
	description string
	modifiers   []modifiers.Modifier
	itemType    string
	immovable   bool
}

func (i *item) UUID() string {
	return i.uuid
}

func (i *item) Name() string {
	return i.name
}

func (i *item) HasKeyword(keyword string) bool {
	for _, candidate := range i.keywords {
		if candidate == keyword {
			return true
		}
	}
	return false
}

func (i *item) AddModifier(modifierType string, value int) {
	modifier := modifiers.Modifier{
		Type:  modifierType,
		Value: value,
	}
	i.modifiers = append(i.modifiers, modifier)
}

func (i *item) Modifiers() []modifiers.Modifier {
	return i.modifiers
}

func (i *item) Description() string {
	parts := make([]string, 0)
	parts = append(parts, fmt.Sprintf("<white>Name: %s", i.Name()))
	parts = append(parts, fmt.Sprintf("Keywords: %s", strings.Join(i.keywords, ", ")))
	for _, modifier := range i.Modifiers() {
		parts = append(parts, fmt.Sprintf("%s: %d", modifier.Type, modifier.Value))
	}
	parts = append(parts, "<reset>")
	parts = append(parts, i.description)
	return strings.Join(parts, "\n")
}

func (i *item) Immovable() bool {
	return i.immovable
}

func (i *item) MakeImmovable() {
	i.immovable = true
}

func (i *item) Type() string {
	return i.itemType
}

type Item interface {
	UUID() string
	Name() string
	HasKeyword(string) bool
	Modifiers() []modifiers.Modifier
	AddModifier(string, int)
	Description() string
	Type() string
	Immovable() bool
	MakeImmovable()
}

const (
	Crush  = "crush"
	Pierce = "pierce"
	Slash  = "slash"
)

const (
	Blade     = "blade"
	Impact    = "impact"
	TwoHanded = "two handed"
	Ranged    = "ranged"
)

type Other struct {
	item
}

func NewItem(UUID string, name string, keywords []string, description string, Type string) *Other {
	return &Other{
		item: item{
			uuid:        UUID,
			name:        name,
			keywords:    keywords,
			description: description,
			modifiers:   make([]modifiers.Modifier, 0),
			itemType:    Type,
		},
	}
}

type Weapon struct {
	item
	DamageType    string
	Attributes    []string
	MaximumDamage int
	MinimumDamage int
	CriticalBonus float64
	CriticalRate  float64
}

func (w *Weapon) IsBlade() bool {
	return w.HasAttribute(Blade)
}

func (w *Weapon) IsImpact() bool {
	return w.HasAttribute(Impact)
}

func (w *Weapon) IsTwoHanded() bool {
	return w.HasAttribute(TwoHanded)
}

func (w *Weapon) IsRange() bool {
	return w.HasAttribute(Ranged)
}

func (w *Weapon) HasAttribute(target string) bool {
	for _, attribute := range w.Attributes {
		if attribute == target {
			return true
		}
	}
	return false
}

func NewWeapon(UUID string, name string, keywords []string, description string, damageType string, attributes []string) *Weapon {
	return &Weapon{
		item: item{
			uuid:        UUID,
			name:        name,
			keywords:    keywords,
			description: description,
			modifiers:   make([]modifiers.Modifier, 0),
			itemType:    WeaponType,
		},
		DamageType: damageType,
		Attributes: attributes,
	}
}

type Armor struct {
	item
	Slot string
}

func NewArmor(UUID string, name string, slot string, keywords []string, description string) *Armor {
	return &Armor{
		item: item{
			uuid:        UUID,
			name:        name,
			keywords:    keywords,
			description: description,
			modifiers:   make([]modifiers.Modifier, 0),
			itemType:    ArmorType,
		},
		Slot: slot,
	}
}

type Container struct {
	Capacity int
	Items    []Item
}

func NewContainer(capacity int) Container {
	return Container{
		Capacity: capacity,
		Items:    make([]Item, 0),
	}
}

func (c *Container) AddItem(item Item) error {
	if len(c.Items) == c.Capacity {
		return errors.New("container full")
	}
	c.Items = append(c.Items, item)
	return nil
}

func (c *Container) IndexOfItem(keyword string) int {
	for i, item := range c.Items {
		if item.HasKeyword(keyword) {
			return i
		}
	}
	return -1
}

func (c *Container) RemItemAtIndex(index int) Item {
	item := c.Items[index]
	c.Items = append(c.Items[:index], c.Items[index+1:]...)
	return item
}

func (c *Container) GetItemAtIndex(index int) Item {
	return c.Items[index]
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
	OffHand  = "OffHand"
	MainHand = "MainHand"
)

type Gear struct {
	Head     *Armor
	Neck     *Armor
	Body     *Armor
	Arms     *Armor
	Hands    *Armor
	Waist    *Armor
	Legs     *Armor
	Feet     *Armor
	Wrist    *Armor
	Fingers  *Armor
	OffHand  *Armor
	MainHand *Weapon
}

func NewGear() *Gear {
	return &Gear{}
}

func (g *Gear) Remove(keyword string) Item {
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

func (g *Gear) Equip(i Item) (Item, error) {
	var previous Item
	if weapon, ok := i.(*Weapon); ok {
		previous = g.MainHand
		g.MainHand = weapon
		return previous, nil
	}
	armor, ok := i.(*Armor)
	if !ok {
		return nil, errors.New("not wearable")
	}
	switch armor.Slot {
	case Head:
		previous = g.Head
		g.Head = armor
	case Neck:
		previous = g.Neck
		g.Neck = armor
	case Body:
		previous = g.Body
		g.Body = armor
	case Arms:
		previous = g.Arms
		g.Arms = armor
	case Hands:
		previous = g.Hands
		g.Hands = armor
	case Waist:
		previous = g.Waist
		g.Waist = armor
	case Legs:
		previous = g.Legs
		g.Legs = armor
	case Feet:
		previous = g.Feet
		g.Feet = armor
	case Wrist:
		previous = g.Wrist
		g.Wrist = armor
	case Fingers:
		previous = g.Fingers
		g.Fingers = armor
	case OffHand:
		previous = g.OffHand
		g.OffHand = armor
	}
	return previous, nil
}

func IsNil(i Item) bool {
	if a, ok := i.(*Armor); ok && a == nil {
		return true
	}
	if w, ok := i.(*Weapon); ok && w == nil {
		return true
	}
	return false
}
