package main

import (
	session2 "github.com/michaelvmata/path/session"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second)
	session := session2.NewSession()
	world := build()
	session.PlayerName = "gaigen"
	player := world.Players[session.PlayerName]
	prompt := make(chan bool)
	done := make(chan bool)
	go handleInput(session.Incoming)
	go handleOutput(session, prompt, done, player)
	prompt <- true

MainLoop:
	for {
		select {
		case text := <-session.Incoming:
			if text == "quit" {
				done <- true
				break MainLoop
			}
			command := determineCommand(text)
			player.Update(false)
			command.Execute(world, session, text)
			prompt <- true
		case <-ticker.C:
			world.Update(true)
		}

	}
	ticker.Stop()
}
