package battle

import (
	"github.com/michaelvmata/path/world"
	"math/rand"
)

func CalculateHitDamage(attacker *world.Player, defender *world.Player) int {
	weapon := attacker.Gear.MainHand

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
	return damage
}

func ApplyDamage(p *world.Player, damage int) {
	p.Health.Current -= damage
	if p.IsDead() {
		p.Health.Current = 1
	}
}

func Simulate(w *world.World) {
	fighting := make(map[string]*world.Player)
	for _, attacker := range w.Players {
		if !attacker.IsFighting() {
			continue
		}
		for _, defender := range attacker.Attacking {
			damage := CalculateHitDamage(attacker, defender)
			ApplyDamage(defender, damage)
			attacker.Show("You do %d %s damage to %s",
				damage,
				attacker.Gear.MainHand.WeaponType,
				defender.Name)
			defender.Show("%s did %d %s damage to you",
				attacker.Name,
				damage,
				attacker.Gear.MainHand.WeaponType)
			attacker.Update(0)
			defender.Update(0)
			fighting[attacker.UUID] = attacker
			fighting[defender.UUID] = defender
			break
		}
	}
	for _, p := range fighting {
		p.ShowPrompt()
	}
}
