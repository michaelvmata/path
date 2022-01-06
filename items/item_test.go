package item

import (
	"github.com/michaelvmata/path/modifiers"
	"testing"
)

func TestItem(t *testing.T) {
	item := NewArmor("f7b83201941a422f95100ac174be587f", "test helmet", Head, []string{"test", "helmet"})
	if len(item.Modifiers()) != 0 {
		t.Fatalf("Default item has non empty modifiers")
	}
	item.AddModifier(modifiers.Dagger, 1)
	if len(item.Modifiers()) != 1 {
		t.Fatalf("Armor does not have modifier")
	}
	if !item.HasKeyword("helmet") {
		t.Fatalf("Armor doesn't have keyword")
	}
}

func TestContainer(t *testing.T) {
	item := NewArmor("f7b83201941a422f95100ac174be587f", "test helmet", Head, []string{"test", "helmet"})
	container := NewContainer(1)
	if index := container.IndexOfItem("test"); index != -1 {
		t.Fatalf("Armor found in empty container")
	}
	if err := container.AddItem(item); err != nil {
		t.Fatalf("Unable to add item to container")
	}
	index := container.IndexOfItem("test")
	if index != 0 {
		t.Fatalf("Armor not found in container")
	}
	container.RemItemAtIndex(index)
	if index := container.IndexOfItem("test"); index != -1 {
		t.Fatalf("Armor still found after rem")
	}
}

func TestGear(t *testing.T) {
	gear := NewGear()
	helmet := NewArmor("f7b83201941a422f95100ac174be587f", "test helmet", Head, []string{})
	if _, err := gear.Equip(helmet); err != nil {
		t.Fatalf("Error equiping %v", helmet)
	}
	if previous, err := gear.Equip(helmet); previous != helmet || err != nil {
		t.Fatalf("Error getting previous %v %v", helmet, err)
	}
	weapon := NewWeapon("Test UUID", "Test weapon", []string{"test", "weapon"}, Crush)
	if _, err := gear.Equip(weapon); err != nil {
		t.Fatalf("Error equiping weapon %v", weapon)
	}
	if i := gear.Remove("test"); i == nil {
		t.Fatalf("Error removing test item")
	}

}
