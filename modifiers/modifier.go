package modifiers

const (
	Power   = "Power"
	Agility = "Agility"
	Insight = "Insight"
	Will    = "Will"

	Dagger = "Dagger"
	Sword  = "Sword"
	Spear  = "Spear"
)

type Modifier struct {
	Type  string
	Value int
}
