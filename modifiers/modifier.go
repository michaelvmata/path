package modifiers

type Type int

const (
	Power Type = iota
	Agility
	Endurance
	Talent
	Insight
	Will

	Dagger
	Sword
	Spear
)

type Modifier struct {
	Type  Type
	Value int
}
