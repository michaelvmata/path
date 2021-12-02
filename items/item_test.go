package item

import (
	"github.com/michaelvmata/path/modifiers"
	"testing"
)

func TestItem(t *testing.T) {
	item := NewItem("test helmet", Armor, Head)
	if len(item.Modifiers) != 0 {
		t.Fatalf("Default item has non empty modifiers")
	}
	item.AddModifier(modifiers.Dagger, 1)
	if len(item.Modifiers) != 1 {
		t.Fatalf("Item does not have modifier")
	}
}

func TestGear(t *testing.T) {
	gear := NewGear()
	if len(gear) == 0 {
		t.Fatalf("Gear is empty")
	}
	helmet := NewItem("test helmet", Armor, Head)
	if _, err := gear.Equip(helmet); err != nil {
		t.Fatalf("Error equiping %v", helmet)
	}
	if previous, err := gear.Equip(helmet); previous != helmet || err != nil {
		t.Fatalf("Error getting previous %v %v", helmet, err)
	}
	tablet := NewItem("test tablet", Tablet, Empty)
	if _, err := gear.Equip(tablet); err == nil {
		t.Fatalf("No Error equiping %v", tablet)
	}
}
