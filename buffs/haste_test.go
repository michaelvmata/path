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
	if h.Upkeep() <= 0 {
		t.Fatalf("Invalid upkeep cost")
	}
	if h.ApplyMessage() == "" || h.UnapplyMessage() == "" || h.AlreadyApplied() == "" {
		t.Fatalf("Didn't set apply related message")
	}
}
