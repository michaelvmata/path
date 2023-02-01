package main

import (
	"github.com/michaelvmata/path/actions"
	"github.com/michaelvmata/path/battle"
	"github.com/michaelvmata/path/events"
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
	done := make(chan bool)
	go handleInput(s.Incoming)
	go handleOutput(s, done, player)
	ctx := Context{
		World:  world,
		Player: player,
	}
	world.SpawnMobiles()
	events.CharacterDeath.Register(actions.RespawnCharacter{})
	events.CharacterDeath.Register(actions.EssenceOnDeath{})

	player.ShowPrompt()
MainLoop:
	for {
		select {
		case text := <-s.Incoming:
			if text == "quit" {
				done <- true
				break MainLoop
			}
			command := determineCommand(text, ctx)
			player.Update(0)
			ctx.Raw = text
			command.Execute(ctx)
			player.ShowPrompt()
		case <-ticker.C:
			world.Update()
			if world.IsBattleTick() {
				battle.Simulate(world)
			}
		}

	}
	ticker.Stop()
}
