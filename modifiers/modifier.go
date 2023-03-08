package modifiers

const (
	Power   = "Power"
	Agility = "Agility"
	Insight = "Insight"
	Will    = "Will"
)

type Modifier struct {
	Type  string
	Value int
}
