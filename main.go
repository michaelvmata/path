package main

import (
	"github.com/michaelvmata/path/actions"
	"github.com/michaelvmata/path/events"
	"github.com/michaelvmata/path/help"
	"github.com/michaelvmata/path/session"
	"github.com/michaelvmata/path/simulate"
	"github.com/michaelvmata/path/title"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second)
	s := session.New()
	done := make(chan bool)
	go handleInput(s.Incoming)
	go handleOutput(s, done)
	w := build("data/areas")
	ctx := Context{
		World: w,
		Help:  help.Build("data/help"),
	}
	w.SpawnMobiles()
	events.CharacterDeath.Init(w)
	events.CharacterDeath.Register(actions.RespawnCharacter{})
	events.CharacterDeath.Register(actions.EssenceOnDeath{})
	events.CharacterDeath.Register(actions.QuestOnDeath{})
	title.ListCharacters(s, w.Players)
MainLoop:
	for {
		select {
		case text := <-s.Incoming:
			if text == "quit" {
				done <- true
				break MainLoop
			}
			if ctx.Player == nil {
				player, err := determinePlayer(text, w)
				if err == nil {
					ctx.Player = player
					ctx.Player.Session = s
					s.PlayerName = text
					ctx.Player.ShowPrompt()
				}
			} else {
				command := determineCommand(text, ctx)
				ctx.Raw = text
				command.Execute(ctx)
				ctx.Player.ShowPrompt()
			}

		case <-ticker.C:
			w.Update()
			if w.IsBattleTick() {
				simulate.Simulate(w)
			}
		}
	}
	ticker.Stop()
}
