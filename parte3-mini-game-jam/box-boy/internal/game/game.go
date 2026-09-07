package game

import (
	"image/color"
	"math"
	"math/rand"
	"sort"

	"box-boy/internal/art"
	"box-boy/internal/audio"
	"box-boy/internal/config"
	"box-boy/internal/customizer"
	"box-boy/internal/entities"
	"box-boy/internal/input"
	"box-boy/internal/render"
	"box-boy/internal/ui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GameState int

const (
	StateTitle GameState = iota
	StateCustomizer
	StatePlaying
	StateGameOver
	StateVictory
)

// RenderItem define um elemento gráfico com chave de profundidade para Y-Sorting.
type RenderItem struct {
	Depth float64
	Draw  func(screen *ebiten.Image)
}

// Game implementa a interface ebiten.Game.
type Game struct {
	state    GameState
	atlas    *art.TextureAtlas
	audio    *audio.AudioSystem
	custom   customizer.Customization
	player   *entities.Player
	customCat int // Categoria selecionada na tela de customização

	// Mundo e Entidades
	houses    []*entities.House
	obstacles []*entities.Obstacle
	packages  []*entities.PackageProjectile
	bosses    []*entities.BossEvent
	particles *entities.ParticleSystem

	// Controle de Câmera e Tremor
	screenShake float64
	redFlash    float64
	titleTimer  float64

	// Opções de renderização reutilizáveis (Zero-Allocation)
	drawOpts ebiten.DrawImageOptions
}

// NewGame inicializa o jogo, atlas procedural e sintetizador de som.
func NewGame() *Game {
	g := &Game{
		state:     StateTitle,
		atlas:     art.NewTextureAtlas(),
		audio:     audio.NewAudioSystem(),
		custom:    customizer.NewDefaultCustomization(),
		particles: entities.NewParticleSystem(300),
	}

	g.audio.PlayBGM("groove")
	return g
}

// startNewMatch inicializa a rota de entrega, casas, perigos e chefes.
func (g *Game) startNewMatch() {
	g.player = entities.NewPlayer(g.custom)
	g.houses = nil
	g.obstacles = nil
	g.packages = nil
	g.bosses = entities.NewBossList()

	// Gerar casas ao longo da rota de 3600 unidades
	houseID := 0
	for y := 200.0; y < config.RouteLength-150.0; y += 120.0 {
		// Alterna entre calçada esquerda e direita
		side := -1
		if houseID%2 == 1 {
			side = 1
		}
		isLocker := (houseID%5 == 4) // Cada 5ª entrega é um Smart Locker!
		g.houses = append(g.houses, entities.NewHouse(houseID, side, y, isLocker))
		houseID++
	}

	// Gerar obstáculos na pista (buracos, cones, cães, poças)
	obsID := 0
	for y := 300.0; y < config.RouteLength-200.0; y += 95.0 {
		oType := entities.ObstacleType(obsID % 6)
		// Posição na pista ou calçada
		laneOffset := -40.0 + rand.Float64()*80.0
		if oType == entities.ObsBarkingDog {
			// Cães ficam perto dos portões das calçadas
			if obsID%2 == 0 {
				laneOffset = -75.0
			} else {
				laneOffset = 75.0
			}
		}
		g.obstacles = append(g.obstacles, entities.NewObstacle(obsID, oType, laneOffset, y))
		obsID++
	}

	g.state = StatePlaying
	g.audio.PlayBGM("groove")
}

// Layout define a resolução virtual fixa de 640x360 (16:9).
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.VirtualWidth, config.VirtualHeight
}

// Update executa a lógica de estado do jogo a 60 FPS.
func (g *Game) Update() error {
	dt := 1.0 / 60.0
	in := input.PollInputs()

	switch g.state {
	case StateTitle:
		g.titleTimer += dt
		if in.JustSelected {
			g.audio.PlaySFX("click")
			g.state = StateCustomizer
		}

	case StateCustomizer:
		// Navegar categorias verticais
		if in.NavY < 0 {
			g.customCat = (g.customCat - 1 + 8) % 8
			g.audio.PlaySFX("click")
		} else if in.NavY > 0 {
			g.customCat = (g.customCat + 1) % 8
			g.audio.PlaySFX("click")
		}

		// Trocar opção horizontal da categoria
		if in.NavX < 0 {
			g.custom.PrevOption(g.customCat)
			g.audio.PlaySFX("click")
		} else if in.NavX > 0 {
			g.custom.NextOption(g.customCat)
			g.audio.PlaySFX("click")
		}

		// Confirmar e iniciar a rota
		if in.JustSelected {
			g.audio.PlaySFX("combo")
			g.startNewMatch()
		}

	case StatePlaying:
		g.updateGameplay(dt, in)

	case StateGameOver:
		if in.JustSelected {
			g.audio.PlaySFX("click")
			g.state = StateTitle
			g.audio.PlayBGM("groove")
		}

	case StateVictory:
		if in.JustSelected {
			g.audio.PlaySFX("click")
			g.state = StateTitle
			g.audio.PlayBGM("groove")
		}
	}

	return nil
}

func (g *Game) updateGameplay(dt float64, in input.InputState) {
	p := g.player

	// 1. Atualizar física e movimento do jogador
	p.Update(dt, in.MoveX, in.Accelerate, in.Brake, in.JustJumped)

	// Som de salto (Bunny-Hop)
	if in.JustJumped {
		g.audio.PlaySFX("bunnyhop")
		g.particles.SpawnDust(p.X, p.Y, 0)
	}

	// Som de buzina / afastar cães
	if in.JustHorn {
		if p.Custom.VehicleType == 0 {
			g.audio.PlaySFX("bell")
		} else {
			g.audio.PlaySFX("horn")
		}

		// Assustar cães próximos
		for _, obs := range g.obstacles {
			if obs.Type == entities.ObsBarkingDog && !obs.IsScared {
				dist := math.Hypot(obs.WorldX-p.X, obs.WorldY-p.Y)
				if dist < 120.0 {
					obs.ScareDog()
					p.Score += 30
				}
			}
		}
	}

	// Emitir poeira das rodas
	if rand.Float64() < 0.25 {
		g.particles.SpawnDust(p.X, p.Y, 0)
	}

	// 2. Arremesso de Pacote
	if in.JustThrew && p.Cargo > 0 {
		// Encontrar a casa ou locker mais próximo à frente
		var targetHouse *entities.House
		minDist := 260.0
		for _, h := range g.houses {
			if h.Status == entities.StatusPending && h.WorldY >= p.Y-20.0 && h.WorldY <= p.Y+minDist {
				targetHouse = h
				break
			}
		}

		// Se achou uma casa alvo, arremessa nela; se não, arremessa na calçada lateral
		tx := p.X - 60.0
		ty := p.Y + 80.0
		pkgType := 0
		if targetHouse != nil {
			tx = targetHouse.TargetX
			ty = targetHouse.TargetY
			if targetHouse.Style == 3 {
				pkgType = 2 // Pacote grande no locker
			}
		}

		proj := entities.NewPackageProjectile(p.X, p.Y, p.Z, tx, ty, pkgType, targetHouse)
		g.packages = append(g.packages, proj)
		p.Cargo--
		g.audio.PlaySFX("throw")
	}

	// 3. Atualizar Pacotes em voo
	for i := len(g.packages) - 1; i >= 0; i-- {
		pkg := g.packages[i]
		landed := pkg.Update(dt)
		if landed {
			if pkg.TargetHouse != nil && pkg.TargetHouse.Status == entities.StatusPending {
				distToTarget := math.Hypot(pkg.CurrentX-pkg.TargetHouse.TargetX, pkg.CurrentY-pkg.TargetHouse.TargetY)
				if distToTarget < 35.0 {
					// ENTREGA PERFEITA!
					pkg.TargetHouse.Status = entities.StatusDelivered
					pkg.TargetHouse.CustomerHappy = true
					p.Successful++
					p.Combo++
					if p.Combo > p.MaxCombo {
						p.MaxCombo = p.Combo
					}
					points := 100 * p.Combo
					p.Score += points
					p.AddReputation(4.0)

					// Feedback sonoro e visual
					if p.Combo >= 3 {
						g.audio.PlaySFX("combo")
					} else {
						g.audio.PlaySFX("deliver")
					}
					g.particles.BurstStars(pkg.CurrentX, pkg.CurrentY, 15.0, 16)
				} else {
					// Errou o alvo
					pkg.TargetHouse.Status = entities.StatusMissed
					p.Missed++
					p.ApplyDamage(6.0)
					g.audio.PlaySFX("crash")
				}
			}
			// Remove o pacote que pousou
			g.packages = append(g.packages[:i], g.packages[i+1:]...)
		}
	}

	// 4. Checar Colisões com Obstáculos
	for _, obs := range g.obstacles {
		if obs.CheckCollision(p.X, p.Y, p.Z) {
			obs.Hit = true
			p.ApplyDamage(12.0)
			g.screenShake = 4.0
			g.audio.PlaySFX("crash")

			if obs.Type == entities.ObsBarkingDog {
				g.audio.PlaySFX("bark")
			}
		}
	}

	// 5. Atualizar Chefes e Eventos de Pânico / Recuperação
	var currentActiveBoss *entities.BossEvent
	for _, b := range g.bosses {
		b.Update(dt, p.Y)
		if b.State == entities.BossPanic || b.State == entities.BossRecover {
			currentActiveBoss = b
			g.screenShake = b.ScreenShake
			g.redFlash = b.RedFlash

			// Durante o pânico, toca a música de tensão
			g.audio.PlayBGM("panic")

			// Ação heroica de recuperação acionada
			if in.JustRecovered || (b.State == entities.BossRecover && (in.JustThrew || in.JustJumped || in.JustHorn)) {
				defeated := b.AddRecoverAction()
				g.audio.PlaySFX("recover")
				g.particles.BurstStars(p.X, p.Y, 20.0, 20)

				if defeated {
					p.BossesBeaten++
					p.Score += 800
					p.AddReputation(20.0)
					g.audio.PlaySFX("combo")
					g.audio.PlayBGM("groove")
				}
			}
		}
	}

	// Se não houver chefe ativo, suaviza o tremor
	if currentActiveBoss == nil {
		if g.screenShake > 0 {
			g.screenShake -= dt * 10.0
			if g.screenShake < 0 {
				g.screenShake = 0
			}
		}
		g.redFlash = 0
	}

	// 6. Atualizar Partículas
	g.particles.Update(dt)

	// 7. Condições de Vitória e Derrota
	if p.Reputation <= 0 {
		g.state = StateGameOver
		g.audio.PlaySFX("crash")
		g.audio.StopBGM()
	} else if p.Y >= config.RouteLength {
		g.state = StateVictory
		g.audio.PlaySFX("recover")
		g.audio.PlayBGM("victory")
	}
}

// Draw renderiza o jogo com projeção isométrica e ordenação por profundidade Y.
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case StateTitle:
		g.drawTitleScreen(screen)

	case StateCustomizer:
		ui.DrawCustomizerUI(screen, &g.custom, g.customCat)
		// Renderiza o preview do entregador montado na bicicleta
		g.drawCharacterPreview(screen)

	case StatePlaying:
		g.drawGameplay(screen)

	case StateGameOver:
		g.drawGameplay(screen)
		ui.DrawGameOverUI(screen, g.player)

	case StateVictory:
		g.drawGameplay(screen)
		ui.DrawVictoryUI(screen, g.player)
	}
}

func (g *Game) drawTitleScreen(screen *ebiten.Image) {
	// Céu noturno com gradiente
	vector.DrawFilledRect(screen, 0, 0, float32(config.VirtualWidth), float32(config.VirtualHeight), color.RGBA{18, 22, 35, 255}, false)

	// Faixa decorativa amarela
	vector.DrawFilledRect(screen, 0, 90, float32(config.VirtualWidth), 110, color.RGBA{255, 230, 0, 255}, false)
	vector.DrawFilledRect(screen, 0, 86, float32(config.VirtualWidth), 4, color.RGBA{30, 90, 220, 255}, false)
	vector.DrawFilledRect(screen, 0, 200, float32(config.VirtualWidth), 4, color.RGBA{30, 90, 220, 255}, false)

	ebitenutil.DebugPrintAt(screen, "=========================================================", 110, 100)
	ebitenutil.DebugPrintAt(screen, "            B O X B O Y :  T U R B O   E X P R E S S", 100, 120)
	ebitenutil.DebugPrintAt(screen, "         A R C A D E   P A P E R B O Y   E D I T I O N", 110, 140)
	ebitenutil.DebugPrintAt(screen, "=========================================================", 110, 160)

	// Animação de bicicleta passando na faixa
	bikeX := math.Mod(g.titleTimer*160.0, float64(config.VirtualWidth)+80.0) - 40.0
	g.drawOpts.GeoM.Reset()
	g.drawOpts.GeoM.Translate(bikeX, 115)
	screen.DrawImage(g.atlas.BicycleFrame, &g.drawOpts)

	// Texto piscante
	if int(g.titleTimer*2)%2 == 0 {
		ebitenutil.DebugPrintAt(screen, ">> PRESSIONE [ESPACO] OU [ENTER] PARA CUSTOMIZAR E JOGAR <<", 105, 240)
	}

	ebitenutil.DebugPrintAt(screen, "Controles: [W/S/A/D] Conduzir | [ESPACO] Arremessar | [SHIFT] Bunny-Hop | [H] Buzina", 65, 300)
	ebitenutil.DebugPrintAt(screen, "Go 1.26+ | Ebitengine v2 | Antigravity Engine", 195, 330)
}

func (g *Game) drawCharacterPreview(screen *ebiten.Image) {
	// Posição no painel esquerdo da Central do Entregador
	cx := float64(130)
	cy := float64(160)

	// Desenha a sombra no chão
	vector.DrawFilledCircle(screen, float32(cx+16), float32(cy+40), 22, color.RGBA{0, 0, 0, 70}, false)

	// Desenha o veículo escolhido
	g.drawOpts.GeoM.Reset()
	g.drawOpts.GeoM.Scale(1.8, 1.8)
	g.drawOpts.GeoM.Translate(cx-18, cy+8)

	switch g.custom.VehicleType {
	case 1:
		screen.DrawImage(g.atlas.Scooter, &g.drawOpts)
	case 2:
		screen.DrawImage(g.atlas.DeliveryVan, &g.drawOpts)
	default:
		screen.DrawImage(g.atlas.BicycleFrame, &g.drawOpts)
	}

	// Desenha o Mascote no cesto dianteiro
	g.drawOpts.GeoM.Reset()
	g.drawOpts.GeoM.Scale(1.4, 1.4)
	g.drawOpts.GeoM.Translate(cx+36, cy+18)
	switch g.custom.Companion {
	case 0:
		screen.DrawImage(g.atlas.CarameloDog, &g.drawOpts)
	case 1:
		screen.DrawImage(g.atlas.Capybara, &g.drawOpts)
	case 2:
		screen.DrawImage(g.atlas.MiniDrone, &g.drawOpts)
	}

	// Desenha o personagem customizado montado
	charSprite := g.atlas.GenerateCustomCharacterSprite(g.custom, 0)
	g.drawOpts.GeoM.Reset()
	g.drawOpts.GeoM.Scale(1.8, 1.8)
	g.drawOpts.GeoM.Translate(cx-6, cy-22)
	screen.DrawImage(charSprite, &g.drawOpts)
}

func (g *Game) drawGameplay(screen *ebiten.Image) {
	p := g.player

	// Tremor de tela de Pânico
	shakeX := 0.0
	shakeY := 0.0
	if g.screenShake > 0 {
		shakeX = (rand.Float64()*2.0 - 1.0) * g.screenShake
		shakeY = (rand.Float64()*2.0 - 1.0) * g.screenShake
	}

	// 1. Fundo de Grama Verdejante
	vector.DrawFilledRect(screen, 0, 0, float32(config.VirtualWidth), float32(config.VirtualHeight), color.RGBA{48, 130, 52, 255}, false)

	// 2. Desenho das Faixas de Asfalto e Calçadas em perspectiva isométrica
	startTileY := int(p.Y-150.0) / 32 * 32
	endTileY := int(p.Y+450.0) / 32 * 32

	for ty := float64(startTileY); ty <= float64(endTileY); ty += 32.0 {
		// Pista Central de Asfalto
		sx, sy := render.IsoProject(0, ty, 0, p.X, p.Y)
		if render.IsOnScreen(sx, sy, 80) {
			g.drawOpts.GeoM.Reset()
			g.drawOpts.GeoM.Scale(2.8, 1.5)
			g.drawOpts.GeoM.Translate(sx-90+shakeX, sy-24+shakeY)
			screen.DrawImage(g.atlas.RoadTile, &g.drawOpts)
		}

		// Calçada Esquerda
		lsx, lsy := render.IsoProject(-config.StreetWidth/2.0-30.0, ty, 0, p.X, p.Y)
		if render.IsOnScreen(lsx, lsy, 80) {
			g.drawOpts.GeoM.Reset()
			g.drawOpts.GeoM.Scale(1.5, 1.2)
			g.drawOpts.GeoM.Translate(lsx-48+shakeX, lsy-18+shakeY)
			screen.DrawImage(g.atlas.SidewalkTile, &g.drawOpts)
		}

		// Calçada Direita
		rsx, rsy := render.IsoProject(config.StreetWidth/2.0+30.0, ty, 0, p.X, p.Y)
		if render.IsOnScreen(rsx, rsy, 80) {
			g.drawOpts.GeoM.Reset()
			g.drawOpts.GeoM.Scale(1.5, 1.2)
			g.drawOpts.GeoM.Translate(rsx-48+shakeX, rsy-18+shakeY)
			screen.DrawImage(g.atlas.SidewalkTile, &g.drawOpts)
		}
	}

	// 3. Coleta de Entidades para Ordenação por Profundidade (Y-Sorting)
	var items []RenderItem

	// Casas e Smart Lockers
	for _, h := range g.houses {
		hx, hy := render.IsoProject(h.WorldX, h.WorldY, 0, p.X, p.Y)
		if render.IsOnScreen(hx, hy, 120) {
			house := h
			items = append(items, RenderItem{
				Depth: hy + 20.0,
				Draw: func(s *ebiten.Image) {
					var opt ebiten.DrawImageOptions
					opt.GeoM.Translate(hx-48+shakeX, hy-90+shakeY)

					if house.Style == 3 {
						s.DrawImage(g.atlas.SmartLocker, &opt)
					} else {
						s.DrawImage(g.atlas.HouseStyles[house.Style], &opt)
					}

					// Retículo de mira no alvo de entrega
					if house.Status == entities.StatusPending {
						rx, ry := render.IsoProject(house.TargetX, house.TargetY, 0, p.X, p.Y)
						var retOpt ebiten.DrawImageOptions
						retOpt.GeoM.Scale(0.7, 0.7)
						retOpt.GeoM.Translate(rx-11+shakeX, ry-11+shakeY)
						s.DrawImage(g.atlas.TargetReticle, &retOpt)
					} else if house.CustomerHappy {
						// Cliente feliz comemorando a entrega
						ebitenutil.DebugPrintAt(s, "Obrigado! 5★", int(hx)-15, int(hy)-105)
					}
				},
			})
		}
	}

	// Obstáculos na Pista
	for _, obs := range g.obstacles {
		ox, oy := render.IsoProject(obs.WorldX, obs.WorldY, 0, p.X, p.Y)
		if render.IsOnScreen(ox, oy, 60) {
			o := obs
			items = append(items, RenderItem{
				Depth: oy,
				Draw: func(s *ebiten.Image) {
					var opt ebiten.DrawImageOptions
					opt.GeoM.Translate(ox-16+shakeX, oy-12+shakeY)
					switch o.Type {
					case entities.ObsPothole:
						s.DrawImage(g.atlas.Pothole, &opt)
					case entities.ObsPuddle:
						s.DrawImage(g.atlas.Puddle, &opt)
					case entities.ObsTrafficCone:
						s.DrawImage(g.atlas.TrafficCone, &opt)
					case entities.ObsRoadBarrier:
						s.DrawImage(g.atlas.RoadBarrier, &opt)
					case entities.ObsBarkingDog:
						if o.IsScared {
							// Cão assustado fugindo
							opt.GeoM.Translate(15, -10)
						}
						s.DrawImage(g.atlas.BarkingDog, &opt)
					case entities.ObsSprinkler:
						s.DrawImage(g.atlas.Sprinkler, &opt)
					}
				},
			})
		}
	}

	// Pacotes em voo
	for _, pkg := range g.packages {
		px, py := render.IsoProject(pkg.CurrentX, pkg.CurrentY, pkg.CurrentZ, p.X, p.Y)
		pkgCopy := pkg
		items = append(items, RenderItem{
			Depth: py + pkg.CurrentZ,
			Draw: func(s *ebiten.Image) {
				var opt ebiten.DrawImageOptions
				// Sombra do pacote no chão
				shadowX, shadowY := render.IsoProject(pkgCopy.CurrentX, pkgCopy.CurrentY, 0, p.X, p.Y)
				vector.DrawFilledCircle(s, float32(shadowX), float32(shadowY), 4, color.RGBA{0, 0, 0, 80}, false)

				opt.GeoM.Translate(px-9+shakeX, py-8+shakeY)
				if pkgCopy.PackageType == 1 {
					s.DrawImage(g.atlas.PackageFragile, &opt)
				} else if pkgCopy.PackageType == 2 {
					s.DrawImage(g.atlas.PackageLarge, &opt)
				} else {
					s.DrawImage(g.atlas.PackageYellow, &opt)
				}
			},
		})
	}

	// Chefes Ativos na Tela
	for _, b := range g.bosses {
		if b.State == entities.BossPanic || b.State == entities.BossRecover {
			bx, by := render.IsoProject(0, b.TriggerY+40.0, 0, p.X, p.Y)
			bossCopy := b
			items = append(items, RenderItem{
				Depth: by + 10.0,
				Draw: func(s *ebiten.Image) {
					var opt ebiten.DrawImageOptions
					opt.GeoM.Scale(1.6, 1.6)
					opt.GeoM.Translate(bx-40+shakeX, by-60+shakeY)
					switch bossCopy.Type {
					case entities.BossTornado:
						s.DrawImage(g.atlas.TornadoFunnel, &opt)
					case entities.BossCrater:
						s.DrawImage(g.atlas.CraterMonster, &opt)
					case entities.BossProtest:
						s.DrawImage(g.atlas.ProtestSmoke, &opt)
					case entities.BossBlackFriday:
						s.DrawImage(g.atlas.ColossusBoss, &opt)
					default:
						s.DrawImage(g.atlas.BarkingDog, &opt)
					}
				},
			})
		}
	}

	// Jogador + Bicicleta + Entregador Customizado
	px, py := render.IsoProject(p.X, p.Y, p.Z, p.X, p.Y)
	items = append(items, RenderItem{
		Depth: py + p.Z,
		Draw: func(s *ebiten.Image) {
			// Sombra no chão
			shadowX, shadowY := render.IsoProject(p.X, p.Y, 0, p.X, p.Y)
			vector.DrawFilledCircle(s, float32(shadowX), float32(shadowY), 16, color.RGBA{0, 0, 0, 90}, false)

			// Efeito de piscar quando invulnerável
			if p.InvulnTimer > 0 && int(p.InvulnTimer*12)%2 == 0 {
				return
			}

			// Veículo
			var vehOpt ebiten.DrawImageOptions
			vehOpt.GeoM.Translate(px-21+shakeX, py-16+shakeY)
			switch p.Custom.VehicleType {
			case 1:
				s.DrawImage(g.atlas.Scooter, &vehOpt)
			case 2:
				s.DrawImage(g.atlas.DeliveryVan, &vehOpt)
			default:
				s.DrawImage(g.atlas.BicycleFrame, &vehOpt)
			}

			// Mascote no cesto
			var compOpt ebiten.DrawImageOptions
			compOpt.GeoM.Translate(px+10+shakeX, py-10+shakeY)
			switch p.Custom.Companion {
			case 0:
				s.DrawImage(g.atlas.CarameloDog, &compOpt)
			case 1:
				s.DrawImage(g.atlas.Capybara, &compOpt)
			case 2:
				s.DrawImage(g.atlas.MiniDrone, &compOpt)
			}

			// Entregador Customizado montado
			charSprite := g.atlas.GenerateCustomCharacterSprite(p.Custom, p.AnimFrame)
			var charOpt ebiten.DrawImageOptions
			charOpt.GeoM.Translate(px-14+shakeX, py-34+shakeY)
			s.DrawImage(charSprite, &charOpt)
		},
	})

	// Partículas
	for _, pt := range g.particles.GetActiveParticles() {
		if pt.Active {
			part := pt
			ptx, pty := render.IsoProject(part.X, part.Y, part.Z, p.X, p.Y)
			if render.IsOnScreen(ptx, pty, 10) {
				items = append(items, RenderItem{
					Depth: pty + part.Z,
					Draw: func(s *ebiten.Image) {
						vector.DrawFilledCircle(s, float32(ptx), float32(pty), float32(part.Size), part.Color, false)
					},
				})
			}
		}
	}

	// 4. Executa o Y-Sorting e Desenha Tudo na Ordem Correta
	sort.Slice(items, func(i, j int) bool {
		return items[i].Depth < items[j].Depth
	})

	for _, it := range items {
		it.Draw(screen)
	}

	// 5. Vinheta de Alarme Vermelho de Pânico (se ativo)
	if g.redFlash > 0.05 {
		alpha := uint8(g.redFlash * 90.0)
		vector.DrawFilledRect(screen, 0, 0, float32(config.VirtualWidth), float32(config.VirtualHeight), color.RGBA{255, 0, 0, alpha}, false)
	}

	// 6. HUD e Minimapa
	var activeBoss *entities.BossEvent
	for _, b := range g.bosses {
		if b.State == entities.BossPanic || b.State == entities.BossRecover {
			activeBoss = b
			break
		}
	}
	ui.DrawHUD(screen, p, activeBoss, config.RouteLength)
}
