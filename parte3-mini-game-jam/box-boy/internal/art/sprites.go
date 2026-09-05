package art

import (
	"image"
	"image/color"
	"math"

	"box-boy/internal/customizer"
	"github.com/hajimehoshi/ebiten/v2"
)

// TextureAtlas armazena todas as texturas procedurais pré-geradas em memória.
type TextureAtlas struct {
	// Terreno e Cenário
	RoadTile      *ebiten.Image
	SidewalkTile  *ebiten.Image
	GrassTile     *ebiten.Image
	CurbTile      *ebiten.Image
	HouseStyles   []*ebiten.Image
	SmartLocker   *ebiten.Image
	Mailbox       *ebiten.Image
	TargetReticle *ebiten.Image

	// Veículos
	BicycleFrame *ebiten.Image
	BicycleWheel *ebiten.Image
	Scooter      *ebiten.Image
	DeliveryVan  *ebiten.Image

	// Mascotes do Cesto
	CarameloDog *ebiten.Image
	Capybara    *ebiten.Image
	MiniDrone   *ebiten.Image

	// Encomendas
	PackageYellow  *ebiten.Image
	PackageFragile *ebiten.Image
	PackageLarge   *ebiten.Image

	// Obstáculos
	Pothole      *ebiten.Image
	Puddle       *ebiten.Image
	TrafficCone  *ebiten.Image
	RoadBarrier  *ebiten.Image
	BarkingDog   *ebiten.Image
	Sprinkler    *ebiten.Image

	// Bosses & Desastres (SimCity Style)
	TornadoFunnel *ebiten.Image
	CraterMonster *ebiten.Image
	ProtestSmoke  *ebiten.Image
	ProtestBanner *ebiten.Image
	ColossusBoss  *ebiten.Image

	// Partículas & UI
	StarIcon    *ebiten.Image
	ParticleDot *ebiten.Image
	SirenIcon   *ebiten.Image
}

// NewTextureAtlas constrói todo o atlas de sprites em memória (zero disco).
func NewTextureAtlas() *TextureAtlas {
	atlas := &TextureAtlas{}
	atlas.buildEnvironment()
	atlas.buildVehicles()
	atlas.buildCompanions()
	atlas.buildPackages()
	atlas.buildObstacles()
	atlas.buildBosses()
	atlas.buildUI()
	return atlas
}

// setPixel define um pixel com segurança nos limites da imagem
func setPixel(img *image.RGBA, x, y int, c color.RGBA) {
	if x >= 0 && x < img.Rect.Dx() && y >= 0 && y < img.Rect.Dy() {
		img.SetRGBA(x, y, c)
	}
}

// drawRect preenche um retângulo com uma cor
func drawRect(img *image.RGBA, x0, y0, w, h int, c color.RGBA) {
	for y := y0; y < y0+h; y++ {
		for x := x0; x < x0+w; x++ {
			setPixel(img, x, y, c)
		}
	}
}

// drawCircle preenche um círculo com uma cor
func drawCircle(img *image.RGBA, cx, cy, r int, c color.RGBA) {
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			dx := x - cx
			dy := y - cy
			if dx*dx+dy*dy <= r*r {
				setPixel(img, x, y, c)
			}
		}
	}
}

func (a *TextureAtlas) buildEnvironment() {
	// 1. Asfalto da Rua Isométrica (64x32 losango)
	roadRaw := image.NewRGBA(image.Rect(0, 0, 64, 32))
	asphalt := color.RGBA{45, 48, 55, 255}
	asphaltLine := color.RGBA{240, 220, 60, 255} // Faixa amarela central
	curbDark := color.RGBA{30, 32, 38, 255}

	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			// Teste de ponto em losango isométrico: |dx/32| + |dy/16| <= 1
			dx := float64(x - 32) / 32.0
			dy := float64(y - 16) / 16.0
			if math.Abs(dx)+math.Abs(dy) <= 1.0 {
				setPixel(roadRaw, x, y, asphalt)
				// Faixa amarela tracejada no meio
				if math.Abs(dx) < 0.08 && (y/4)%2 == 0 {
					setPixel(roadRaw, x, y, asphaltLine)
				}
			}
		}
	}
	a.RoadTile = ebiten.NewImageFromImage(roadRaw)

	// 2. Calçada Portuguesa / Piso Urbano (64x32)
	sideRaw := image.NewRGBA(image.Rect(0, 0, 64, 32))
	sideColor1 := color.RGBA{195, 190, 185, 255}
	sideColor2 := color.RGBA{170, 165, 160, 255}
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			dx := float64(x - 32) / 32.0
			dy := float64(y - 16) / 16.0
			if math.Abs(dx)+math.Abs(dy) <= 1.0 {
				if (x/4+y/2)%2 == 0 {
					setPixel(sideRaw, x, y, sideColor1)
				} else {
					setPixel(sideRaw, x, y, sideColor2)
				}
			}
		}
	}
	a.SidewalkTile = ebiten.NewImageFromImage(sideRaw)

	// 3. Gramado com Flores (64x32)
	grassRaw := image.NewRGBA(image.Rect(0, 0, 64, 32))
	grassBase := color.RGBA{58, 145, 62, 255}
	grassHigh := color.RGBA{72, 175, 78, 255}
	flowerCol := color.RGBA{255, 220, 80, 255}
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			dx := float64(x - 32) / 32.0
			dy := float64(y - 16) / 16.0
			if math.Abs(dx)+math.Abs(dy) <= 1.0 {
				if (x*7+y*13)%5 == 0 {
					setPixel(grassRaw, x, y, grassHigh)
				} else if (x*11+y*3)%17 == 0 {
					setPixel(grassRaw, x, y, flowerCol)
				} else {
					setPixel(grassRaw, x, y, grassBase)
				}
			}
		}
	}
	a.GrassTile = ebiten.NewImageFromImage(grassRaw)

	// 4. Meio-fio / Guia de Calçada
	curbRaw := image.NewRGBA(image.Rect(0, 0, 64, 16))
	drawRect(curbRaw, 0, 0, 64, 16, curbDark)
	drawRect(curbRaw, 0, 0, 64, 4, color.RGBA{220, 220, 220, 255})
	a.CurbTile = ebiten.NewImageFromImage(curbRaw)

	// 5. Estilos de Casas e Fachadas de Bairro (96x110)
	houseColors := []struct {
		wall color.RGBA
		roof color.RGBA
		door color.RGBA
	}{
		{color.RGBA{245, 235, 210, 255}, color.RGBA{180, 60, 45, 255}, color.RGBA{70, 40, 25, 255}},  // Casa Aconchegante
		{color.RGBA{210, 235, 250, 255}, color.RGBA{40, 70, 140, 255}, color.RGBA{50, 60, 80, 255}},   // Casa Moderna
		{color.RGBA{255, 240, 175, 255}, color.RGBA{200, 100, 30, 255}, color.RGBA{85, 45, 20, 255}}, // Casa Tropical
	}

	for _, hc := range houseColors {
		hRaw := image.NewRGBA(image.Rect(0, 0, 96, 110))
		// Paredes principais
		drawRect(hRaw, 8, 35, 80, 70, hc.wall)
		// Contorno e sombras
		drawRect(hRaw, 8, 35, 3, 70, color.RGBA{0, 0, 0, 45})
		drawRect(hRaw, 85, 35, 3, 70, color.RGBA{0, 0, 0, 65})
		// Telhado triangular / empena
		for py := 0; py < 35; py++ {
			span := py * 48 / 35
			drawRect(hRaw, 48-span, py, span*2, 1, hc.roof)
		}
		// Porta principal (Alvo de entrega!)
		drawRect(hRaw, 38, 72, 20, 33, hc.door)
		drawCircle(hRaw, 54, 88, 2, color.RGBA{255, 215, 0, 255}) // Maçaneta dourada
		// Janelas com vidros azuis reluzentes
		drawRect(hRaw, 16, 48, 18, 18, color.RGBA{135, 206, 250, 255})
		drawRect(hRaw, 62, 48, 18, 18, color.RGBA{135, 206, 250, 255})
		// Cruzeta branca nas janelas
		drawRect(hRaw, 24, 48, 2, 18, color.RGBA{255, 255, 255, 220})
		drawRect(hRaw, 16, 56, 18, 2, color.RGBA{255, 255, 255, 220})
		drawRect(hRaw, 70, 48, 2, 18, color.RGBA{255, 255, 255, 220})
		drawRect(hRaw, 62, 56, 18, 2, color.RGBA{255, 255, 255, 220})
		// Varandinha / Tapete de entrada amarelo
		drawRect(hRaw, 34, 102, 28, 6, color.RGBA{255, 230, 0, 255})

		a.HouseStyles = append(a.HouseStyles, ebiten.NewImageFromImage(hRaw))
	}

	// 6. Smart Locker (Amarelo Canário Icônico com tela azul)
	lRaw := image.NewRGBA(image.Rect(0, 0, 44, 60))
	drawRect(lRaw, 2, 4, 40, 54, color.RGBA{255, 230, 0, 255}) // Amarelo Express
	drawRect(lRaw, 2, 4, 40, 2, color.RGBA{30, 80, 160, 255})  // Faixa Azul Express
	// Portas dos armários
	for gy := 10; gy < 52; gy += 10 {
		for gx := 6; gx < 38; gx += 16 {
			drawRect(lRaw, gx, gy, 14, 8, color.RGBA{240, 210, 0, 255})
			drawRect(lRaw, gx+10, gy+3, 2, 2, color.RGBA{50, 50, 50, 255}) // Trincos
		}
	}
	// Tela digital azul central
	drawRect(lRaw, 15, 22, 14, 10, color.RGBA{30, 80, 200, 255})
	drawRect(lRaw, 17, 24, 10, 2, color.RGBA{180, 235, 255, 255})
	a.SmartLocker = ebiten.NewImageFromImage(lRaw)

	// 7. Caixa de Correio Residencial
	mRaw := image.NewRGBA(image.Rect(0, 0, 20, 36))
	drawRect(mRaw, 8, 16, 4, 20, color.RGBA{90, 85, 95, 255})   // Haste
	drawRect(mRaw, 2, 4, 16, 14, color.RGBA{220, 50, 45, 255})  // Caixa vermelha
	drawRect(mRaw, 15, 8, 3, 6, color.RGBA{255, 215, 0, 255})   // Bandeirinha
	a.Mailbox = ebiten.NewImageFromImage(mRaw)

	// 8. Mira / Target Reticle de Arremesso (Anel pulsante)
	tRaw := image.NewRGBA(image.Rect(0, 0, 32, 32))
	drawCircle(tRaw, 16, 16, 12, color.RGBA{255, 230, 0, 180})
	drawCircle(tRaw, 16, 16, 8, color.RGBA{0, 0, 0, 0})
	drawRect(tRaw, 14, 2, 4, 28, color.RGBA{30, 110, 240, 220})
	drawRect(tRaw, 2, 14, 28, 4, color.RGBA{30, 110, 240, 220})
	a.TargetReticle = ebiten.NewImageFromImage(tRaw)
}

func (a *TextureAtlas) buildVehicles() {
	// 1. Bicicleta Amarela de Entrega (42x32)
	bikeRaw := image.NewRGBA(image.Rect(0, 0, 42, 32))
	yellowFrame := color.RGBA{255, 225, 0, 255}
	chrome := color.RGBA{190, 195, 205, 255}
	darkRubber := color.RGBA{35, 35, 40, 255}

	// Rodas com pneus pretos e raios cromados
	drawCircle(bikeRaw, 9, 23, 7, darkRubber)
	drawCircle(bikeRaw, 9, 23, 5, chrome)
	drawCircle(bikeRaw, 9, 23, 2, darkRubber)

	drawCircle(bikeRaw, 33, 23, 7, darkRubber)
	drawCircle(bikeRaw, 33, 23, 5, chrome)
	drawCircle(bikeRaw, 33, 23, 2, darkRubber)

	// Quadro amarelo triangular
	drawRect(bikeRaw, 9, 21, 14, 3, yellowFrame)
	drawRect(bikeRaw, 21, 13, 3, 11, yellowFrame)
	drawRect(bikeRaw, 21, 13, 12, 3, yellowFrame)
	drawRect(bikeRaw, 21, 13, 14, 12, yellowFrame)

	// Guidão & Selim
	drawRect(bikeRaw, 31, 8, 3, 10, chrome)
	drawRect(bikeRaw, 29, 6, 7, 3, darkRubber)
	drawRect(bikeRaw, 18, 10, 7, 3, darkRubber) // Selim

	// Cesto dianteiro amarelo para mascote / pacotes
	drawRect(bikeRaw, 33, 10, 8, 7, color.RGBA{240, 200, 0, 255})
	a.BicycleFrame = ebiten.NewImageFromImage(bikeRaw)

	// 2. Scooter Elétrica (46x34)
	scooterRaw := image.NewRGBA(image.Rect(0, 0, 46, 34))
	drawCircle(scooterRaw, 10, 26, 6, darkRubber)
	drawCircle(scooterRaw, 36, 26, 6, darkRubber)
	drawRect(scooterRaw, 8, 14, 30, 14, yellowFrame)
	drawRect(scooterRaw, 34, 6, 4, 16, yellowFrame)
	drawRect(scooterRaw, 32, 4, 8, 3, darkRubber)
	// Farol LED frontal azul
	drawRect(scooterRaw, 39, 10, 4, 4, color.RGBA{100, 200, 255, 255})
	// Baú traseiro grande
	drawRect(scooterRaw, 4, 8, 16, 14, color.RGBA{255, 230, 0, 255})
	drawRect(scooterRaw, 4, 13, 16, 3, color.RGBA{30, 90, 200, 255}) // Faixa azul
	a.Scooter = ebiten.NewImageFromImage(scooterRaw)

	// 3. Furgão Compacto Elétrico (64x38)
	vanRaw := image.NewRGBA(image.Rect(0, 0, 64, 38))
	drawCircle(vanRaw, 14, 31, 6, darkRubber)
	drawCircle(vanRaw, 50, 31, 6, darkRubber)
	drawRect(vanRaw, 4, 8, 56, 23, yellowFrame)
	drawRect(vanRaw, 42, 11, 14, 9, color.RGBA{140, 215, 255, 255}) // Pára-brisa
	drawRect(vanRaw, 4, 20, 56, 4, color.RGBA{30, 90, 200, 255})     // Faixa azul lateral
	drawRect(vanRaw, 58, 23, 4, 4, color.RGBA{255, 255, 200, 255})   // Farol
	a.DeliveryVan = ebiten.NewImageFromImage(vanRaw)
}

func (a *TextureAtlas) buildCompanions() {
	// 1. Cão Caramelo (16x14) - Orelhas caídas e focinho preto
	dogRaw := image.NewRGBA(image.Rect(0, 0, 16, 14))
	caramelo := color.RGBA{218, 145, 60, 255}
	carameloDark := color.RGBA{170, 105, 35, 255}
	black := color.RGBA{20, 20, 20, 255}

	drawCircle(dogRaw, 8, 7, 5, caramelo)
	drawCircle(dogRaw, 12, 8, 3, caramelo)     // Focinho
	setPixel(dogRaw, 14, 7, black)             // Nariz
	setPixel(dogRaw, 10, 5, black)             // Olho
	drawRect(dogRaw, 4, 3, 4, 6, carameloDark) // Orelha caída
	// Coleira azul
	drawRect(dogRaw, 7, 10, 3, 2, color.RGBA{30, 100, 230, 255})
	a.CarameloDog = ebiten.NewImageFromImage(dogRaw)

	// 2. Capivara Zen (18x14)
	capyRaw := image.NewRGBA(image.Rect(0, 0, 18, 14))
	capyBrown := color.RGBA{160, 110, 70, 255}
	drawRect(capyRaw, 3, 4, 13, 9, capyBrown)
	drawCircle(capyRaw, 13, 6, 4, capyBrown)
	// Óculos escuros na capivara
	drawRect(capyRaw, 10, 4, 6, 3, color.RGBA{20, 20, 20, 255})
	a.Capybara = ebiten.NewImageFromImage(capyRaw)

	// 3. Mini Drone (16x12)
	droneRaw := image.NewRGBA(image.Rect(0, 0, 16, 12))
	drawRect(droneRaw, 4, 5, 8, 5, color.RGBA{255, 230, 0, 255})
	// Hélices
	drawRect(droneRaw, 1, 3, 5, 2, color.RGBA{200, 210, 225, 255})
	drawRect(droneRaw, 10, 3, 5, 2, color.RGBA{200, 210, 225, 255})
	// Luz de navegação piscante
	drawCircle(droneRaw, 8, 7, 2, color.RGBA{50, 170, 255, 255})
	a.MiniDrone = ebiten.NewImageFromImage(droneRaw)
}

func (a *TextureAtlas) buildPackages() {
	// 1. Pacote Amarelo Pequeno com fita azul
	p1Raw := image.NewRGBA(image.Rect(0, 0, 18, 16))
	boxYellow := color.RGBA{255, 220, 40, 255}
	tapeBlue := color.RGBA{25, 95, 215, 255}
	drawRect(p1Raw, 1, 1, 16, 14, boxYellow)
	drawRect(p1Raw, 0, 0, 18, 1, color.RGBA{210, 175, 20, 255})
	// Fita adesiva azul cruzando a caixa
	drawRect(p1Raw, 7, 1, 4, 14, tapeBlue)
	drawRect(p1Raw, 1, 6, 16, 3, tapeBlue)
	a.PackageYellow = ebiten.NewImageFromImage(p1Raw)

	// 2. Pacote Frágil (Com ícone de taça vermelha)
	p2Raw := image.NewRGBA(image.Rect(0, 0, 20, 18))
	drawRect(p2Raw, 1, 1, 18, 16, color.RGBA{230, 200, 140, 255}) // Papelão Kraft
	drawRect(p2Raw, 8, 1, 4, 16, tapeBlue)
	// Selo de Taça Frágil
	drawRect(p2Raw, 3, 4, 5, 4, color.RGBA{220, 40, 30, 255})
	drawRect(p2Raw, 5, 8, 1, 3, color.RGBA{220, 40, 30, 255})
	a.PackageFragile = ebiten.NewImageFromImage(p2Raw)

	// 3. Pacote Grande Especial
	p3Raw := image.NewRGBA(image.Rect(0, 0, 26, 22))
	drawRect(p3Raw, 1, 1, 24, 20, boxYellow)
	drawRect(p3Raw, 11, 1, 5, 20, tapeBlue)
	drawRect(p3Raw, 1, 8, 24, 4, tapeBlue)
	a.PackageLarge = ebiten.NewImageFromImage(p3Raw)
}

func (a *TextureAtlas) buildObstacles() {
	// 1. Buraco no Asfalto (Cratera menor)
	potRaw := image.NewRGBA(image.Rect(0, 0, 30, 18))
	drawCircle(potRaw, 15, 9, 8, color.RGBA{20, 20, 25, 255})
	drawCircle(potRaw, 15, 9, 5, color.RGBA{10, 10, 12, 255})
	a.Pothole = ebiten.NewImageFromImage(potRaw)

	// 2. Poça de Água com Reflexo
	pudRaw := image.NewRGBA(image.Rect(0, 0, 32, 16))
	drawCircle(pudRaw, 16, 8, 7, color.RGBA{80, 160, 240, 190})
	drawCircle(pudRaw, 15, 7, 4, color.RGBA{180, 220, 255, 220})
	a.Puddle = ebiten.NewImageFromImage(pudRaw)

	// 3. Cone de Trânsito Laranja e Branco
	coneRaw := image.NewRGBA(image.Rect(0, 0, 18, 24))
	for y := 0; y < 20; y++ {
		w := 3 + y*11/20
		c := color.RGBA{255, 90, 20, 255}
		if y >= 6 && y <= 9 || y >= 13 && y <= 15 {
			c = color.RGBA{255, 255, 255, 255}
		}
		drawRect(coneRaw, 9-w/2, y, w, 1, c)
	}
	drawRect(coneRaw, 2, 20, 14, 4, color.RGBA{40, 40, 40, 255}) // Base preta
	a.TrafficCone = ebiten.NewImageFromImage(coneRaw)

	// 4. Barricada de Obras Listrada
	barRaw := image.NewRGBA(image.Rect(0, 0, 36, 26))
	drawRect(barRaw, 4, 14, 4, 12, color.RGBA{80, 80, 80, 255})
	drawRect(barRaw, 28, 14, 4, 12, color.RGBA{80, 80, 80, 255})
	for x := 0; x < 36; x++ {
		c := color.RGBA{240, 80, 30, 255}
		if (x/6)%2 == 0 {
			c = color.RGBA{255, 255, 255, 255}
		}
		drawRect(barRaw, x, 4, 1, 10, c)
	}
	a.RoadBarrier = ebiten.NewImageFromImage(barRaw)

	// 5. Cão Latindo no Portão (Feroz com olhos vermelhos)
	barkRaw := image.NewRGBA(image.Rect(0, 0, 24, 20))
	drawRect(barkRaw, 4, 6, 14, 10, color.RGBA{70, 50, 40, 255})
	drawCircle(barkRaw, 16, 7, 5, color.RGBA{70, 50, 40, 255})
	setPixel(barkRaw, 17, 6, color.RGBA{255, 30, 30, 255}) // Olho vermelho
	drawRect(barkRaw, 18, 9, 5, 3, color.RGBA{255, 255, 255, 255}) // Dentes
	a.BarkingDog = ebiten.NewImageFromImage(barkRaw)

	// 6. Aspersor de Jardim Giratório
	spRaw := image.NewRGBA(image.Rect(0, 0, 20, 20))
	drawCircle(spRaw, 10, 10, 4, color.RGBA{70, 160, 240, 255})
	drawRect(spRaw, 8, 2, 4, 16, color.RGBA{40, 110, 190, 255})
	a.Sprinkler = ebiten.NewImageFromImage(spRaw)
}

func (a *TextureAtlas) buildBosses() {
	// 1. O Tornado Metropolitano (SimCity Style) (56x72)
	torRaw := image.NewRGBA(image.Rect(0, 0, 56, 72))
	for y := 0; y < 72; y++ {
		// Funil se alarga no topo
		width := 8 + y*44/72
		cx := 28 + int(math.Sin(float64(y)*0.2)*4.0)
		for x := cx - width/2; x < cx+width/2; x++ {
			noiseAlpha := uint8(100 + (x*17+y*23)%120)
			setPixel(torRaw, x, 71-y, color.RGBA{200, 215, 235, noiseAlpha})
		}
	}
	// Detritos girando no tornado (pedaços de papelão e placas)
	setPixel(torRaw, 18, 30, color.RGBA{255, 220, 0, 255})
	setPixel(torRaw, 38, 20, color.RGBA{255, 80, 40, 255})
	setPixel(torRaw, 25, 12, color.RGBA{255, 220, 0, 255})
	a.TornadoFunnel = ebiten.NewImageFromImage(torRaw)

	// 2. A Cratera Voraz / Monstro do Asfalto (64x36)
	cratRaw := image.NewRGBA(image.Rect(0, 0, 64, 36))
	drawCircle(cratRaw, 32, 18, 16, color.RGBA{15, 15, 20, 255})
	drawCircle(cratRaw, 32, 18, 11, color.RGBA{120, 60, 20, 255}) // Lama vulcânica
	// Tubulação estourada jorrando água
	drawRect(cratRaw, 22, 14, 20, 4, color.RGBA{70, 170, 250, 220})
	a.CraterMonster = ebiten.NewImageFromImage(cratRaw)

	// 3. Fumaça e Barricada da Manifestação
	smkRaw := image.NewRGBA(image.Rect(0, 0, 48, 36))
	drawCircle(smkRaw, 24, 18, 14, color.RGBA{255, 120, 40, 160})
	drawCircle(smkRaw, 20, 14, 10, color.RGBA{255, 210, 40, 190})
	a.ProtestSmoke = ebiten.NewImageFromImage(smkRaw)

	banRaw := image.NewRGBA(image.Rect(0, 0, 60, 24))
	drawRect(banRaw, 2, 2, 56, 20, color.RGBA{220, 30, 40, 255})
	drawRect(banRaw, 4, 8, 52, 6, color.RGBA{255, 255, 255, 255}) // Texto estilizado da faixa
	a.ProtestBanner = ebiten.NewImageFromImage(banRaw)

	// 4. Mega Boss: O Monstro do Atraso da Black Friday (72x80)
	colRaw := image.NewRGBA(image.Rect(0, 0, 72, 80))
	drawRect(colRaw, 14, 22, 44, 48, color.RGBA{50, 55, 65, 255}) // Corpo mecânico
	// Relógio de contagem regressiva gigante no peito
	drawCircle(colRaw, 36, 46, 16, color.RGBA{240, 240, 245, 255})
	drawCircle(colRaw, 36, 46, 13, color.RGBA{255, 40, 40, 255}) // Mostrador vermelho
	drawRect(colRaw, 35, 36, 2, 10, color.RGBA{20, 20, 20, 255})   // Ponteiro
	// Olhos de LED escaneadores
	drawRect(colRaw, 24, 12, 8, 4, color.RGBA{255, 20, 20, 255})
	drawRect(colRaw, 40, 12, 8, 4, color.RGBA{255, 20, 20, 255})
	// Braços de esteira rolante
	drawRect(colRaw, 2, 34, 12, 24, color.RGBA{85, 90, 100, 255})
	drawRect(colRaw, 58, 34, 12, 24, color.RGBA{85, 90, 100, 255})
	a.ColossusBoss = ebiten.NewImageFromImage(colRaw)
}

func (a *TextureAtlas) buildUI() {
	// 1. Estrela Dourada (16x16)
	starRaw := image.NewRGBA(image.Rect(0, 0, 16, 16))
	gold := color.RGBA{255, 215, 0, 255}
	drawCircle(starRaw, 8, 8, 6, gold)
	drawRect(starRaw, 7, 1, 2, 14, gold)
	drawRect(starRaw, 1, 7, 14, 2, gold)
	a.StarIcon = ebiten.NewImageFromImage(starRaw)

	// 2. Partícula Ponto Brilhante (4x4)
	partRaw := image.NewRGBA(image.Rect(0, 0, 4, 4))
	drawRect(partRaw, 0, 0, 4, 4, color.RGBA{255, 255, 255, 255})
	a.ParticleDot = ebiten.NewImageFromImage(partRaw)

	// 3. Ícone de Sirene de Pânico (18x18)
	sirRaw := image.NewRGBA(image.Rect(0, 0, 18, 18))
	drawCircle(sirRaw, 9, 9, 7, color.RGBA{255, 30, 30, 255})
	drawRect(sirRaw, 7, 3, 4, 12, color.RGBA{255, 255, 255, 255})
	a.SirenIcon = ebiten.NewImageFromImage(sirRaw)
}

// GenerateCustomCharacterSprite monta o sprite do entregador combinando as camadas customizadas.
func (a *TextureAtlas) GenerateCustomCharacterSprite(c customizer.Customization, frame int) *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, 32, 42))

	skinCol := customizer.SkinColors[c.SkinTone]
	hairCol := customizer.HairColors[c.HairColor]

	// Balanço sutil de pedalada / corrida baseado no frame
	bob := 0
	if frame%2 == 1 {
		bob = 1
	}

	// 1. Tronco / Roupa do Entregador
	jacketYellow := color.RGBA{255, 225, 0, 255}
	blueStripe := color.RGBA{25, 95, 215, 255}
	vestNeon := color.RGBA{180, 255, 20, 255}
	hoodieGray := color.RGBA{65, 70, 80, 255}

	var outfitColor color.RGBA
	switch c.Outfit {
	case 1:
		outfitColor = vestNeon
	case 2:
		outfitColor = hoodieGray
	case 3:
		outfitColor = blueStripe
	default:
		outfitColor = jacketYellow
	}

	// Cabeça & Rosto
	headY := 7 + bob
	drawCircle(img, 16, headY, 6, skinCol)
	// Olhos
	setPixel(img, 18, headY, color.RGBA{20, 20, 20, 255})
	setPixel(img, 20, headY, color.RGBA{20, 20, 20, 255})

	// Óculos
	switch c.Glasses {
	case 1: // Óculos Ciclista Espelhado Neon
	drawRect(img, 17, headY-1, 6, 3, color.RGBA{50, 230, 255, 255})
	case 2: // Óculos Nerd
		drawRect(img, 17, headY-1, 6, 3, color.RGBA{30, 30, 30, 255})
	case 3: // Escuros
		drawRect(img, 17, headY-1, 6, 3, color.RGBA{10, 10, 15, 255})
	}

	// Cabelo
	switch c.HairStyle {
	case 0: // Degradê
		drawRect(img, 12, headY-6, 8, 3, hairCol)
	case 1: // Coque Samurai
		drawRect(img, 12, headY-6, 8, 3, hairCol)
		drawCircle(img, 12, headY-7, 3, hairCol)
	case 2: // Black Power Afro
		drawCircle(img, 15, headY-4, 7, hairCol)
	case 3: // Moicano
		drawRect(img, 15, headY-8, 3, 6, hairCol)
	case 4: // Dreadlocks
		drawRect(img, 11, headY-6, 9, 7, hairCol)
	case 5: // Longo
		drawRect(img, 11, headY-5, 8, 9, hairCol)
	}

	// Acessório de Cabeça (Headgear)
	switch c.Headgear {
	case 0: // Boné Aba Reta (virado pra trás)
		drawRect(img, 12, headY-6, 9, 3, jacketYellow)
		drawRect(img, 9, headY-4, 4, 2, blueStripe) // Aba para trás
	case 1: // Viseira Retrô
		drawRect(img, 13, headY-4, 9, 2, jacketYellow)
	case 2: // Capacete
		drawCircle(img, 16, headY-3, 7, color.RGBA{240, 50, 40, 255})
	case 3: // Bandana
		drawRect(img, 12, headY-4, 9, 2, blueStripe)
	}

	// Corpo / Jaqueta
	bodyY := 14 + bob
	drawRect(img, 12, bodyY, 9, 13, outfitColor)
	drawRect(img, 12, bodyY+5, 9, 3, blueStripe) // Faixa azul

	// Mochila / Baú Térmico Amarelo nas costas
	drawRect(img, 6, bodyY+1, 6, 12, jacketYellow)
	drawRect(img, 6, bodyY+6, 6, 2, blueStripe)

	// Braços segurando o guidão
	drawRect(img, 18, bodyY+3, 6, 3, skinCol)

	// Pernas / Pedalando
	drawRect(img, 13, bodyY+13, 3, 8-bob*2, color.RGBA{40, 45, 55, 255})
	drawRect(img, 17, bodyY+13, 3, 6+bob*2, color.RGBA{40, 45, 55, 255})

	// Tênis amarelo
	drawRect(img, 13, bodyY+20-bob*2, 4, 3, jacketYellow)
	drawRect(img, 17, bodyY+18+bob*2, 4, 3, jacketYellow)

	return ebiten.NewImageFromImage(img)
}
