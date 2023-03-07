package buffs

import "github.com/michaelvmata/path/world"

var HasteName = "haste"

type Haste struct {
	CoolDown
	Level int
}

func (h Haste) ApplyMessage() string {
	return "The world slows perceptibly."
}

func (h Haste) UnapplyMessage() string {
	return "The world speeds up perceptibly."
}

func (h Haste) AlreadyApplied() string {
	return "You already move with haste."
}

func (h Haste) Upkeep() int {
	return h.Level + 1
}

func (h Haste) NumberOfAttacks() int {
	return h.Level
}

func NewHaste(character *world.Character) *Haste {
	level := character.Skills.Haste.Value()
	return &Haste{CoolDown: NewCoolDown(60, HasteName), Level: level}
}
