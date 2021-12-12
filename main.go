package main

import (
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second)
	session := NewSession()
	world := build()
	session.player = world.Players["gaigen"]
	prompt := make(chan bool)
	done := make(chan bool)
	go handleInput(session.incoming)
	go handleOutput(session, prompt, done)
	prompt <- true

MainLoop:
	for {
		select {
		case text := <-session.incoming:
			if text == "quit" {
				done <- true
				break MainLoop
			}
			command := determineCommand(text)
			session.player.Update(false)
			command.Execute(world, session, text)
			prompt <- true
		case <-ticker.C:
			session.player.Update(true)
		}

	}
	ticker.Stop()
}
