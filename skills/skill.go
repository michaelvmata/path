package skills

import (
	"fmt"
	"github.com/michaelvmata/path/stats"
	"github.com/michaelvmata/path/symbols"
	"strings"
)

type SkillType int

const (
	POWER SkillType = iota
	AGILITY
	ENDURANCE
	TALENT
	INSIGHT
	WILL

	DAGGER
	SWORD
	SPEAR
)

func (st SkillType) String() string {
	return symbols.TRIANGULAR_BULLET + " " + []string{"Dagger", "Sword", "Spear"}[st]
}

type Skills map[SkillType]*stats.Stat

func (s Skills) Describe() string {
	parts := make([]string, 0)
	for i := POWER; i <= SPEAR; i++ {
		description := fmt.Sprintf("%s: %d", i.String(), s[i].Value())
		parts = append(parts, description)
	}
	return strings.Join(parts, "\n")
}

func NewSkills() Skills {
	s := make(Skills)
	for i := DAGGER; i <= SPEAR; i++ {
		stat := stats.NewStat(0, 0)
		s[i] = &stat
	}
	return s
}

type Modifier struct {
	Skill SkillType
	Value int
}
