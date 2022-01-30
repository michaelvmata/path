package buffs

var StunName = "haste"

type Stun struct {
	Level    int
	Lifetime int
}

func (s Stun) Name() string {
	return StunName
}

func (s *Stun) Update(tick int) {
	s.Lifetime--
}

func (s *Stun) Expire() {
	s.Lifetime = 0
}
func (s Stun) IsExpired() bool {
	return s.Lifetime <= 0
}

func (s Stun) ApplyMessage() string {
	return "You see stars."
}

func (s Stun) UnapplyMessage() string {
	return "You no longer see stars."
}

func (s Stun) AlreadyApplied() string {
	return ""
}

func (s Stun) Upkeep() int {
	return 0
}

func NewStun(level int) *Stun {
	return &Stun{Lifetime: level, Level: level}
}
