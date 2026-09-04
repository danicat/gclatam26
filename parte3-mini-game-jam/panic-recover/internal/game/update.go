package game

import (
	"math"
	"time"
)

const calmMoveSpeed = 45.0

const (
	panicStartSpeed = 78.0
	panicEndSpeed   = 120.0
)

type Input struct {
	Move           Vec2
	PanicPressed   bool
	StartPressed   bool
	RestartPressed bool
}

func (m *Model) Update(input Input, dt time.Duration) {
	if m.Scene == SceneTitle {
		if input.StartPressed {
			m.Start()
		}
		return
	}
	if m.Scene == SceneVictory || m.Scene == SceneGameOver {
		if input.RestartPressed {
			m.Start()
		}
		return
	}
	m.updatePlayer(input.Move, dt)
	m.updateBugs(dt)
	if m.Phase == PhaseCalm {
		if m.playerTouchesLiveBug() {
			m.beginPanic(true)
			return
		}
		if input.PanicPressed {
			m.beginPanic(false)
			return
		}
		m.CalmRemaining -= dt
		if m.CalmRemaining <= 0 {
			m.CalmRemaining = 0
			m.beginPanic(false)
		}
		return
	}
	if m.Phase == PhasePanic || m.Phase == PhaseRecoverAvailable {
		m.eliminateTouchedBugs()
		if m.Phase == PhasePanic && m.Eliminations >= m.Config.Cycles[m.Cycle].Quota {
			m.activateRecover()
		}
		if m.Phase == PhaseRecoverAvailable && m.Recover.Active && overlaps(m.Player.Position, m.Player.Radius, m.Recover.Position, m.Recover.Radius) {
			m.completeRecover()
			return
		}
		m.PanicRemaining -= dt
		if m.PanicRemaining <= 0 {
			m.PanicRemaining = 0
			m.Scene = SceneGameOver
		}
	}
}

func (m *Model) completeRecover() {
	m.Recover.Active = false
	if m.Cycle == len(m.Config.Cycles)-1 {
		m.Scene = SceneVictory
		return
	}
	m.Cycle++
	m.startCycle()
}

func (m *Model) eliminateTouchedBugs() {
	for i := range m.Bugs {
		bug := &m.Bugs[i]
		if bug.Alive && overlaps(m.Player.Position, m.Player.Radius, bug.Position, bug.Radius) {
			bug.Alive = false
			m.Eliminations++
		}
	}
}

func (m *Model) activateRecover() {
	candidates := [...]Vec2{
		{X: 24, Y: 24}, {X: VirtualWidth / 2, Y: 24}, {X: VirtualWidth - 24, Y: 24},
		{X: VirtualWidth - 24, Y: VirtualHeight / 2}, {X: VirtualWidth - 24, Y: VirtualHeight - 24},
		{X: VirtualWidth / 2, Y: VirtualHeight - 24}, {X: 24, Y: VirtualHeight - 24},
		{X: 24, Y: VirtualHeight / 2},
	}
	best := -1
	bestDistance := -1.0
	for i, candidate := range candidates {
		if overlaps(candidate, m.Recover.Radius, m.Player.Position, m.Player.Radius) || m.candidateTouchesLiveBug(candidate) {
			continue
		}
		if distance := distanceSquared(candidate, m.Player.Position); distance > bestDistance {
			best = i
			bestDistance = distance
		}
	}
	if best == -1 {
		for i, candidate := range candidates {
			if overlaps(candidate, m.Recover.Radius, m.Player.Position, m.Player.Radius) {
				continue
			}
			if distance := distanceSquared(candidate, m.Player.Position); distance > bestDistance {
				best = i
				bestDistance = distance
			}
		}
		if best == -1 {
			best = 0
		}
		for i := range m.Bugs {
			if m.Bugs[i].Alive && overlaps(candidates[best], m.Recover.Radius, m.Bugs[i].Position, m.Bugs[i].Radius) {
				m.Bugs[i].Alive = false
			}
		}
	}
	m.Recover.Position = candidates[best]
	m.Recover.Active = true
	m.Phase = PhaseRecoverAvailable
}

func (m *Model) candidateTouchesLiveBug(candidate Vec2) bool {
	for i := range m.Bugs {
		if m.Bugs[i].Alive && overlaps(candidate, m.Recover.Radius, m.Bugs[i].Position, m.Bugs[i].Radius) {
			return true
		}
	}
	return false
}

func distanceSquared(a, b Vec2) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return dx*dx + dy*dy
}

func (m *Model) playerTouchesLiveBug() bool {
	for i := range m.Bugs {
		if m.Bugs[i].Alive && overlaps(m.Player.Position, m.Player.Radius, m.Bugs[i].Position, m.Bugs[i].Radius) {
			return true
		}
	}
	return false
}

func overlaps(a Vec2, aRadius float64, b Vec2, bRadius float64) bool {
	dx := a.X - b.X
	dy := a.Y - b.Y
	radius := aRadius + bRadius
	return dx*dx+dy*dy <= radius*radius
}

func (m *Model) updatePlayer(direction Vec2, dt time.Duration) {
	length := math.Hypot(direction.X, direction.Y)
	if length > 1 {
		direction.X /= length
		direction.Y /= length
	}
	seconds := dt.Seconds()
	if m.Phase == PhaseCalm {
		m.Player.Velocity = Vec2{
			X: direction.X * calmMoveSpeed,
			Y: direction.Y * calmMoveSpeed,
		}
	} else {
		progress := 1 - float64(m.PanicRemaining)/float64(m.Config.PanicDuration)
		progress = clamp(progress, 0, 1)
		speed := panicStartSpeed + (panicEndSpeed-panicStartSpeed)*progress
		responsiveness := 10 - 7*progress
		blend := clamp(responsiveness*seconds, 0, 1)
		targetX := direction.X * speed
		targetY := direction.Y * speed
		m.Player.Velocity.X += (targetX - m.Player.Velocity.X) * blend
		m.Player.Velocity.Y += (targetY - m.Player.Velocity.Y) * blend
	}
	m.Player.Position.X += m.Player.Velocity.X * seconds
	m.Player.Position.Y += m.Player.Velocity.Y * seconds
	m.Player.Position.X = clamp(m.Player.Position.X, m.Player.Radius, VirtualWidth-m.Player.Radius)
	m.Player.Position.Y = clamp(m.Player.Position.Y, m.Player.Radius, VirtualHeight-m.Player.Radius)
}

func (m *Model) updateBugs(dt time.Duration) {
	seconds := dt.Seconds()
	for i := range m.Bugs {
		bug := &m.Bugs[i]
		if !bug.Alive {
			continue
		}
		dx := m.Player.Position.X - bug.Position.X
		dy := m.Player.Position.Y - bug.Position.Y
		length := math.Hypot(dx, dy)
		if length == 0 {
			continue
		}
		bug.Position.X += dx / length * bug.Speed * seconds
		bug.Position.Y += dy / length * bug.Speed * seconds
	}
}

func clamp(value, minimum, maximum float64) float64 {
	return math.Max(minimum, math.Min(maximum, value))
}

func (m *Model) beginPanic(forced bool) {
	m.Phase = PhasePanic
	m.PanicRemaining = m.Config.PanicDuration
	if forced {
		m.PanicRemaining = time.Duration(math.Round(float64(m.PanicRemaining) * m.Config.ForcedPanicFraction))
	}
}
