package config

const (
	// Resolução Virtual (16:9 Retro-HD)
	VirtualWidth  = 640
	VirtualHeight = 360

	// Configuração de Janela
	WindowWidth  = 1280
	WindowHeight = 720
	WindowTitle  = "BoxBoy: Turbo Express (Entrega Turbo)"

	// Física e Movimento Isométrico
	TargetFPS      = 60
	IsoTileRatio   = 0.5 // Razão isométrica clássica 2:1
	StreetSpeed    = 130.0 // Velocidade base de rolagem da rota em px/s
	MaxSpeedBoost  = 220.0
	BrakeSpeed     = 70.0
	LateralSpeed   = 160.0 // Velocidade de manobra lateral entre pista e calçada

	// Dimensões do Mundo
	StreetWidth    = 180.0 // Largura da rua em coordenadas locais
	SidewalkWidth  = 140.0 // Largura das calçadas
	RouteLength    = 3600.0 // Comprimento total da rota da fase em unidades de mundo

	// Arremesso de Pacotes
	PackageThrowSpeed = 320.0
	PackageGravity    = 520.0
	MaxCargoCapacity  = 20

	// Bunny Hop (Pulo da Bicicleta)
	BunnyHopJumpForce = 210.0
	Gravity           = 580.0
)
