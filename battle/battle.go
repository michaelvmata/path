package battle

import (
	"github.com/michaelvmata/path/buffs"
	"github.com/michaelvmata/path/events"
	"github.com/michaelvmata/path/world"
	"math/rand"
	"strings"
)

func CalculateHitDamage(attacker *world.Character, defender *world.Character) int {
	weapon := attacker.Weapon()

	// Get Weapon base damage
	damage := rand.Intn(weapon.MaximumDamage-weapon.MinimumDamage) + weapon.MinimumDamage

	// Check if a critical hit
	agilityDiff := float64(attacker.Core.Agility.Value() - defender.Core.Agility.Value())
	criticalRate := weapon.CriticalRate + (agilityDiff * .01)
	if rand.Float64() <= criticalRate {
		// Apply critical bonus
		insightDiff := float64(attacker.Core.Insight.Value() - defender.Core.Insight.Value())
		criticalBonus := weapon.CriticalBonus + (insightDiff * .01)
		damage = int(float64(damage) * (1.0 + criticalBonus))
	}

	return damage + (10 * attacker.Core.Power.Value())
}

func DoAttack(world *world.World, attacker *world.Character, defender *world.Character) {
	damage := CalculateHitDamage(attacker, defender)
	defender.Health.Current -= damage
	attacker.Showln("You do %d %s damage to %s.",
		damage,
		strings.ToLower(attacker.Weapon().WeaponType),
		defender.Name)
	defender.Showln("%s does %d %s damage to you.",
		attacker.Name,
		damage,
		strings.ToLower(attacker.Weapon().WeaponType))
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
		if buff.Name() == buffs.HasteName {
			attackNumber++
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
			DoAttack(w, attacker, defender)
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
