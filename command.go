package main

import (
	"fmt"
	"github.com/michaelvmata/path/buffs"
	"github.com/michaelvmata/path/items"
	"github.com/michaelvmata/path/symbols"
	"github.com/michaelvmata/path/world"
	"math"
	"strings"
)

type Context struct {
	World  *world.World
	Player *world.Character
	Raw    string
}

type Attack struct{}

func (a Attack) Execute(ctx Context) {
	attacker := ctx.Player
	parts := strings.SplitN(ctx.Raw, " ", 2)
	if len(parts) == 1 {
		attacker.Showln("Attack who?")
		return
	}
	handle := parts[1]
	index := attacker.Room.IndexOfPlayerHandle(handle)
	if index == -1 {
		attacker.Showln("You don't see '%s'.", handle)
		return
	}
	defender := attacker.Room.Players[index]
	attacker.StartAttacking(defender)
	attacker.Showln("You attack %s", defender.Name)
	defender.StartAttacking(attacker)
	defender.Showln("%s attacks you.", attacker.Name)
}

func (a Attack) Label() string {
	return "attack"
}

type Bash struct{}

func (b Bash) Execute(ctx Context) {
	attacker := ctx.Player
	level := attacker.Skills.Bash.Value()
	if level <= 0 {
		attacker.Showln("You try and fail.")
		return

	}
	if !attacker.Spirit.IsAvailable(level) {
		attacker.Showln("You're spirit isn't strong enough.")
		return
	}
	if attacker.OnCoolDown("bash") {
		attacker.Showln("It's on cool down.")
		return
	}
	attacker.Spirit.Consume(level)
	parts := strings.SplitN(ctx.Raw, " ", 2)
	if len(parts) == 1 {
		attacker.Showln("Bash who?")
		return
	}
	handle := parts[1]
	index := attacker.Room.IndexOfPlayerHandle(handle)
	if index == -1 {
		attacker.Showln("You don't see '%s'.", handle)
		return
	}
	defender := attacker.Room.Players[index]
	attacker.StartAttacking(defender)
	attacker.Showln("You bash %s.", defender.Name)
	defender.StartAttacking(attacker)
	defender.Showln("%s bashes you.", attacker.Name)
	defender.Stun(1)

	coolDown := buffs.NewCoolDown(9, "bash")
	attacker.ApplyCoolDown(&coolDown)
}

func (b Bash) Label() string {
	return "bash"
}

type Drop struct{}

func (d Drop) Execute(ctx Context) {
	player := ctx.Player
	parts := strings.SplitN(ctx.Raw, " ", 2)
	if len(parts) == 1 {
		player.Showln("Drop what?")
		return
	}
	keyword := parts[1]
	i := player.Discard(keyword)
	if i == nil {
		player.Showln("You don't have '%s'.", keyword)
		return
	}
	if err := player.Room.Accept(i); err != nil {
		player.Showln("You can't drop %s.", i.Name())
		player.Receive(i)
		return
	}
	player.Showln("You drop %s.", i.Name())
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

	player.Showln("")
	player.Showln("     [Head]: %s", g.SafeName(gear.Head))
	player.Showln("     [Neck]: %s", g.SafeName(gear.Neck))
	player.Showln("     [Body]: %s", g.SafeName(gear.Body))
	player.Showln("     [Arms]: %s", g.SafeName(gear.Arms))
	player.Showln("    [Hands]: %s", g.SafeName(gear.Hands))
	player.Showln("    [Waist]: %s", g.SafeName(gear.Waist))
	player.Showln("     [Legs]: %s", g.SafeName(gear.Legs))
	player.Showln("     [Feet]: %s", g.SafeName(gear.Feet))
	player.Showln("     Wrist]: %s", g.SafeName(gear.Wrist))
	player.Showln("  [Fingers]: %s", g.SafeName(gear.Fingers))
	player.Showln(" [Off Hand]: %s", g.SafeName(gear.OffHand))
	player.Showln("[Main Hand]: %s", g.SafeName(gear.MainHand))
	player.Showln("")
}

func (g Gear) Label() string {
	return "gear"
}

type Get struct{}

func (g Get) Execute(ctx Context) {
	player := ctx.Player
	parts := strings.SplitN(ctx.Raw, " ", 2)
	if len(parts) == 1 {
		player.Showln("Get what?")
		return
	}
	keyword := parts[1]
	i, err := player.Room.PickupItem(keyword)
	if err != nil {
		player.Showln("You don't see '%s'.", keyword)
		return
	}
	if err := player.Receive(i); err != nil {
		player.Showln("You can't get %s.", i.Name())
		player.Room.Accept(i)
		return
	} else {
		player.Showln("You get %s.", i.Name())
	}
}

func (g Get) Label() string {
	return "get"
}

type Haste struct{}

func (h Haste) Execute(ctx Context) {
	player := ctx.Player
	if !player.Skills.Haste.IsAvailable() {
		player.Showln("You fail to move with haste.")
		return
	}
	buff := buffs.NewHaste(1)
	player.Apply(buff)
}

func (h Haste) Label() string {
	return "haste"
}

type Inventory struct{}

func (i Inventory) Execute(ctx Context) {
	player := ctx.Player
	player.Showln("")
	for _, i := range player.Inventory.Items {
		player.Showln(i.Name())
	}
	player.Showln("")
}

func (i Inventory) Label() string {
	return "inventory"
}

type Invest struct{}

func spendEssence(p *world.Character, amount int) bool {
	cost := essenceCost(amount)
	if p.Essence < cost {
		p.Showln("You need %d more essence.", cost-p.Essence)
		return false
	}
	p.Essence -= cost
	return true
}

func essenceCost(amount int) int {
	return int(math.Pow(1.2, float64(amount))) + amount
}

func (i Invest) Execute(ctx Context) {
	player := ctx.Player
	parts := strings.SplitN(ctx.Raw, " ", 2)
	if len(parts) == 1 {
		player.Showln("Invest what?")
		return
	}
	keyword := strings.ToLower(parts[1])
	core := &player.Core
	skills := &player.Skills
	switch keyword {
	case "power":
		if spendEssence(player, core.Power.Base) {
			core.Power.Increment()
			player.Showln("Power courses through you.")
		}
	case "agility":
		if spendEssence(player, core.Agility.Base) {
			core.Agility.Increment()
			player.Showln("Balance flows through you.")
		}
	case "insight":
		if spendEssence(player, core.Insight.Base) {
			core.Insight.Increment()
			player.Showln("The world becomes clearer.")
		}
	case "will":
		if spendEssence(player, core.Will.Base) {
			core.Will.Increment()
			player.Showln("Reality itself warps before you.")
		}
	case "evasion":
		if spendEssence(player, skills.Evasion.Base) {
			skills.Evasion.Increment()
			player.Showln("You'll evade more alacrity.")
		}
	case "parry":
		if spendEssence(player, skills.Parry.Base) {
			skills.Parry.Increment()
			player.Showln("You'll parry with ease.")
		}
	case "bash":
		if spendEssence(player, skills.Bash.Base) {
			skills.Bash.Increment()
			player.Showln("Your mastery of bash improves.")
		}
	case "haste":
		if spendEssence(player, skills.Haste.Base) {
			skills.Haste.Increment()
			player.Showln("Your mastery of haste improves.")
		}
	default:
		player.Showln("Invest what?")
	}

}

func (i Invest) Label() string {
	return "invest"
}

type Look struct{}

func (l Look) Execute(ctx Context) {
	ctx.Player.Showln("")
	ctx.Player.Showln(ctx.Player.Room.Describe(ctx.Player))
	ctx.Player.Showln("")
}

func (l Look) Label() string {
	return "look"
}

type Remove struct{}

func (r Remove) Execute(ctx Context) {
	player := ctx.Player
	parts := strings.SplitN(ctx.Raw, " ", 2)
	if len(parts) == 1 {
		player.Showln("Remove what?")
		return
	}
	keyword := parts[1]
	i := player.Gear.Remove(keyword)
	if i == nil {
		player.Showln("You don't have a '%s'", keyword)
		return
	}
	player.Inventory.AddItem(i)
	player.Showln("You remove a %s", i.Name())
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
	p.Showln("")
	p.Showln(strings.Join(parts, "   "))
	p.Showln("")
	p.Showln(p.Core.Describe())
	p.Showln("")
	p.Showln("<yellow>%s<reset> Essence %d", symbols.FIVE_STAR, p.Essence)
	p.Showln("")
	p.Showln(p.Skills.Describe())
	p.Showln("")
}

func (sc Score) Label() string {
	return "score"
}

type StunLocked struct{}

func (s StunLocked) Execute(ctx Context) {
	ctx.Player.Showln("You are too stunned.")
}

func (s StunLocked) Label() string {
	return ""
}

type Typo struct{}

func (t Typo) Execute(ctx Context) {
	ctx.Player.Showln("The typo monster strikes again")
}

func (t Typo) Label() string {
	return "typo"
}

type Noop struct{}

func (n Noop) Execute(ctx Context) {
	ctx.Player.Showln("")
}

func (n Noop) Label() string {
	return ""
}

type Wear struct{}

func (wr Wear) Execute(ctx Context) {
	player := ctx.Player
	parts := strings.SplitN(ctx.Raw, " ", 2)
	if len(parts) == 1 {
		player.Showln("Wear what?")
		return
	}
	index := player.Inventory.IndexOfItem(parts[1])
	if index == -1 {
		player.Showln("You don't have a '%s'", parts[1])
		return
	}
	i := player.Inventory.RemItemAtIndex(index)
	previous, err := player.Gear.Equip(i)
	if err != nil {
		player.Showln("You can't wear %s.", i.Name())
		player.Inventory.AddItem(i)
	} else {
		player.Showln("You wear %s", i.Name())
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
		Bash{},
		Drop{},
		Gear{},
		Get{},
		Haste{},
		Inventory{},
		Invest{},
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
	aliases[Noop{}.Label()] = Noop{}
	return aliases
}

func determineCommand(raw string, ctx Context) Executor {
	if ctx.Player.IsStunned() {
		return StunLocked{}
	}
	rawCmd := strings.SplitN(raw, " ", 2)[0]
	command, ok := commands[rawCmd]
	if !ok {
		return commands[Typo{}.Label()]
	}

	return command
}
