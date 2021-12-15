package battle

import (
	item "github.com/michaelvmata/path/items"
	"github.com/michaelvmata/path/world"
	"testing"
)

func TestCalculateHitDamage(t *testing.T) {
	attacker := world.NewPlayer("Test Attacker")
	w := item.NewWeapon("TestUUID", "Test Weapon", []string{"test"}, item.Crush)
	w.MaximumDamage = 10
	w.MinimumDamage = 5
	attacker.Gear.Equip(w)

	defender := world.NewPlayer("Test Defender")

	if damage := CalculateHitDamage(attacker, defender); damage < w.MinimumDamage || damage > w.MaximumDamage {
		t.Fatalf("Damage max(%d), min(%d), actual(%d)", w.MaximumDamage, w.MinimumDamage, damage)
	}

	w.CriticalRate = 1.0
	w.CriticalBonus = 1.0
	if damage := CalculateHitDamage(attacker, defender); damage < w.MinimumDamage*2.0 || damage > w.MaximumDamage*2.0 {
		t.Fatalf("Damage max(%d), min(%d), actual(%d)", w.MaximumDamage*2.0, w.MinimumDamage*2.0, damage)
	}
}
