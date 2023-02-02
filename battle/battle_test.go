package battle

import (
	"github.com/michaelvmata/path/buffs"
	item "github.com/michaelvmata/path/items"
	"github.com/michaelvmata/path/world"
	"testing"
)

func TestCalculateHitDamage(t *testing.T) {
	attacker := world.NewPlayer("Test Attacker", "Test Handle")
	w := item.NewWeapon("TestUUID", "Test Weapon", []string{"test"}, item.Crush)
	w.MaximumDamage = 10
	w.MinimumDamage = 5
	w.CriticalBonus = 0.0
	w.CriticalRate = 0.0

	_, err := attacker.Gear.Equip(w)
	if err != nil {
		t.Fatalf("Unable to equip weapon")
	}

	defender := world.NewPlayer("TestUUID2", "Test Defender")

	if damage := CalculateHitDamage(attacker, defender); damage.Amount <= 0 || damage.Critical == true {
		t.Fatalf("Damage max(%d), min(%d), actual(%d)", w.MaximumDamage, w.MinimumDamage, damage.Amount)
	}

	w.CriticalRate = 1.0
	w.CriticalBonus = 1.0
	if damage := CalculateHitDamage(attacker, defender); damage.Amount <= 0 || damage.Critical == false {
		t.Fatalf("Damage max(%d), min(%d), actual(%d)", w.MaximumDamage*2.0, w.MinimumDamage*2.0, damage.Amount)
	}
}

func TestNumberOfAttacks(t *testing.T) {
	character := world.NewPlayer("Test UUID", "Test Handle")

	buff := buffs.NewHaste(1)
	character.Core.Will.Modify(buff.Upkeep())
	character.Spirit.Maximum = buff.Upkeep()
	character.Spirit.Current = buff.Upkeep()
	character.Apply(buff)
	if NumberOfAttacks(character) != 2 {
		t.Fatalf("Haste didn't increase number of attacks")
	}
}

func TestEvade(t *testing.T) {
	character := world.NewPlayer("Test UUID", "Test Handle")
	character.Core.Agility.Increment()
	if ShouldEvade(character) {
		t.Fatalf("Character evaded without skill")
	}
	character.Skills.Evasion.Modifier = 100
	if !ShouldEvade(character) {
		t.Fatalf("Character with 100 evasion failed to evade")
	}
	DoEvade(character, character)
}

func TestDoDamage(t *testing.T) {
	character := world.NewPlayer("Test UUID", "Test Handle")
	DoDamage(nil, character, character, 1)
}
