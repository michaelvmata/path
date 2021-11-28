package stats

type Stat struct {
	Base     int
	Modifier int
}

func NewStat(base int, modifier int) *Stat {
	return &Stat{
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

type ConsumableType int

const (
	Health ConsumableType = iota
	Energy
	Spirit
)

var ConsumableLabels = []string{"Health", "Energy", "Spirit"}

func (ct ConsumableType) String() string {
	return ConsumableLabels[ct]
}

type Consumables []Stat

func NewConsumables() Consumables {
	c := make([]Stat, len(ConsumableLabels))
	return c
}

func (c Consumables) Health() *Stat {
	return &c[Health]
}

func (c Consumables) Energy() *Stat {
	return &c[Energy]
}

func (c Consumables) Spirit() *Stat {
	return &c[Spirit]
}

type Line struct {
	Natural int // Base limit
	Maximum int // Adjusted limit
	Current int
	Recover int
}

type CoreType int

const (
	Power CoreType = iota
	Agility
	Endurance
	Talent
	Insight
	Will
)

var CoreLabels = []string{"Power", "Agility", "Endurance", "Talent", "Insight", "Will"}

func (ct CoreType) String() string {
	return CoreLabels[ct]
}

type Cores []Line

func NewCores() Cores {
	c := make([]Line, len(CoreLabels))
	return c
}

func (c Cores) Power() Line {
	return c[Power]
}

func (c Cores) Agility() Line {
	return c[Agility]
}

func (c Cores) Endurance() Line {
	return c[Endurance]
}

func (c Cores) Talent() Line {
	return c[Talent]
}

func (c Cores) Insight() Line {
	return c[Insight]
}

func (c Cores) Will() Line {
	return c[Will]
}
