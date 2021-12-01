package item

import "github.com/michaelvmata/path/skills"

type Type int

const (
	Armor Type = iota
	Weapon
)

type Item struct {
	Name      string
	Type      Type
	Slot      Slot
	Modifiers []skills.Modifier
}

func (i *Item) AddModifier(skill skills.SkillType, value int) {
	modifier := skills.Modifier{
		Skill: skill,
		Value: value,
	}
	i.Modifiers = append(i.Modifiers, modifier)
}

func NewItem(name string, itemType Type, slot Slot) Item {
	return Item{
		Name:      name,
		Type:      itemType,
		Slot:      slot,
		Modifiers: make([]skills.Modifier, 0),
	}
}

type Slot int

const (
	Head Slot = iota
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

type Worn []Item

func NewWorn() Worn {
	e := make([]Item, MainHand-1)
	return e
}
