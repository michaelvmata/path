package stats

import "testing"

func TestCore(t *testing.T) {
	c := NewCores()
	if c.Power().Current != 0 {
		t.Fatalf("Core power zero value is not 0")
	}
	c.Will()
}

func TestConsumable(t *testing.T) {
	c := NewConsumables()
	c.Health()
	c.Energy()
	c.Spirit()
}
