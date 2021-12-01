package item

import (
	"github.com/michaelvmata/path/skills"
	"testing"
)

func TestItem(t *testing.T) {
	item := NewItem("test helmet", Armor, Head)
	if len(item.Modifiers) != 0 {
		t.Fatalf("Default item has non empty modifiers")
	}
	item.AddModifier(skills.DAGGER, 1)
	if len(item.Modifiers) != 1 {
		t.Fatalf("Item does not have modifier")
	}
}

func TestWorn(t *testing.T) {
	worn := NewWorn()
	if len(worn) == 0 {
		t.Fatalf("Worn is empty")
	}

}
