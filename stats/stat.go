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
	total := s.Base + s.Modifier
	if total <= 0 {
		total = 1
	}
	return total
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
