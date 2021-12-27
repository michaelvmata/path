package main

import (
	"fmt"
	"github.com/michaelvmata/path/items"
	"github.com/michaelvmata/path/symbols"
	"github.com/michaelvmata/path/world"
	"strings"
)

type Context struct {
	World  *world.World
	Player *world.Player
	Raw    string
}

type Attack struct{}

func (a Attack) Execute(ctx Context) {
	attacker := ctx.Player
	parts := strings.SplitN(ctx.Raw, " ", 2)
	if len(parts) == 1 {
		attacker.Show("Attack who?")
		return
	}
	handle := parts[1]
	index := attacker.Room.IndexOfPlayerHandle(handle)
	if index == -1 {
		attacker.Show("You don't see '%s'.", handle)
		return
	}
	defender := attacker.Room.Players[index]
	attacker.StartAttacking(defender.Name)
	attacker.Show("You attack %s", defender.Name)
	defender.StartAttacking(attacker.Name)
	defender.Show("%s attacks you.", attacker.Name)
}

func (a Attack) Label() string {
	return "attack"
}

type Drop struct{}

func (d Drop) Execute(ctx Context) {
	player := ctx.Player
	parts := strings.SplitN(ctx.Raw, " ", 2)
	if len(parts) == 1 {
		player.Show("Drop what?")
		return
	}
	keyword := parts[1]
	i := player.Discard(keyword)
	if i == nil {
		player.Show("You don't have '%s'.", keyword)
		return
	}
	if err := player.Room.Accept(i); err != nil {
		player.Show("You can't drop %s.", i.Name())
		player.Receive(i)
		return
	}
	player.Show("You drop %s.", i.Name())
}

func (d Drop) Label() string {
	return "drop"
}

type Gear struct{}

func (g Gear) SafeName(i item.Item) string {
	if item.IsNil(i) {
		return ""
	}
	return i.Name()
}

func (g Gear) Execute(ctx Context) {
	player := ctx.Player
	gear := player.Gear

	player.Show("")
	player.Show("     [Head]: %s", g.SafeName(gear.Head))
	player.Show("     [Neck]: %s", g.SafeName(gear.Neck))
	player.Show("     [Body]: %s", g.SafeName(gear.Body))
	player.Show("     [Arms]: %s", g.SafeName(gear.Arms))
	player.Show("    [Hands]: %s", g.SafeName(gear.Hands))
	player.Show("    [Waist]: %s", g.SafeName(gear.Waist))
	player.Show("     [Legs]: %s", g.SafeName(gear.Legs))
	player.Show("     [Feet]: %s", g.SafeName(gear.Feet))
	player.Show("     Wrist]: %s", g.SafeName(gear.Wrist))
	player.Show("  [Fingers]: %s", g.SafeName(gear.Fingers))
	player.Show(" [Off Hand]: %s", g.SafeName(gear.OffHand))
	player.Show("[Main Hand]: %s", g.SafeName(gear.MainHand))
	player.Show("")
}

func (g Gear) Label() string {
	return "gear"
}

type Get struct{}

func (g Get) Execute(ctx Context) {
	player := ctx.Player
	parts := strings.SplitN(ctx.Raw, " ", 2)
	if len(parts) == 1 {
		player.Show("Get what?")
		return
	}
	keyword := parts[1]
	i, err := player.Room.PickupItem(keyword)
	if err != nil {
		player.Show("You don't see '%s'.", keyword)
		return
	}
	if err := player.Receive(i); err != nil {
		player.Show("You can't get %s.", i.Name())
		player.Room.Accept(i)
		return
	} else {
		player.Show("You get %s.", i.Name())
	}
}

func (g Get) Label() string {
	return "get"
}

type Inventory struct{}

func (i Inventory) Execute(ctx Context) {
	player := ctx.Player
	player.Show("")
	for _, i := range player.Inventory.Items {
		player.Show(i.Name())
	}
	player.Show("")
}

func (i Inventory) Label() string {
	return "inventory"
}

type Look struct{}

func (l Look) Execute(ctx Context) {
	ctx.Player.Show("")
	ctx.Player.Show(ctx.Player.Room.Describe(ctx.Player))
	ctx.Player.Show("")
}

func (l Look) Label() string {
	return "look"
}

type Remove struct{}

func (r Remove) Execute(ctx Context) {
	player := ctx.Player
	parts := strings.SplitN(ctx.Raw, " ", 2)
	if len(parts) == 1 {
		player.Show("Remove what?")
		return
	}
	keyword := parts[1]
	i := player.Gear.Remove(keyword)
	if i == nil {
		player.Show("You don't have a '%s'", keyword)
		return
	}
	player.Inventory.AddItem(i)
	player.Show("You remove a %s", i.Name())
}

func (r Remove) Label() string {
	return "remove"
}

type Score struct{}

func (sc Score) Execute(ctx Context) {
	p := ctx.Player
	parts := make([]string, 0)
	parts = append(parts, fmt.Sprintf("<red>%s<reset> Health %d(%d)+%d", symbols.HEART, p.Health.Current, p.Health.Maximum, p.Health.RecoverRate))
	parts = append(parts, fmt.Sprintf("<green>%s<reset> Spirit %d(%d)+%d", symbols.TWELVE_STAR, p.Spirit.Current, p.Spirit.Maximum, p.Spirit.RecoverRate))
	p.Show("")
	p.Show(strings.Join(parts, "   "))
	p.Show("")
	p.Show(p.Core.Describe())
	p.Show("")
	p.Show(p.Skills.Describe())
	p.Show("")
}

func (sc Score) Label() string {
	return "score"
}

type Typo struct{}

func (t Typo) Execute(ctx Context) {
	ctx.Player.Show("The typo monster strikes again")
}

func (t Typo) Label() string {
	return "typo"
}

type Noop struct{}

func (n Noop) Execute(ctx Context) {
	ctx.Player.Show("")
}

func (n Noop) Label() string {
	return ""
}

type Wear struct{}

func (wr Wear) Execute(ctx Context) {
	player := ctx.Player
	parts := strings.SplitN(ctx.Raw, " ", 2)
	if len(parts) == 1 {
		player.Show("Wear what?")
		return
	}
	index := player.Inventory.IndexOfItem(parts[1])
	if index == -1 {
		player.Show("You don't have a '%s'", parts[1])
		return
	}
	i := player.Inventory.RemItemAtIndex(index)
	previous, err := player.Gear.Equip(i)
	if err != nil {
		player.Show("You can't wear %s.", i.Name())
		player.Inventory.AddItem(i)
	} else {
		player.Show("You wear %s", i.Name())
	}
	if !item.IsNil(previous) {
		player.Inventory.AddItem(previous)
	}
}

func (wr Wear) Label() string {
	return "wear"
}

type Executor interface {
	Execute(ctx Context)
	Label() string
}

var commands = buildCommands()

func buildCommands() map[string]Executor {
	commands := []Executor{
		Attack{},
		Drop{},
		Gear{},
		Get{},
		Inventory{},
		Look{},
		Noop{},
		Remove{},
		Score{},
		Typo{},
		Wear{},
	}
	aliases := make(map[string]Executor)
	for _, command := range commands {
		for i := range command.Label() {
			alias := command.Label()[:i+1]
			aliases[alias] = command
		}
	}
	return aliases
}

func determineCommand(raw string) Executor {
	rawCmd := strings.SplitN(raw, " ", 2)[0]
	command, ok := commands[rawCmd]
	if !ok {
		return commands[Typo{}.Label()]
	}

	return command
}
