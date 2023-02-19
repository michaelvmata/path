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
		t.Fatalf("Character handle(%s) expect (%s)", c.Name, handle)
	}

	uuid := "b8712a40130e41dabb7e17adb2d1aef7"
	name := "The Void"
	description := "An unending abyss."
	size := 1
	r := NewRoom(uuid, name, description, size)
	c.Move(r)
	if c.Room != r {
		t.Fatalf("Character room(%v) expected(%v)", c.Room, r)
	}

	target := NewPlayer("SomeUUID", "Somehandle")
	if c.IsAttacking(target) {
		t.Fatalf("Character attacking %v by default", target)
	}
	c.StartAttacking(target)
	if !c.IsAttacking(target) {
		t.Fatalf("Character not attacking %v", target)
	}
	c.StopAttacking(target)
	if c.IsAttacking(target) {
		t.Fatalf("Character still attacking %v", target)
	}

	if !c.HasKeyword(handle) {
		t.Fatalf("Character handle is not a keyword")
	}
}

type TestBuff struct {
	Expired bool
}

func (t TestBuff) Update(tick int)        {}
func (t TestBuff) IsExpired() bool        { return t.Expired }
func (t *TestBuff) Expire()               { t.Expired = true }
func (t TestBuff) Name() string           { return "TestBuff" }
func (t TestBuff) ApplyMessage() string   { return "TestBuff" }
func (t TestBuff) UnapplyMessage() string { return "TestBuff" }
func (t TestBuff) AlreadyApplied() string { return "TestBuff" }
func (t TestBuff) Upkeep() int            { return 0 }

func TestPlayerBuff(t *testing.T) {
	player := NewPlayer("Test UUID", "Test Handle")
	buff := TestBuff{Expired: false}
	player.Apply(&buff)
	if len(player.Buffs) != 1 {
		t.Fatalf("Unable to apply buff %v", buff)
	}
	player.Update(1)
	if len(player.Buffs) != 1 {
		t.Fatalf("Buff removed before expired")
	}
	buff.Expire()
	player.Update(1)
	if len(player.Buffs) != 0 {
		t.Fatalf("Expired buff still applied")
	}
}

func TestRoom(t *testing.T) {
	uuid := "b8712a40130e41dabb7e17adb2d1aef7"
	name := "The Void"
	description := "An unending abyss."
	size := 2
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
	if r.GetPlayer(c.Name) != nil {
		t.Fatalf("Character unexpectedly present by handle")
	}
	if r.IsFull() {
		t.Fatalf("Room is unexpectedly full")
	}
	if err := r.Enter(c); err != nil {
		t.Fatalf("Character unable to enter empty room")
	}
	if r.IndexOfPlayer(c) == -1 {
		t.Fatalf("Character not found in room")
	}
	if r.GetPlayer(c.Name) != c {
		t.Fatalf("Character not found by handle")
	}

	c2 := NewPlayer("Test UUID", "Tester 2")
	if err := r.Enter(c2); err != nil {
		t.Fatalf("2nd character unable to enter room %v", err)
	}
	if r.IndexOfPlayer(c2) == -1 {
		t.Fatalf("Character not found in room")
	}
	if !r.IsFull() {
		t.Fatalf("Room is unexpectedly not full")
	}

	c3 := NewPlayer("Test UUID", "Tester 3")
	if err := r.Enter(c3); err == nil {
		t.Fatalf("3rd character entered full room")
	}

	r.Size += 1
	if r.IsFull() {
		t.Fatalf("Room is unexpectedly full after resize")
	}
	if err := r.Enter(c); err == nil {
		t.Fatalf("Character able to enter room twice")
	}

	if err := r.Exit(c); err != nil {
		t.Fatalf("Character unable to exit room")
	}
	if r.IndexOfPlayer(c) != -1 {
		t.Fatalf("Character present after exiting")
	}
	if err := r.Exit(c); err == nil {
		t.Fatalf("Character able to exit room twice")
	}

	r.Enter(c)
	r.Enter(c2)
	r.Enter(c3)
	message := Message{
		FirstPerson:         c,
		SecondPerson:        c2,
		FirstPersonMessage:  "1st person test",
		SecondPersonMessage: "2nd person test",
		ThirdPersonMessage:  "3rd person test",
	}
	if err := r.ShowMessage(message); err != nil {
		t.Fatalf("Unexpected error showing message %v", err)
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
