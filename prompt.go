package main

import (
	"fmt"
	"github.com/michaelvmata/path/symbols"
)

type Prompt struct {
	Session *Session
}

func (p *Prompt) Render() string {
	border := "<grey_62>>> "
	if !p.Session.HasPlayer() {
		return border
	}
	player := p.Session.player
	return fmt.Sprintf("%s%s <red>%d%s <green>%d%s %s",
		border,
		player.Name,
		player.Health.Current, symbols.HEART,
		player.Spirit.Current, symbols.TWELVE_STAR,
		border)
}

func NewPrompt(session *Session) *Prompt {
	return &Prompt{Session: session}
}
