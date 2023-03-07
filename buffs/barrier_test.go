package buffs

import (
	"github.com/michaelvmata/path/world"
	"testing"
)

func TestBarrier(t *testing.T) {
	character := world.NewPlayer("Test UUID", "Test Handle")
	character.Skills.Barrier.Increment()
	b := NewBarrier(character)
	if b.DamageReduction() == 0 {
		t.Fatalf("Barrier damage reduction is 0")
	}
}
