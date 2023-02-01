package buffs

import "testing"

func TestCoolDown(t *testing.T) {
	name := "test cool down"
	lifetime := 2
	cd := NewCoolDown(lifetime, name)
	if cd.IsExpired() {
		t.Fatalf("Cool down immediately expired")
	}
	cd.Update(0)
	if cd.Lifetime >= lifetime {
		t.Fatalf("Lifetime didn't tick down")
	}
	cd.Expire()
	if !cd.IsExpired() {
		t.Fatalf("Cooldown didn't expire")
	}
	if cd.Name() != name {
		t.Fatalf("Cooldown name not set")
	}
}
