package internal

import "math/rand/v2"

type Member struct {
	alive  bool
	sick   bool
	name   string
	health int
	hunger int
	warmth int
}

func (m *Member) doUpdate() bool {
	if !m.alive {
		return false
	}
	m.hunger--
	m.warmth--

	if m.warmth <= 7 {
		m.hunger--
	}

	if m.hunger <= 5 || m.warmth <= 5 {
		m.health--

		if !m.sick {
			m.sick = doSickRoll(m.hunger + m.warmth)
		}
	}
	if !m.sick {
		m.health--
	}

	if m.health <= 0 {
		m.alive = false
	}
	return m.alive
}

func doSickRoll(health int) bool {
	r := rand.IntN(20)

	return r > health
}
