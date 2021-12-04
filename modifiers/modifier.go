package modifiers

const (
	Power     = "Power"
	Agility   = "Agility"
	Endurance = "Endurance"
	Talent    = "Talent"
	Insight   = "Insight"
	Will      = "Will"

	Dagger = "Dagger"
	Sword  = "Sword"
	Spear  = "Spear"
)

type Modifier struct {
	Type  string
	Value int
}
