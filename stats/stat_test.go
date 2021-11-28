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

	oldBase := c.Health().Base
	c.Health().Increment()
	if c.Health().Base <= oldBase {
		t.Fatalf("Health increment didn't increase base")
	}

	oldModifier := c.Health().Modifier
	c.Health().Modify(1)
	if c.Health().Modifier <= oldModifier {
		t.Fatalf("Health modify didn't change the value")
	}
	c.Health().Reset()
	if c.Health().Modifier != oldModifier {
		t.Fatalf("Reset didn't set to initial state")
	}
}
