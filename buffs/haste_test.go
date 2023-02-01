package buffs

import "testing"

func TestHaste(t *testing.T) {
	h := NewHaste(1)
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
