package buffs

var BarrierName = "barrier"

type Barrier struct {
	CoolDown
	Level int
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

func NewBarrier(level int) *Barrier {
	return &Barrier{CoolDown: NewCoolDown(60, BarrierName), Level: level}
}
