package entities

import (
	"box-boy/internal/config"
	"box-boy/internal/customizer"
)

// Player representa o entregador BoxBoy, seu veículo e estado de gameplay.
type Player struct {
	Custom customizer.Customization

	// Posição no mundo
	// X: Lateral na pista (-config.StreetWidth/2 até config.StreetWidth/2 + calçadas)
	// Y: Distância percorrida na rota (avança com a rolagem)
	// Z: Altura do salto (Bunny-Hop)
	X  float64
	Y  float64
	Z  float64
	Vz float64 // Velocidade vertical do pulo

	Speed       float64 // Velocidade atual para frente
	LateralVel  float64 // Velocidade lateral de manobra
	IsJumping   bool
	IsSkidding  bool
	SkidTimer   float64
	InvulnTimer float64 // Tempo de imunidade após bater

	// Inventário e Status
	Cargo        int // Encomendas restantes no bagageiro
	MaxCargo     int
	Reputation   float64 // 0.0 a 100.0 (determina as estrelas de 1 a 5)
	Combo        int     // Entregas consecutivas perfeitas
	MaxCombo     int
	Score        int
	Successful   int // Total de entregas feitas
	Missed       int // Total de entregas perdidas
	BossesBeaten int // Quantidade de chefes superados

	// Animação
	AnimTimer float64
	AnimFrame int
}

// NewPlayer inicializa o jogador com a customização escolhida.
func NewPlayer(c customizer.Customization) *Player {
	maxCap := config.MaxCargoCapacity
	if c.VehicleType == 2 { // Furgão tem mais capacidade
		maxCap = 35
	} else if c.VehicleType == 1 { // Scooter tem capacidade intermediária
		maxCap = 25
	}

	return &Player{
		Custom:     c,
		X:          0, // Começa no meio da pista
		Y:          0,
		Z:          0,
		Vz:         0,
		Speed:      config.StreetSpeed,
		Cargo:      maxCap,
		MaxCargo:   maxCap,
		Reputation: 100.0, // Começa com reputação máxima (5 estrelas Ouro)
		Combo:      0,
	}
}

// Update processa física, salto (Bunny-Hop) e movimentação com delta time.
func (p *Player) Update(dt float64, moveX float64, accelerate bool, brake bool, jump bool) {
	// 1. Controle de Velocidade para Frente
	targetSpeed := config.StreetSpeed
	if accelerate {
		targetSpeed = config.MaxSpeedBoost
		if p.Custom.VehicleType == 1 { // Scooter é mais rápida no turbo
			targetSpeed *= 1.15
		}
	} else if brake {
		targetSpeed = config.BrakeSpeed
	}

	// Interpolação suave de aceleração
	p.Speed += (targetSpeed - p.Speed) * 5.0 * dt
	p.Y += p.Speed * dt

	// 2. Movimento Lateral
	maxLateral := (config.StreetWidth + config.SidewalkWidth) / 2.0
	p.X += moveX * config.LateralSpeed * dt
	if p.X < -maxLateral {
		p.X = -maxLateral
	} else if p.X > maxLateral {
		p.X = maxLateral
	}

	// 3. Salto (Bunny-Hop da Bicicleta / Desvio Alto)
	if jump && !p.IsJumping {
		p.IsJumping = true
		p.Vz = config.BunnyHopJumpForce
		if p.Custom.VehicleType == 0 { // Bike tem melhor bunny-hop
			p.Vz *= 1.1
		}
	}

	if p.IsJumping {
		p.Z += p.Vz * dt
		p.Vz -= config.Gravity * dt
		if p.Z <= 0 {
			p.Z = 0
			p.Vz = 0
			p.IsJumping = false
		}
	}

	// 4. Temporizadores de estado
	if p.SkidTimer > 0 {
		p.SkidTimer -= dt
		if p.SkidTimer <= 0 {
			p.IsSkidding = false
		}
	}
	if p.InvulnTimer > 0 {
		p.InvulnTimer -= dt
	}

	// 5. Ciclo de Animação das Pernas / Rodas
	p.AnimTimer += dt * (p.Speed / config.StreetSpeed) * 8.0
	if p.AnimTimer >= 1.0 {
		p.AnimTimer -= 1.0
		p.AnimFrame = (p.AnimFrame + 1) % 4
	}
}

// ApplyDamage aplica perda de reputação por bater em obstáculos ou errar pedidos
func (p *Player) ApplyDamage(amount float64) {
	if p.InvulnTimer > 0 {
		return
	}
	p.Reputation -= amount
	if p.Reputation < 0 {
		p.Reputation = 0
	}
	p.Combo = 0
	p.InvulnTimer = 1.0 // 1 segundo de invulnerabilidade piscante
}

// AddReputation restaura reputação por entregas perfeitas e recuperação de chefes
func (p *Player) AddReputation(amount float64) {
	p.Reputation += amount
	if p.Reputation > 100.0 {
		p.Reputation = 100.0
	}
}

// GetStarRating retorna o número de estrelas de reputação (1 a 5)
func (p *Player) GetStarRating() int {
	if p.Reputation >= 85.0 {
		return 5
	} else if p.Reputation >= 65.0 {
		return 4
	} else if p.Reputation >= 45.0 {
		return 3
	} else if p.Reputation >= 20.0 {
		return 2
	}
	return 1
}
