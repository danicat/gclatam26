package entities

import (
	"box-boy/internal/config"
)

type HouseStatus int

const (
	StatusPending HouseStatus = iota
	StatusDelivered
	StatusMissed
	StatusDamaged
)

// House representa um ponto de entrega ao longo do percurso (residência ou Smart Locker).
type House struct {
	ID         int
	Side       int     // -1: Calçada Esquerda, +1: Calçada Direita
	WorldX     float64 // Coordenada X no mundo
	WorldY     float64 // Distância na rota
	Style      int     // 0..2: Casas residenciais, 3: Smart Locker
	TargetType int     // 0: Porta/Varanda, 1: Caixa de Correio, 2: Compartimento Locker

	TargetX float64 // Posição exata do alvo de arremesso
	TargetY float64

	Status        HouseStatus
	DeliveredTime float64
	CustomerHappy bool
}

// NewHouse instancia um novo ponto de entrega na calçada especificada.
func NewHouse(id int, side int, worldY float64, isLocker bool) *House {
	style := (id % 3)
	targetType := 0
	if isLocker {
		style = 3
		targetType = 2
	} else if id%4 == 1 {
		targetType = 1 // Caixa de correio
	}

	// Posição na calçada lateral
	var posX float64
	if side < 0 {
		posX = -(config.StreetWidth/2.0 + 55.0)
	} else {
		posX = (config.StreetWidth/2.0 + 55.0)
	}

	// Alvo de arremesso (tapete de entrada, correio ou locker)
	tX := posX
	tY := worldY
	if side < 0 {
		tX += 14.0 // Mais próximo da rua
	} else {
		tX -= 14.0
	}

	return &House{
		ID:         id,
		Side:       side,
		WorldX:     posX,
		WorldY:     worldY,
		Style:      style,
		TargetType: targetType,
		TargetX:    tX,
		TargetY:    tY,
		Status:     StatusPending,
	}
}
