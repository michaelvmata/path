package main

import (
	"fmt"
	"github.com/michaelvmata/path/session"
	"github.com/michaelvmata/path/symbols"
	"github.com/michaelvmata/path/world"
)

type Prompt struct {
	Session *session.Session
	Player  *world.Character
}

func (p *Prompt) Render() string {
	border := "<grey_62>>> "
	if !p.Session.HasPlayer() {
		return border
	}
	return fmt.Sprintf("%s%s <red>%d%s <green>%d%s %s",
		border,
		p.Player.Name,
		p.Player.Health.Current, symbols.HEART,
		p.Player.Spirit.Current, symbols.TWELVE_STAR,
		border)
}

func NewPrompt(session *session.Session, p *world.Character) *Prompt {
	return &Prompt{Session: session, Player: p}
}
