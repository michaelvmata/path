package buffs

import (
	"github.com/michaelvmata/path/world"
	"testing"
)

func TestAllBuffs(t *testing.T) {
	character := world.NewPlayer("Test UUID", "Test Handle")
	buffs := []world.Buff{
		NewBarrier(character),
		NewHaste(character),
	}
	for _, buff := range buffs {
		if buff.Name() == "" {
			t.Fatalf("%v name not set", buff)
		}
		if buff.ApplyMessage() == "" {
			t.Fatalf("%s apply not set", buff.Name())
		}
		if buff.UnapplyMessage() == "" {
			t.Fatalf("%s unapply not set", buff.Name())
		}
		if buff.AlreadyApplied() == "" {
			t.Fatalf("%s already applied not set", buff.Name())
		}
		if buff.Upkeep() <= 0 {
			t.Fatalf("%s upkeep not set", buff.Name())
		}

		remaining := buff.Remaining()
		buff.Update(1)
		if remaining == buff.Remaining() {
			t.Fatalf("%s remaining constant", buff.Name())
		}
	}
}
