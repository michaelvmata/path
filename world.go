package main

import (
	"errors"
	"strings"
)

type Stat struct {
	Base     int
	Modifier int
}

func NewStat(base int, modifier int) *Stat {
	return &Stat{
		Base:     base,
		Modifier: modifier,
	}
}

func (s *Stat) Value() int {
	total := s.Base + s.Modifier
	if total <= 0 {
		total = 1
	}
	return total
}

type Line struct {
	Natural int // Base limit
	Maximum int // Adjusted limit
	Current int
	Recover int
}

type Player struct {
	Name   string
	Room   *Room
	Health Line
	Spirit Line
	Energy Line
}

func NewPlayer(handle string) *Player {
	return &Player{
		Name:   handle,
		Health: Line{},
		Spirit: Line{},
		Energy: Line{},
	}
}

func (c *Player) Move(r *Room) {
	c.Room = r
}

type Room struct {
	uuid        string
	name        string
	description string
	players     []*Player
	size        int
}

func NewRoom(uuid string, name string, description string, size int) *Room {
	room := Room{
		uuid:        uuid,
		name:        name,
		description: description,
		size:        size,
		players:     make([]*Player, 0, size),
	}
	return &room
}

func (r *Room) Describe() string {
	parts := make([]string, 0)
	parts = append(parts, r.name)
	parts = append(parts, "")
	parts = append(parts, r.description)
	return strings.Join(parts, "\n")
}

func (r *Room) IsFull() bool {
	return r.size == len(r.players)
}

func (r *Room) Enter(c *Player) error {
	if r.IsFull() {
		return errors.New("room is full")
	}
	if r.IndexOfPlayer(c) != -1 {
		return errors.New("player already in room")
	}
	r.players = append(r.players, c)
	return nil
}

func (r *Room) Exit(c *Player) error {
	i := r.IndexOfPlayer(c)
	if i == -1 {
		return errors.New("player not in room")
	}
	copy(r.players[i:], r.players[:i+1])
	length := len(r.players) - 1
	r.players[length] = nil
	r.players = r.players[:length]
	return nil
}

func (r *Room) IndexOfPlayer(target *Player) int {
	for i, p := range r.players {
		if p == target {
			return i
		}
	}
	return -1
}

type World struct {
	Players map[string]*Player
	Rooms   map[string]*Room
}

func NewWorld() *World {
	w := World{
		Players: make(map[string]*Player),
		Rooms:   make(map[string]*Room),
	}
	return &w
}
