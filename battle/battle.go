package battle

import (
	"github.com/michaelvmata/path/buffs"
	"github.com/michaelvmata/path/events"
	"github.com/michaelvmata/path/world"
	"math/rand"
	"strings"
)

type Damage struct {
	Amount   int
	Type     string
	Critical bool
}

func CalculateHitDamage(attacker *world.Character, defender *world.Character) Damage {
	weapon := attacker.Weapon()
	damage := Damage{
		Type:     strings.ToLower(weapon.WeaponType),
		Amount:   rand.Intn(weapon.MaximumDamage-weapon.MinimumDamage) + weapon.MinimumDamage,
		Critical: false,
	}

	// Check if a critical hit
	agilityDiff := float64(attacker.Core.Agility.Value() - defender.Core.Agility.Value())
	criticalRate := weapon.CriticalRate + (agilityDiff * .01)
	if rand.Float64() <= criticalRate {
		// Apply critical bonus
		insightDiff := float64(attacker.Core.Insight.Value() - defender.Core.Insight.Value())
		criticalBonus := weapon.CriticalBonus + (insightDiff * .01)
		damage.Amount = int(float64(damage.Amount) * (1.0 + criticalBonus))
		damage.Critical = true
	}

	damage.Amount = damage.Amount + (10 * attacker.Core.Power.Value())
	return damage
}

func ShouldEvade(defender *world.Character) bool {
	evasionLevel := defender.Skills.Evasion.Value()
	if evasionLevel <= 0 {
		return false
	}
	evasionRate := float64(evasionLevel)*.01 + float64(defender.Core.Agility.Value())*.01
	return rand.Float64() <= evasionRate
}

func DoEvade(attacker *world.Character, defender *world.Character) {
	attacker.Showln("%s evades you.", defender.Name)
	defender.Showln("You evade %s.", attacker.Name)
}

func DoAttack(world *world.World, attacker *world.Character, defender *world.Character) {
	damage := CalculateHitDamage(attacker, defender)
	defender.Health.Current -= damage.Amount
	highlight := "white"
	if damage.Critical {
		highlight = "orange_3"
	}

	attacker.Showln("You do <%s>%d<reset> %s damage to %s.",
		highlight,
		damage.Amount,
		damage.Type,
		defender.Name)
	defender.Showln("%s does <%s>%d<reset> %s damage to you.",
		attacker.Name,
		highlight,
		damage.Amount,
		damage.Type)

	if defender.IsDead() {
		events.CharacterDeath.Emit(events.CharacterDeathPayload{
			Character: defender,
			Killer:    attacker,
			World:     world,
		})
	}
	attacker.Update(0)
	defender.Update(0)
}

func NumberOfAttacks(character *world.Character) int {
	attackNumber := 1
	for _, buff := range character.Buffs {
		if buffHaste, ok := buff.(*buffs.Haste); ok && !buff.IsExpired() {
			upkeep := buff.Upkeep()
			if character.Spirit.IsAvailable(upkeep) {
				character.Spirit.Consume(upkeep)
				attackNumber += buffHaste.NumberOfAttacks()
			} else {
				buff.Expire()
			}
		}
	}
	return attackNumber
}

func Round(w *world.World, fighting map[string]*world.Character, attacker *world.Character) {
	if !attacker.IsFighting() {
		return
	}
	if _, ok := fighting[attacker.UUID]; !ok {
		attacker.Showln("")
	}

	for _, defender := range attacker.Attacking {
		if _, ok := fighting[defender.UUID]; !ok {
			defender.Showln("")
		}
		for i := 1; i <= NumberOfAttacks(attacker); i++ {
			if ShouldEvade(defender) {
				DoEvade(attacker, defender)
			} else {
				DoAttack(w, attacker, defender)
			}
		}
		fighting[attacker.UUID] = attacker
		fighting[defender.UUID] = defender
		break
	}
}

func Simulate(w *world.World) {
	fighting := make(map[string]*world.Character)
	for _, attacker := range w.Players {
		Round(w, fighting, attacker)
	}
	for _, attacker := range w.Mobiles.Instances {
		Round(w, fighting, attacker)
	}
	for _, p := range fighting {
		p.ShowPrompt()
	}
}
