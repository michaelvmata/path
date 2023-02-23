package buffs

import "github.com/michaelvmata/path/world"

var BarrierName = "barrier"

type Barrier struct {
	CoolDown
	Level     int
	Character *world.Character
}

func (b Barrier) ApplyMessage() string {
	return "A layer of spirit solidifies around you."
}

func (b Barrier) UnapplyMessage() string {
	return "A layer of spirit around you fades."
}

func (b Barrier) AlreadyApplied() string {
	return "Your spirit already protects you."
}

func (b Barrier) Upkeep() int {
	return b.Level
}

func (b Barrier) DamageReduction() int {
	return b.Level * 2
}

func (b Barrier) Update(tick int) {

}

func NewBarrier(character *world.Character) *Barrier {
	level := character.Skills.Barrier.Value()
	return &Barrier{CoolDown: NewCoolDown(60, BarrierName), Level: level, Character: character}
}
