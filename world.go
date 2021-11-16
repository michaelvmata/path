package main

import (
	"errors"
)

type Character struct {
	Handle string
	Room   *Room
}

func NewCharacter(handle string) *Character {
	return &Character{
		Handle: handle,
	}
}

func (c *Character) Move(r *Room) {
	c.Room = r
}

type Room struct {
	uuid        string
	name        string
	description string
	characters  []*Character
	size        int
}

func NewRoom(uuid string, name string, description string, size int) *Room {
	room := Room{
		uuid:        uuid,
		name:        name,
		description: description,
		size:        size,
		characters:  make([]*Character, 0, size),
	}
	return &room
}

func (r *Room) IsFull() bool {
	return r.size == len(r.characters)
}

func (r *Room) EnterCharacter(c *Character) error {
	if r.IsFull() {
		return errors.New("room is full")
	}
	if r.IndexOfCharacter(c) != -1 {
		return errors.New("character already in room")
	}
	r.characters = append(r.characters, c)
	return nil
}

func (r *Room) ExitCharacter(c *Character) error {
	i := r.IndexOfCharacter(c)
	if i == -1 {
		return errors.New("character not in room")
	}
	copy(r.characters[i:], r.characters[:i+1])
	length := len(r.characters) - 1
	r.characters[length] = nil
	r.characters = r.characters[:length]
	return nil
}

func (r *Room) IndexOfCharacter(target *Character) int {
	for i, c := range r.characters {
		if c == target {
			return i
		}
	}
	return -1
}

type World struct {
	Characters map[string]*Character
	Rooms      map[string]*Room
}

func NewWorld() *World {
	w := World{
		Characters: make(map[string]*Character),
		Rooms:      make(map[string]*Room),
	}
	return &w
}
