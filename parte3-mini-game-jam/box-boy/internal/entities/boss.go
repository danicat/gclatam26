package entities

import "math"

type BossType int

const (
	BossDogPack BossType = iota // Mini-Boss: Cliente Ausente & Matilha de Cães
	BossCrater                  // Mini-Boss: A Grande Panela de Asfalto (Cratera Voraz)
	BossProtest                 // Boss Maior: O Megabloqueio da Manifestação
	BossTornado                 // Boss Maior: O Tornado Metropolitano (SimCity)
	BossBlackFriday             // Mega-Boss Final: O Monstro do Atraso da Black Friday
)

type BossState int

const (
	BossInactive BossState = iota
	BossApproaching
	BossPanic
	BossRecover
	BossDefeated
)

// BossEvent gerencia o ciclo dramático de Pânico e Recuperação dos chefes urbanos.
type BossEvent struct {
	Type               BossType
	Name               string
	PanicDescription   string
	RecoverDescription string
	TriggerY           float64 // Distância em que o evento se inicia
	WorldX             float64
	WorldY             float64

	State           BossState
	PanicDuration   float64
	PanicTimeLeft   float64
	RecoverProgress int
	RecoverTarget   int

	ScreenShake float64 // Intensidade de tremor na tela durante Pânico
	RedFlash    float64 // Intensidade de vinheta vermelha de alarme
	AnimTimer   float64
}

// NewBossList cria os 5 bosses e seus marcos de distância ao longo da rota.
func NewBossList() []*BossEvent {
	return []*BossEvent{
		{
			Type:               BossDogPack,
			Name:               "O CLIENTE FANTASMA & A MATILHA DO PORTÃO",
			PanicDescription:   "PÂNICO: NINGUÉM ATENDE! O PORTÃO ABRIU E A MATILHA ESTÁ SOLTA!",
			RecoverDescription: "RECUPERAÇÃO: BUZINE PARA AFASTAR OS CÃES E ACERTE O SMART LOCKER!",
			TriggerY:           650.0,
			PanicDuration:      9.0,
			RecoverTarget:      3,
			State:              BossInactive,
		},
		{
			Type:               BossCrater,
			Name:               "A GRANDE PANELA DE ASFALTO (CRATERA VORAZ)",
			PanicDescription:   "PÂNICO: O ASFALTO CEDEU! UMA CRATERA GIGANTE ESTÁ ENGOLINDO A PISTA!",
			RecoverDescription: "RECUPERAÇÃO: EXECUTE O BUNNY-HOP NA RAMPA E ESTABILIZE O SOLO COM PACOTES!",
			TriggerY:           1350.0,
			PanicDuration:      10.0,
			RecoverTarget:      3,
			State:              BossInactive,
		},
		{
			Type:               BossProtest,
			Name:               "O MEGABLOQUEIO DA MANIFESTAÇÃO SURPRESA",
			PanicDescription:   "PÂNICO: BARRICADAS E FUMAÇA! A AVENIDA PRINCIPAL ESTÁ TOTALMENTE FECHADA!",
			RecoverDescription: "RECUPERAÇÃO: ENCONTRE O ATALHO NA CALÇADA E LANCE BRINDES NA MULTIDÃO!",
			TriggerY:           2100.0,
			PanicDuration:      11.0,
			RecoverTarget:      4,
			State:              BossInactive,
		},
		{
			Type:               BossTornado,
			Name:               "O TORNADO METROPOLITANO (SIMCITY STYLE)",
			PanicDescription:   "PÂNICO: VÓRTICE DE VENTO EXTREMO! PACOTES ESTÃO SENDO SUGADOS PRO CÉU!",
			RecoverDescription: "RECUPERAÇÃO: ATIVE A TRAVA MAGNÉTICA DAS RODAS E RESGATE OS PACOTES NO AR!",
			TriggerY:           2800.0,
			PanicDuration:      12.0,
			RecoverTarget:      4,
			State:              BossInactive,
		},
		{
			Type:               BossBlackFriday,
			Name:               "O MONSTRO DO ATRASO DA BLACK FRIDAY",
			PanicDescription:   "PÂNICO TOTAL: A CONTAGEM PARA A MEIA-NOITE ESTÁ ACELERADA! O COLOSSO ATACA!",
			RecoverDescription: "RECUPERAÇÃO HEROICA: ENTREGA TURBO! ACERTE OS SENSORES DO MONSTRO COM PACOTES!",
			TriggerY:           3450.0,
			PanicDuration:      14.0,
			RecoverTarget:      5,
			State:              BossInactive,
		},
	}
}

// Update avança a máquina de estados do chefe.
func (b *BossEvent) Update(dt float64, playerY float64) {
	b.AnimTimer += dt

	switch b.State {
	case BossInactive:
		// Se o jogador estiver a 120 unidades do chefe, inicia a aproximação
		if playerY >= b.TriggerY-120.0 && playerY < b.TriggerY {
			b.State = BossApproaching
		}
	case BossApproaching:
		if playerY >= b.TriggerY {
			b.StartPanic()
		}
	case BossPanic:
		b.PanicTimeLeft -= dt
		b.ScreenShake = math.Sin(b.AnimTimer*24.0) * 4.5
		b.RedFlash = math.Abs(math.Sin(b.AnimTimer * 6.0))

		// Ao esgotar o tempo de pânico sem zerar, entra na fase de recuperação heroica
		if b.PanicTimeLeft <= 0 {
			b.State = BossRecover
			b.ScreenShake = 1.0
			b.RedFlash = 0.2
		}
	case BossRecover:
		// Tremor diminui durante a recuperação
		if b.ScreenShake > 0 {
			b.ScreenShake -= dt * 2.0
		}
		if b.RedFlash > 0 {
			b.RedFlash -= dt * 0.5
		}
	case BossDefeated:
		b.ScreenShake = 0
		b.RedFlash = 0
	}
}

// StartPanic ativa a crise imprevista do chefe.
func (b *BossEvent) StartPanic() {
	b.State = BossPanic
	b.PanicTimeLeft = b.PanicDuration
	b.RecoverProgress = 0
	b.ScreenShake = 5.0
	b.RedFlash = 0.8
}

// AddRecoverAction registra uma manobra de sucesso do jogador durante o evento.
func (b *BossEvent) AddRecoverAction() bool {
	if b.State != BossPanic && b.State != BossRecover {
		return false
	}

	b.RecoverProgress++
	// Reduz o pânico progressivamente
	b.ScreenShake = math.Max(0, b.ScreenShake-1.5)
	b.RedFlash = math.Max(0, b.RedFlash-0.2)

	if b.RecoverProgress >= b.RecoverTarget {
		b.State = BossDefeated
		return true // Derrotou o chefe!
	}
	return false
}
