package main

import (
	"fmt"
	"github.com/michaelvmata/path/stats"
	"github.com/michaelvmata/path/symbols"
	"strings"
)

type SkillType int

const (
	DAGGER SkillType = iota
	SWORD
	SPEAR
)

func (st SkillType) String() string {
	return symbols.TRIANGULAR_BULLET + " " + []string{"Dagger", "Sword", "Spear"}[st]
}

type Skills map[SkillType]*stats.Stat

func (s Skills) Describe() string {
	parts := make([]string, 0)
	for i := DAGGER; i <= SPEAR; i++ {
		description := fmt.Sprintf("%s: %d", i.String(), s[i].Value())
		parts = append(parts, description)
	}
	return strings.Join(parts, "\n")
}

func NewSkills() Skills {
	s := make(Skills)
	for i := DAGGER; i <= SPEAR; i++ {
		s[i] = stats.NewStat(0, 0)
	}
	return s
}
