package internal

import "math/rand/v2"

type Member struct {
	alive            bool
	sick             bool
	isBuyingMedicine bool
	name             string
	health           int
	hunger           int
	warmth           int
}

func (m *Member) doUpdate(food, warmth bool) bool {
	if !m.alive {
		return false
	}
	if food {
		m.hunger += 2
	} else {
		m.hunger--
	}
	if warmth {
		m.warmth += 2
	} else {
		m.warmth--
	}

	if !food && m.warmth <= 7 {
		m.hunger--
	}

	if !food && !warmth && (m.hunger <= 5 || m.warmth <= 5) {
		m.health--

		if !m.sick {
			m.sick = doSickRoll(m.hunger + m.warmth)
		}
	}
	if m.isBuyingMedicine {
		m.sick = false
		m.isBuyingMedicine = false
	}
	if m.sick {
		m.health--
	}

	if m.health <= 0 {
		m.alive = false
	}
	if m.alive && food && warmth && !m.sick {
		m.health++
	}
	return m.alive
}

func doSickRoll(health int) bool {
	r := rand.IntN(20)

	return r > health
}
