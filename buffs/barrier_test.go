package buffs

import (
	"github.com/michaelvmata/path/world"
	"testing"
)

func TestBarrier(t *testing.T) {
	character := world.NewPlayer("Test UUID", "Test Handle")
	character.Skills.Barrier.Increment()
	b := NewBarrier(character)
	if b.ApplyMessage() == "" {
		t.Fatalf("Barrier apply not set")
	}
	if b.UnapplyMessage() == "" {
		t.Fatalf("Barrier unapply not set")
	}
	if b.AlreadyApplied() == "" {
		t.Fatalf("Barrier already applied not set")
	}
	if b.Upkeep() <= 0 {
		t.Fatalf("Barrier upkeep not set")
	}
	if b.DamageReduction() == 0 {
		t.Fatalf("Barrier damage reduction is 0")
	}
}
