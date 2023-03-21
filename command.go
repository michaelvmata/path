package main

import (
	"errors"
	"fmt"
	"github.com/michaelvmata/path/buffs"
	"github.com/michaelvmata/path/help"
	"github.com/michaelvmata/path/items"
	"github.com/michaelvmata/path/simulate"
	"github.com/michaelvmata/path/symbols"
	"github.com/michaelvmata/path/world"
	"log"
	"math"
	"sort"
	"strings"
)

type Context struct {
	World  *world.World
	Player *world.Character
	Help   map[string]help.YAMLHelp
	Raw    string
}

func IsWieldingRange(character *world.Character) bool {
	return !item.IsNil(character.Gear.MainHand) && character.Gear.MainHand.IsRange()
}

func IsWieldingBlade(attacker *world.Character) bool {
	return !item.IsNil(attacker.Gear.MainHand) && attacker.Gear.MainHand.IsBlade()
}

func FindTarget(attacker *world.Character, command string) *world.Character {
	parts := strings.SplitN(command, " ", 2)
	if len(parts) == 1 {
		return attacker.ImmediateDefender()
	}
	handle := parts[1]
	return attacker.Room.GetPlayer(handle)
}

func CanUseSkill(attacker *world.Character, skill string, level int, cost int) bool {
	if level == 0 {
		attacker.Showln("You gotta learn how to %s first.", skill)
		return false
	}
	if !attacker.Spirit.IsAvailable(cost) {
		attacker.Showln("Your spirit isn't strong enough to %s.", skill)
		return false
	}
	if attacker.OnCoolDown(skill) {
		attacker.Showln("You need a moment before you can %s again.", skill)
		return false
	}
	return true
}

func InitBattleSkill(attacker *world.Character, defender *world.Character, spirit int, skill string, coolDownDuration int) {
	attacker.Spirit.Consume(spirit)
	attacker.StartAttacking(defender)
	defender.StartAttacking(attacker)
	coolDown := buffs.NewCoolDown(coolDownDuration, skill)
	attacker.ApplyCoolDown(&coolDown)
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
	defender := attacker.Room.GetPlayer(handle)
	if defender == nil {
		attacker.Showln("You don't see '%s'.", handle)
		return
	}

	attacker.StartAttacking(defender)
	defender.StartAttacking(attacker)

	message := world.Message{
		FirstPerson:         attacker,
		FirstPersonMessage:  fmt.Sprintf("You attack %s", defender.Name),
		SecondPerson:        defender,
		SecondPersonMessage: fmt.Sprintf("%s attacks you.", attacker.Name),
		ThirdPersonMessage:  fmt.Sprintf("%s attacks %s.", attacker.Name, defender.Name),
	}
	attacker.Room.ShowMessage(message)
}

func (a Attack) Label() string {
	return "attack"
}

type Affect struct{}

func (a Affect) Execute(ctx Context) {
	player := ctx.Player
	if len(player.Buffs) == 0 {
		player.Showln("You are not affected by any technique.")
		return
	}

	player.Showln("Techniques")
	player.ShowDivider()
	player.ShowNewline()
	for _, buff := range player.Buffs {
		player.Showln(fmt.Sprintf("Name: %s\tDuration: %d", buff.Name(), buff.Remaining()))
	}
}

func (a Affect) Label() string {
	return "affect"
}

type Backstab struct{}

func (b Backstab) Execute(ctx Context) {
	attacker := ctx.Player
	level := attacker.Skills.Backstab.Value()
	if !CanUseSkill(attacker, b.Label(), level, level) {
		return
	}

	if attacker.IsFighting() {
		attacker.Showln("You can't backstab while in combat.")
		return
	}

	if !IsWieldingBlade(attacker) {
		attacker.Showln("You can't backstab without a blade.")
		return
	}

	defender := FindTarget(attacker, ctx.Raw)
	if defender == nil {
		attacker.Showln("Backstab who?")
		return
	}

	cost := level
	coolDown := level
	InitBattleSkill(attacker, defender, cost, b.Label(), coolDown)

	b.DoBackstab(attacker, defender, level)
}

func (b Backstab) DoBackstab(attacker *world.Character, defender *world.Character, level int) {
	defender.Memory.AddGameEvent(b.Label(), 18)
	hitDamage := simulate.CalculateHitDamage(attacker, defender)
	amount := hitDamage.Amount * 10

	message := world.Message{
		FirstPersonMessage:  fmt.Sprintf("You backstab %s for %d damage.", defender.Name, amount),
		FirstPerson:         attacker,
		SecondPersonMessage: fmt.Sprintf("%s backstabs you for %d damage.", attacker.Name, amount),
		SecondPerson:        defender,
		ThirdPersonMessage:  fmt.Sprintf("%s backstabs %s for %d damage.", attacker.Name, defender.Name, amount),
	}
	if err := attacker.Room.ShowMessage(message); err != nil {
		log.Fatalf("Problem showing backstab message: %v", err)
	}
	simulate.DoDamage(attacker, defender, amount)
}

func (b Backstab) Label() string {
	return "backstab"
}

type Barrier struct{}

func (b Barrier) Execute(ctx Context) {
	player := ctx.Player
	if !player.Skills.Barrier.IsAvailable() {
		player.Showln("You fail to form a barrier.")
		return
	}
	buff := buffs.NewBarrier(player)
	if player.HasBuff(buff.Name()) {
		player.Unapply(buff.Name())
	} else {
		player.Apply(buff)
	}
}

func (b Barrier) Label() string {
	return "barrier"
}

type Bash struct{}

func (b Bash) Execute(ctx Context) {
	attacker := ctx.Player
	level := attacker.Skills.Bash.Value()
	if !CanUseSkill(attacker, b.Label(), level, level) {
		return
	}
	defender := FindTarget(attacker, ctx.Raw)
	if defender == nil {
		attacker.Showln("Bash who?")
		return
	}
	attacker.Spirit.Consume(level)
	attacker.StartAttacking(defender)

	amount := b.CalculateDamage(level)
	attacker.Showln("You bash %s for %d damage.", defender.Name, amount)
	defender.StartAttacking(attacker)
	defender.Showln("%s bashes you for %d damage.", attacker.Name, amount)
	defender.Stun(1)

	simulate.DoDamage(attacker, defender, amount)

	coolDown := buffs.NewCoolDown(9, "bash")
	attacker.ApplyCoolDown(&coolDown)
}

func (b Bash) CalculateDamage(level int) int {
	return 10 + (level * 2)
}

func (b Bash) Label() string {
	return "bash"
}

type Bleed struct{}

func (b Bleed) Execute(ctx Context) {
	attacker := ctx.Player
	level := attacker.Skills.Blitz.Value()
	if !CanUseSkill(attacker, b.Label(), level, level) {
		return
	}

	defender := FindTarget(attacker, ctx.Raw)
	if defender == nil {
		attacker.Showln("Bleed who?")
		return
	}

	InitBattleSkill(attacker, defender, level, b.Label(), level)
	b.DoBleed(ctx.World, attacker, defender, level)
}

func (b Bleed) DoBleed(World *world.World, attacker *world.Character, defender *world.Character, level int) {
	message := world.Message{
		FirstPerson:         attacker,
		SecondPerson:        defender,
		FirstPersonMessage:  fmt.Sprintf("You strike %s precisely, inflicting a grievous wound.", defender.Name),
		SecondPersonMessage: fmt.Sprintf("%s strikes you precisely, inflicting a grievous wound.", attacker.Name),
		ThirdPersonMessage:  fmt.Sprintf("%s stirkes %s precisely, inflicting a grievous wound.", attacker.Name, defender.Name),
	}
	attacker.Room.ShowMessage(message)
	buff := buffs.NewBleed(defender, attacker)
	defender.Apply(buff)
}

func (b Bleed) Label() string {
	return "bleed"
}

type Blitz struct{}

func (b Blitz) Execute(ctx Context) {
	attacker := ctx.Player
	level := attacker.Skills.Blitz.Value()
	if !CanUseSkill(attacker, b.Label(), level, level) {
		return
	}

	defender := FindTarget(attacker, ctx.Raw)
	if defender == nil {
		attacker.Showln("Blitz who?")
		return
	}

	InitBattleSkill(attacker, defender, level, b.Label(), level)
	b.DoBlitz(attacker, defender, level)
}

func (b Blitz) DoBlitz(attacker *world.Character, defender *world.Character, level int) {
	defender.Memory.AddGameEvent(b.Label(), 18)
	for i := 0; i <= level; i++ {
		hitDamage := simulate.CalculateHitDamage(attacker, defender)
		amount := int(float64(10+level) / float64(100) * float64(hitDamage.Amount))

		message := world.Message{
			FirstPersonMessage:  fmt.Sprintf("You blitz %s for %d damage.", defender.Name, amount),
			FirstPerson:         attacker,
			SecondPersonMessage: fmt.Sprintf("%s blitzes you for %d damage.", attacker.Name, amount),
			SecondPerson:        defender,
			ThirdPersonMessage:  fmt.Sprintf("%s blitzes %s for %d damage.", attacker.Name, defender.Name, amount),
		}
		if err := attacker.Room.ShowMessage(message); err != nil {
			log.Fatalf("Problem showing blitz message: %v", err)
		}
		simulate.DoDamage(attacker, defender, amount)
	}
	defender.Stun(1)
}

func (b Blitz) Label() string {
	return "blitz"
}

type Circle struct{}

func (c Circle) Execute(ctx Context) {
	attacker := ctx.Player
	defender := FindTarget(attacker, ctx.Raw)
	if defender == nil {
		attacker.Showln("Circle who?")
		return
	}

	level := attacker.Skills.Circle.Value()
	if !CanUseSkill(attacker, c.Label(), level, level) {
		return
	}

	attacker.Spirit.Consume(level)
	attacker.StartAttacking(defender)
	defender.StartAttacking(attacker)
	coolDown := buffs.NewCoolDown(12, c.Label())
	attacker.ApplyCoolDown(&coolDown)

	if c.TargetExpectsCircle(defender) {
		c.HandleExpectedCircle(attacker, defender)
		return
	}

	defender.Memory.AddGameEvent(c.Label(), 18)
	hitDamage := simulate.CalculateHitDamage(attacker, defender)
	amount := c.CalculateDamage(level, hitDamage.Amount)

	message := world.Message{
		FirstPersonMessage:  fmt.Sprintf("You circle %s for %d damage.", defender.Name, amount),
		FirstPerson:         attacker,
		SecondPersonMessage: fmt.Sprintf("%s circles you for %d damage.", attacker.Name, amount),
		SecondPerson:        defender,
		ThirdPersonMessage:  fmt.Sprintf("%s circles %s for %d damage.", attacker.Name, defender.Name, amount),
	}
	if err := attacker.Room.ShowMessage(message); err != nil {
		log.Fatalf("Problem showing cirle message: %v", err)
	}
	simulate.DoDamage(attacker, defender, amount)
	defender.Stun(1)
}

func (c Circle) CalculateDamage(level int, hitDamage int) int {
	multiplier := 3.0 + (level / 10.0)
	return multiplier * hitDamage
}

func (c Circle) TargetExpectsCircle(target *world.Character) bool {
	return target.Memory.Occurrences(c.Label()) != 0
}

func (c Circle) HandleExpectedCircle(attacker *world.Character, defender *world.Character) {
	message := world.Message{
		FirstPersonMessage:  fmt.Sprintf("You try to circle %s.  Clearly expecting it, %s dodges easily.", defender.Name, defender.Name),
		FirstPerson:         attacker,
		SecondPersonMessage: fmt.Sprintf("%s tries to circle you.  Expecting it, you dodge easily.", attacker.Name),
		SecondPerson:        defender,
		ThirdPersonMessage:  fmt.Sprintf("%s tries to circle %s.  Clearly expecting it, %s dodges easily.", attacker.Name, defender.Name, defender.Name),
	}
	if err := attacker.Room.ShowMessage(message); err != nil {
		log.Fatalf("Problem showing expected circle message: %v", err)
	}
}

func (c Circle) Label() string {
	return "circle"
}

type Die struct{}

func (d Die) Execute(ctx Context) {
	player := ctx.Player
	amount := player.Health.Current
	message := world.Message{
		FirstPersonMessage: fmt.Sprintf("Focusing inward, you do %d damage to yourself.", amount),
		FirstPerson:        player,
		ThirdPersonMessage: fmt.Sprintf("Focusing inward, %s self inflicts %d damage.", player.Name, amount),
	}
	if err := player.Room.ShowMessage(message); err != nil {
		log.Fatalf("Problem showing die message: %v", err)
	}
	simulate.DoDamage(player, player, amount)
}

func (d Die) Label() string {
	return "die"
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

type Flee struct{}

func (f Flee) Execute(ctx Context) {
	player := ctx.Player
	roomUUID := player.Room.Exits.FirstExit()
	room, ok := ctx.World.Rooms[roomUUID]
	if !ok || room.IsFull() {
		player.Showln("There is no where for you to flee.")
		return
	}
	if err := player.Room.Exit(player); err != nil {
		log.Fatalf("Player %s not in room %s", player.UUID, room.UUID)
	}
	if err := room.Enter(player); err != nil {
		log.Fatalf("Room %s is full.", room.UUID)
	}
	player.Room = room
	player.Showln("You flee!")
	for _, opponent := range player.Attacking {
		player.StopAttacking(opponent)
		opponent.StopAttacking(player)
		opponent.Showln("%s flees from you.", player.Name)
	}
}

func (f Flee) Label() string {
	return "flee"
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

	buff := buffs.NewHaste(player)

	if player.HasBuff(buff.Name()) {
		player.Unapply(buff.Name())
	} else {
		player.Apply(buff)
	}
}

func (h Haste) Label() string {
	return "haste"
}

type Inspect struct{}

func (i Inspect) Execute(ctx Context) {
	player := ctx.Player
	parts := strings.SplitN(ctx.Raw, " ", 2)
	if len(parts) == 1 {
		player.Showln("Inspect what?")
		return
	}
	index := player.Inventory.IndexOfItem(parts[1])
	if index == -1 {
		player.Showln("You don't have a '%s'", parts[1])
		return
	}
	item := player.Inventory.GetItemAtIndex(index)
	player.Showln(item.Description())
}

func (i Inspect) Label() string {
	return "inspect"
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
	p.DebitEssence(cost)
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
	case "barrier":
		if spendEssence(player, skills.Barrier.Base) {
			skills.Barrier.Increment()
			player.Showln("Your mastery of barrier improves.")
		}
	case "bash":
		if spendEssence(player, skills.Bash.Base) {
			skills.Bash.Increment()
			player.Showln("Your mastery of bash improves.")
		}
	case "bleed":
		if spendEssence(player, skills.Bleed.Base) {
			skills.Bleed.Increment()
			player.Showln("Your mastery of bleed improves.")
		}
	case "blitz":
		if spendEssence(player, skills.Blitz.Base) {
			skills.Blitz.Increment()
			player.Showln("Your mastery of blitz improves.")
		}
	case "circle":
		if spendEssence(player, skills.Circle.Base) {
			skills.Circle.Increment()
			player.Showln("Your mastery of circle improves.")
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
	player := ctx.Player
	parts := strings.SplitN(ctx.Raw, " ", 2)
	if len(parts) == 1 {
		player.Showln(ctx.Player.Room.Describe(ctx.Player))
		return
	}

	keyword := parts[1]
	item, err := player.Room.IndexOfItem(keyword)
	if err != nil {
		player.Showln("You don't see '%s'.", keyword)
		return
	}
	player.Showln(item.Description())
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

type Save struct{}

func (s Save) Execute(ctx Context) {
	savePlayers(ctx.World.Players)
}

func (s Save) Label() string {
	return "save"
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

type Sweep struct{}

func (s Sweep) Execute(ctx Context) {
	attacker := ctx.Player
	level := attacker.Skills.Sweep.Value()
	cost := level
	if !CanUseSkill(attacker, s.Label(), level, cost) {
		return
	}

	if !IsWieldingRange(attacker) {
		attacker.Showln("You can't sweep without a ranged weapon.")
		return
	}

	defender := FindTarget(attacker, ctx.Raw)
	if defender == nil {
		attacker.Showln("Sweep who?")
		return
	}

	coolDown := level
	InitBattleSkill(attacker, defender, cost, s.Label(), coolDown)

	s.DoSweep(attacker, level)
}

func (s Sweep) DoSweep(attacker *world.Character, level int) {
	for _, defender := range attacker.Attacking {
		defender.Memory.AddGameEvent(s.Label(), 10)
		hitDamage := simulate.CalculateHitDamage(attacker, defender)
		amount := hitDamage.Amount

		message := world.Message{
			FirstPersonMessage:  fmt.Sprintf("You sweep %s for %d damage.", defender.Name, amount),
			FirstPerson:         attacker,
			SecondPersonMessage: fmt.Sprintf("%s sweeps you for %d damage.", attacker.Name, amount),
			SecondPerson:        defender,
			ThirdPersonMessage:  fmt.Sprintf("%s sweeps %s for %d damage.", attacker.Name, defender.Name, amount),
		}
		if err := attacker.Room.ShowMessage(message); err != nil {
			log.Fatalf("Problem showing sweep message: %v", err)
		}
		simulate.DoDamage(attacker, defender, amount)
	}
}

func (s Sweep) Label() string {
	return "sweep"
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

func MovePlayer(ctx Context, roomUUID string, direction string) {
	player := ctx.Player
	if player.IsFighting() {
		player.Showln("You can't go %s while fighting.", direction)
		return
	}

	if roomUUID == "" {
		player.Showln("You can't go %s", direction)
		return
	}

	room, ok := ctx.World.Rooms[roomUUID]
	if !ok {
		player.Showln("An unseen force prevents you from going %s", direction)
		return
	}

	oldRoom := player.Room
	if err := player.Room.Exit(player); err != nil {
		player.Showln("You're unable to leave %s", player.Room.Name())
		return
	}

	if err := room.Enter(player); err != nil {
		oldRoom.Enter(player)
		player.Showln("You're unable to enter %s", room.Name())
	}
	player.Room = room
	player.Room = room
	player.Showln("You go %s", direction)
	Look{}.Execute(ctx)
}

type East struct{}

func (e East) Execute(ctx Context) {
	roomUUID := ctx.Player.Room.Exits.East
	MovePlayer(ctx, roomUUID, e.Label())
}

func (e East) Label() string {
	return "east"
}

type North struct{}

func (n North) Execute(ctx Context) {
	roomUUID := ctx.Player.Room.Exits.North
	MovePlayer(ctx, roomUUID, n.Label())
}

func (n North) Label() string {
	return "north"
}

type South struct{}

func (s South) Execute(ctx Context) {
	roomUUID := ctx.Player.Room.Exits.South
	MovePlayer(ctx, roomUUID, s.Label())
}

func (s South) Label() string {
	return "south"
}

type West struct{}

func (w West) Execute(ctx Context) {
	roomUUID := ctx.Player.Room.Exits.West
	MovePlayer(ctx, roomUUID, w.Label())
}

func (w West) Label() string {
	return "west"
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

type Help struct{}

func (h Help) Execute(ctx Context) {
	parts := strings.SplitN(ctx.Raw, " ", 2)
	if len(parts) == 1 {
		h.ShowAllCommands(ctx)
	} else {
		h.ShowCommand(ctx, parts[1])
	}
}

func (h Help) ShowCommand(ctx Context, command string) {
	help, found := ctx.Help[command]
	if !found {
		ctx.Player.Showln("No help for '%s'", command)
		return
	}
	ctx.Player.ShowNewline()
	ctx.Player.Showln(help.Describe())
	ctx.Player.ShowNewline()
}

func (h Help) ShowAllCommands(ctx Context) {
	player := ctx.Player
	aliases := make([]string, 0)
	for alias, command := range commands {
		if alias == command.Label() {
			aliases = append(aliases, alias)
		}
	}
	sort.Strings(aliases)
	for _, alias := range aliases {
		player.Showln(alias)
	}
	player.ShowNewline()
}

func (h Help) Label() string {
	return "help"
}

type Executor interface {
	Execute(ctx Context)
	Label() string
}

var commands = buildCommands()

func buildCommands() map[string]Executor {
	commands := []Executor{
		Attack{},
		Affect{},
		Backstab{},
		Barrier{},
		Bash{},
		Bleed{},
		Blitz{},
		Circle{},
		Die{},
		Drop{},
		East{},
		Flee{},
		Gear{},
		Get{},
		Haste{},
		Help{},
		Inspect{},
		Inventory{},
		Invest{},
		Look{},
		Noop{},
		North{},
		Remove{},
		Save{},
		Score{},
		South{},
		Sweep{},
		Typo{},
		Wear{},
		West{},
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

func determinePlayer(name string, w *world.World) (*world.Character, error) {
	player, found := w.Players[name]
	if !found {
		return nil, errors.New("player not found")
	}
	return player, nil
}
