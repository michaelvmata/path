package buffs

var HasteName = "haste"

type Haste struct {
	CoolDown
	Level int
}

func (h Haste) Name() string {
	return HasteName
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
	return h.Level
}

func (h Haste) NumberOfAttacks() int {
	return h.Level
}

func NewHaste(level int) *Haste {
	return &Haste{CoolDown: CoolDown{Lifetime: 60}, Level: level}
}
