package skills

import (
	"fmt"
	"github.com/michaelvmata/path/stats"
	"github.com/michaelvmata/path/symbols"
	"strings"
)

type Skills struct {
	Dagger stats.Stat
	Sword  stats.Stat
	Spear  stats.Stat

	Evasion stats.Stat
	Parry   stats.Stat

	Barrier  stats.Stat
	Backstab stats.Stat
	Bash     stats.Stat
	Bleed    stats.Stat
	Blitz    stats.Stat
	Circle   stats.Stat
	Haste    stats.Stat
}

func (s Skills) Describe() string {
	parts := []string{
		fmt.Sprintf("%s Dagger: %d", symbols.TRIANGULAR_BULLET, s.Dagger.Value()),
		fmt.Sprintf("%s Sword: %d", symbols.TRIANGULAR_BULLET, s.Sword.Value()),
		fmt.Sprintf("%s Spear: %d", symbols.TRIANGULAR_BULLET, s.Spear.Value()),
		fmt.Sprintf("%s Evasion: %d", symbols.TRIANGULAR_BULLET, s.Evasion.Value()),
		fmt.Sprintf("%s Parry: %d", symbols.TRIANGULAR_BULLET, s.Parry.Value()),
		fmt.Sprintf("%s Backstab: %d", symbols.TRIANGULAR_BULLET, s.Backstab.Value()),
		fmt.Sprintf("%s Barrier: %d", symbols.TRIANGULAR_BULLET, s.Barrier.Value()),
		fmt.Sprintf("%s Bash: %d", symbols.TRIANGULAR_BULLET, s.Bash.Value()),
		fmt.Sprintf("%s Bleed: %d", symbols.TRIANGULAR_BULLET, s.Bleed.Value()),
		fmt.Sprintf("%s Blitz: %d", symbols.TRIANGULAR_BULLET, s.Blitz.Value()),
		fmt.Sprintf("%s Circle: %d", symbols.TRIANGULAR_BULLET, s.Circle.Value()),
		fmt.Sprintf("%s Haste: %d", symbols.TRIANGULAR_BULLET, s.Haste.Value()),
	}
	return strings.Join(parts, "\n")
}

func NewSkills() Skills {
	return Skills{
		Dagger: stats.NewStat(0, 0),
		Sword:  stats.NewStat(0, 0),
		Spear:  stats.NewStat(0, 0),

		Evasion: stats.NewStat(0, 0),
		Parry:   stats.NewStat(0, 0),

		Backstab: stats.NewStat(0, 0),
		Bash:     stats.NewStat(0, 0),
		Bleed:    stats.NewStat(0, 0),
		Blitz:    stats.NewStat(0, 0),
		Circle:   stats.NewStat(0, 0),
		Haste:    stats.NewStat(0, 0),
	}
}
