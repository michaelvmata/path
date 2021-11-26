package main

type SkillType int

const (
	DAGGER SkillType = iota
	SWORD
	SPEAR
)

func (st SkillType) String() string {
	return []string{"Dagger", "Sword", "Spear"}[st]
}

type Skills map[SkillType]*Stat

func NewSkills() Skills {
	s := make(Skills)
	for i := DAGGER; i <= SPEAR; i++ {
		s[i] = NewStat(0, 0)
	}
	return s
}
