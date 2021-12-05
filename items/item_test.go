package item

import (
	"github.com/michaelvmata/path/modifiers"
	"testing"
)

func TestItem(t *testing.T) {
	item := NewItem("f7b83201941a422f95100ac174be587f", "test helmet", Armor, Head, []string{"test", "helmet"})
	if len(item.Modifiers) != 0 {
		t.Fatalf("Default item has non empty modifiers")
	}
	item.AddModifier(modifiers.Dagger, 1)
	if len(item.Modifiers) != 1 {
		t.Fatalf("Item does not have modifier")
	}
	if !item.HasKeyword("helmet") {
		t.Fatalf("Item doesn't have keyword")
	}
}

func TestGear(t *testing.T) {
	gear := NewGear()
	helmet := NewItem("f7b83201941a422f95100ac174be587f", "test helmet", Armor, Head, []string{})
	if _, err := gear.Equip(helmet); err != nil {
		t.Fatalf("Error equiping %v", helmet)
	}
	if previous, err := gear.Equip(helmet); previous != helmet || err != nil {
		t.Fatalf("Error getting previous %v %v", helmet, err)
	}
	tablet := NewItem("329f203e98c64a2fa511385f55a7abcb", "test tablet", Tablet, Empty, []string{})
	if _, err := gear.Equip(tablet); err == nil {
		t.Fatalf("No Error equiping %v", tablet)
	}
}
