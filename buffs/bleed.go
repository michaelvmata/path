package buffs

import (
	"github.com/michaelvmata/path/world"
)

var BleedName = "bleed"

type Bleed struct {
	CoolDown
	Level     int
	Character *world.Character
	Applier   *world.Character
}

func (b *Bleed) ApplyMessage() string {
	return "You start bleeding from a gaping wound."
}

func (b *Bleed) UnapplyMessage() string {
	return "Your wound stops bleeding."
}

func (b *Bleed) AlreadyApplied() string {
	return "Your gaping wound continues to bleed."
}

func (b *Bleed) Upkeep() int {
	return b.Level + 1
}

func NewBleed(character *world.Character, applier *world.Character) *Bleed {
	level := character.Skills.Bleed.Value()
	return &Bleed{CoolDown: NewCoolDown(6, BleedName), Level: level, Character: character, Applier: applier}
}
