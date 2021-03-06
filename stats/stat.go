package stats

import (
	"fmt"
	"github.com/michaelvmata/path/symbols"
)

type Stat struct {
	Base     int
	Modifier int
}

func (s Stat) IsAvailable() bool {
	return s.Base > 0
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
	Power   Stat
	Agility Stat
	Insight Stat
	Will    Stat
}

func (c *Core) Describe() string {
	return fmt.Sprintf("%s Power %d  %s Agility %d  %s Insight %d  %s Will %d",
		symbols.HEAVY_GREEK_CROSS, c.Power.Value(),
		symbols.HEAVY_GREEK_CROSS, c.Agility.Value(),
		symbols.HEAVY_GREEK_CROSS, c.Insight.Value(),
		symbols.HEAVY_GREEK_CROSS, c.Will.Value())
}

func (c *Core) ResetModifier() {
	c.Power.Reset()
	c.Agility.Reset()
	c.Insight.Reset()
	c.Will.Reset()
}

func NewCore() Core {
	return Core{
		Power:   NewStat(1, 0),
		Agility: NewStat(1, 0),
		Insight: NewStat(1, 0),
		Will:    NewStat(1, 0),
	}
}

type Line struct {
	Maximum     int
	Current     int
	RecoverRate int
}

func (l *Line) EnforceMaximum() {
	if l.Current > l.Maximum {
		l.Current = l.Maximum
	}
}

func (l *Line) Recover() {
	if l.Current >= l.Maximum {
		return
	}
	l.Current += l.RecoverRate
	l.EnforceMaximum()
}

func (l *Line) IsAvailable(amount int) bool {
	return l.Current >= amount
}

func (l *Line) Consume(amount int) {
	l.Current -= amount
}
