package stats

type Stat struct {
	Base     int
	Modifier int
}

func NewStat(base int, modifier int) Stat {
	return Stat{
		Base:     base,
		Modifier: modifier,
	}
}

func (s *Stat) Value() int {
	// The value is a sum of the base and modifier.  It can be negative.
	total := s.Base + s.Modifier
	return total
}

func (s *Stat) Increment() {
	// The base value should only ever increase.  It should be a relatively
	// rare event.
	s.Base += 1
}

func (s *Stat) Modify(value int) {
	// The modifier can be positive or negative.  It will be a common event.
	s.Modifier += value
}

func (s *Stat) Reset() {
	// Reset the Modifier to a zero state.  Used when recalculating the total.
	s.Modifier = 0
}

type Core struct {
	Power     Stat
	Agility   Stat
	Endurance Stat
	Talent    Stat
	Insight   Stat
	Will      Stat
}

func NewCore() Core {
	return Core{
		Power:     NewStat(1, 0),
		Agility:   NewStat(1, 0),
		Endurance: NewStat(1, 0),
		Talent:    NewStat(1, 0),
		Insight:   NewStat(1, 0),
		Will:      NewStat(1, 0),
	}
}

type Line struct {
	Natural     int // Base limit
	Maximum     int // Adjusted limit
	Current     int
	RecoverRate int
}

func (l *Line) Recover() {
	if l.Current >= l.Maximum {
		return
	}
	diff := l.Maximum - l.Current
	if l.RecoverRate > diff {
		l.Current = l.Maximum
	} else {
		l.Current += l.RecoverRate
	}
}

type Consumable struct {
	Health Line
	Energy Line
	Spirit Line
}

func NewConsumable() Consumable {
	return Consumable{
		Health: Line{},
		Energy: Line{},
		Spirit: Line{},
	}
}
