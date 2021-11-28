package stats

import "testing"

func TestStats(t *testing.T) {
	c := NewCores()
	if c.Power().Current != 0 {
		t.Fatalf("Core power zero value is not 0")
	}
	c.Will()
}
