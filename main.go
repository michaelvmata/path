package main

import (
	"github.com/michaelvmata/path/session"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second)
	s := session.New()
	s.PlayerName = "gaigen"
	world := build()
	player := world.Players[s.PlayerName]
	player.Session = s
	prompt := make(chan bool)
	done := make(chan bool)
	go handleInput(s.Incoming)
	go handleOutput(s, prompt, done, player)
	prompt <- true
	ctx := Context{
		World:  world,
		Player: player,
	}
MainLoop:
	for {
		select {
		case text := <-s.Incoming:
			if text == "quit" {
				done <- true
				break MainLoop
			}
			command := determineCommand(text)
			player.Update(0)
			ctx.Raw = text
			command.Execute(ctx)
			prompt <- true
		case <-ticker.C:
			world.Update()
		}

	}
	ticker.Stop()
}
