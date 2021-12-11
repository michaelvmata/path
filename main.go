package main

import (
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second)
	ticks := 0
	session := NewSession()
	world := build()
	session.player = world.Players["gaigen"]
	prompt := make(chan bool)
	done := make(chan bool)
	go handleInput(session.incoming)
	go handleOutput(session, prompt, done)
	prompt <- true
	for {
		select {
		case text := <-session.incoming:
			command := determineCommand(text)
			session.player.Update(false)
			command.Execute(world, session, text)
			if text == "quit" {
				done <- true
				break
			}
			prompt <- true
		case <-ticker.C:
			ticks += 1
			isTock := ticks%20 == 0
			session.player.Update(isTock)
			if isTock == true {
				session.outgoing <- ""
				prompt <- true
			}
		}

	}
	ticker.Stop()
}
