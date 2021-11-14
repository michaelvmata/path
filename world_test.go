package main

import (
	"testing"
)

func TestCharacter(t *testing.T) {
	handle := "Tester"
	c := NewCharacter(handle)
	if c.handle != handle {
		t.Fatalf("Character handle(%s) expect (%s)", c.handle, handle)
	}
}

func TestRoom(t *testing.T) {
	uuid := "b8712a40130e41dabb7e17adb2d1aef7"
	name := "The Void"
	description := "An unending abyss."
	size := 1
	r := NewRoom(uuid, name, description, size)
	if r.uuid != uuid {
		t.Fatalf("Rooom uuid(%s) expected(%s)", r.uuid, uuid)
	}
	if r.name != name {
		t.Fatalf("Room name(%s) expected(%s)", r.name, name)
	}
	if r.description != description {
		t.Fatalf("Room description(%s) expected(%s)", r.description, description)
	}
	if r.size != size {
		t.Fatalf("Room size(%d) expected(%d)", r.size, size)
	}

	c := NewCharacter("Tester")
	if r.IsFull() {
		t.Fatalf("Room is unexpectedly full")
	}
	if err := r.EnterCharacter(c); err != nil {
		t.Fatalf("Character unable to enter empty room")
	}
	if r.IndexOfCharacter(c) == -1 {
		t.Fatalf("Character not found in room")
	}
	if !r.IsFull() {
		t.Fatalf("Room is unexpectedly not full")
	}
	if err := r.EnterCharacter(c); err == nil {
		t.Fatalf("Character entered full room")
	}

	r.size += 1
	if r.IsFull() {
		t.Fatalf("Room is unexpectedly full after resize")
	}
	if err := r.EnterCharacter(c); err == nil {
		t.Fatalf("Character able to enter room twice")
	}

	if err := r.ExitCharacter(c); err != nil {
		t.Fatalf("Character unable to exit room")
	}
	if r.IndexOfCharacter(c) != -1 {
		t.Fatalf("Character present after exiting")
	}
	if err := r.ExitCharacter(c); err == nil {
		t.Fatalf("Character able to exit room twice")
	}
}
