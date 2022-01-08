package stats

import "testing"

func TestCore(t *testing.T) {
	c := NewCore()

	oldBase := c.Power.Base
	c.Power.Increment()
	if c.Power.Base <= oldBase {
		t.Fatalf("Power increment didn't increase base")
	}

	oldModifier := c.Power.Modifier
	c.Power.Modify(1)
	if c.Power.Modifier <= oldModifier {
		t.Fatalf("Power modify didn't change the value")
	}
	c.Power.Reset()
	if c.Power.Modifier != oldModifier {
		t.Fatalf("Reset didn't set to initial state")
	}
	c.ResetModifier()
	c.Describe()
}
