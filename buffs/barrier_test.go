package buffs

import "testing"

func TestBarrier(t *testing.T) {
	b := NewBarrier(1)
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
