package world

import (
	"github.com/michaelvmata/path/items"
	"github.com/michaelvmata/path/stats"
	"testing"
)

func TestStat(t *testing.T) {
	s := stats.NewStat(10, 10)
	if s.Value() != 20 {
		t.Fatalf("Stat value expected(%d) actual(%d)", 20, s.Value())
	}
}

func TestPlayer(t *testing.T) {
	UUID := "TestUUID"
	handle := "Tester"
	c := NewPlayer(UUID, handle)
	if c.Name != handle {
		t.Fatalf("Player handle(%s) expect (%s)", c.Name, handle)
	}

	uuid := "b8712a40130e41dabb7e17adb2d1aef7"
	name := "The Void"
	description := "An unending abyss."
	size := 1
	r := NewRoom(uuid, name, description, size)
	c.Move(r)
	if c.Room != r {
		t.Fatalf("Player room(%v) expected(%v)", c.Room, r)
	}

	target := "SomeUUID"
	if c.IsAttacking(target) {
		t.Fatalf("Player attacking %s by default", target)
	}
	c.StartAttacking(target)
	if !c.IsAttacking(target) {
		t.Fatalf("Player not attacking %s", target)
	}
	c.StopAttacking(target)
	if c.IsAttacking(target) {
		t.Fatalf("Player still attacking %s", target)
	}
}

func TestRoom(t *testing.T) {
	uuid := "b8712a40130e41dabb7e17adb2d1aef7"
	name := "The Void"
	description := "An unending abyss."
	size := 1
	r := NewRoom(uuid, name, description, size)
	if r.UUID != uuid {
		t.Fatalf("Rooom UUID(%s) expected(%s)", r.UUID, uuid)
	}
	if r.name != name {
		t.Fatalf("Room name(%s) expected(%s)", r.name, name)
	}
	if r.description != description {
		t.Fatalf("Room description(%s) expected(%s)", r.description, description)
	}
	if r.Size != size {
		t.Fatalf("Room Size(%d) expected(%d)", r.Size, size)
	}

	c := NewPlayer("Test UUID", "Tester")
	if r.IsFull() {
		t.Fatalf("Room is unexpectedly full")
	}
	if err := r.Enter(c); err != nil {
		t.Fatalf("Player unable to enter empty room")
	}
	if r.IndexOfPlayer(c) == -1 {
		t.Fatalf("Player not found in room")
	}
	if !r.IsFull() {
		t.Fatalf("Room is unexpectedly not full")
	}
	if err := r.Enter(c); err == nil {
		t.Fatalf("Player entered full room")
	}

	r.Size += 1
	if r.IsFull() {
		t.Fatalf("Room is unexpectedly full after resize")
	}
	if err := r.Enter(c); err == nil {
		t.Fatalf("Player able to enter room twice")
	}

	if err := r.Exit(c); err != nil {
		t.Fatalf("Player unable to exit room")
	}
	if r.IndexOfPlayer(c) != -1 {
		t.Fatalf("Player present after exiting")
	}
	if err := r.Exit(c); err == nil {
		t.Fatalf("Player able to exit room twice")
	}

	i := item.NewArmor("test UUID", "test item", item.Head, []string{"test"})
	err := r.Accept(i)
	if err != nil {
		t.Fatalf("Unable to drop item in room")
	}
	if _, ok := r.PickupItem("test"); ok != nil {
		t.Fatalf("Unable to pickup item in room")
	}
	if _, ok := r.PickupItem("test"); ok == nil {
		t.Fatalf("Able to pickup item in room, twice")
	}
}
