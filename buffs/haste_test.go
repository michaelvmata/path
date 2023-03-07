package buffs

import (
	"github.com/michaelvmata/path/world"
	"testing"
)

func TestHaste(t *testing.T) {
	character := world.NewPlayer("Test UUID", "Test Handle")
	character.Skills.Haste.Increment()
	h := NewHaste(character)
	if h.NumberOfAttacks() <= 0 {
		t.Fatalf("Invalid number of attacks")
	}
}
