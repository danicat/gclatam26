package entities_test

import (
	"testing"

	"box-boy/internal/customizer"
	"box-boy/internal/entities"
)

func TestPlayerMovementAndBunnyHop(t *testing.T) {
	c := customizer.NewDefaultCustomization()
	p := entities.NewPlayer(c)

	initialY := p.Y
	// 1 frame a 60fps com aceleração
	p.Update(1.0/60.0, 1.0, true, false, false)

	if p.Y <= initialY {
		t.Errorf("jogador deveria ter avançado em Y, y=%f", p.Y)
	}
	if p.X <= 0 {
		t.Errorf("jogador deveria ter se movido lateralmente para a direita (+X), x=%f", p.X)
	}

	// Teste do Bunny-Hop
	p.Update(1.0/60.0, 0, false, false, true)
	if !p.IsJumping || p.Z <= 0 {
		t.Errorf("salto (Bunny-Hop) deveria estar ativo com Z > 0, isJumping=%v, z=%f", p.IsJumping, p.Z)
	}
}

func TestPlayerDamageAndReputation(t *testing.T) {
	c := customizer.NewDefaultCustomization()
	p := entities.NewPlayer(c)

	if p.GetStarRating() != 5 {
		t.Errorf("esperava 5 estrelas iniciais, obteve %d", p.GetStarRating())
	}

	p.ApplyDamage(40.0)
	if p.Reputation != 60.0 {
		t.Errorf("esperava reputação 60.0 após dano, obteve %f", p.Reputation)
	}
	if p.GetStarRating() != 3 {
		t.Errorf("esperava 3 estrelas com 60%% de reputação, obteve %d", p.GetStarRating())
	}

	p.AddReputation(30.0)
	if p.Reputation != 90.0 {
		t.Errorf("esperava reputação 90.0 após bônus, obteve %f", p.Reputation)
	}
	if p.GetStarRating() != 5 {
		t.Errorf("esperava 5 estrelas restauradas, obteve %d", p.GetStarRating())
	}
}
