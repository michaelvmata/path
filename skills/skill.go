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

	Haste stats.Stat
}

func (s Skills) Describe() string {
	parts := []string{
		fmt.Sprintf("%s Dagger: %d", symbols.TRIANGULAR_BULLET, s.Dagger.Value()),
		fmt.Sprintf("%s Sword: %d", symbols.TRIANGULAR_BULLET, s.Sword.Value()),
		fmt.Sprintf("%s Spear: %d", symbols.TRIANGULAR_BULLET, s.Spear.Value()),
		fmt.Sprintf("%s Haste: %d", symbols.TRIANGULAR_BULLET, s.Haste.Value()),
	}
	return strings.Join(parts, "\n")
}

func NewSkills() Skills {
	return Skills{
		Dagger: stats.NewStat(0, 0),
		Sword:  stats.NewStat(0, 0),
		Spear:  stats.NewStat(0, 0),

		Haste: stats.NewStat(0, 0),
	}
}
