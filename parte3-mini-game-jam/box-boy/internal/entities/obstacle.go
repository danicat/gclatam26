package entities

import "math"

type ObstacleType int

const (
	ObsPothole ObstacleType = iota // Buraco no asfalto (requer Bunny-Hop)
	ObsPuddle                      // Poça d'água (faz escorregar)
	ObsTrafficCone                 // Cone de trânsito
	ObsRoadBarrier                 // Barricada de obras
	ObsBarkingDog                  // Cão latindo no portão / rua
	ObsSprinkler                   // Aspersor de água
)

// Obstacle representa um perigo urbano na pista ou calçada.
type Obstacle struct {
	ID       int
	Type     ObstacleType
	WorldX   float64
	WorldY   float64
	Width    float64
	Height   float64
	Radius   float64
	IsScared bool    // Cão assustado pela buzina
	Hit      bool    // Já atingiu o jogador
	AnimTime float64
}

// NewObstacle cria um obstáculo configurado com raio de colisão.
func NewObstacle(id int, oType ObstacleType, wx, wy float64) *Obstacle {
	radius := 14.0
	switch oType {
	case ObsPothole:
		radius = 16.0
	case ObsPuddle:
		radius = 18.0
	case ObsRoadBarrier:
		radius = 20.0
	case ObsBarkingDog:
		radius = 15.0
	}

	return &Obstacle{
		ID:     id,
		Type:   oType,
		WorldX: wx,
		WorldY: wy,
		Radius: radius,
	}
}

// CheckCollision verifica sobreposição com o jogador considerando altura Z do pulo.
func (o *Obstacle) CheckCollision(px, py, pz float64) bool {
	if o.Hit || o.IsScared {
		return false
	}

	// Se o jogador estiver alto o bastante no Bunny-Hop (Z > 12), passa por cima de buracos e cones!
	if pz > 12.0 && (o.Type == ObsPothole || o.Type == ObsTrafficCone || o.Type == ObsPuddle) {
		return false
	}

	dist := math.Hypot(px-o.WorldX, py-o.WorldY)
	return dist < (o.Radius + 10.0)
}

// ScareDog afasta o cão com a buzina da bicicleta/scooter.
func (o *Obstacle) ScareDog() {
	if o.Type == ObsBarkingDog {
		o.IsScared = true
	}
}
