package ui

import (
	"fmt"
	"image/color"

	"box-boy/internal/config"
	"box-boy/internal/customizer"
	"box-boy/internal/entities"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// DrawHUD renderiza a interface superior e o minimapa inferior durante o gameplay.
func DrawHUD(screen *ebiten.Image, p *entities.Player, activeBoss *entities.BossEvent, totalDistance float64) {
	// 1. Barra Superior Translúcida
	vector.DrawFilledRect(screen, 0, 0, float32(config.VirtualWidth), 34, color.RGBA{15, 20, 28, 220}, false)
	// Borda amarela inferior
	vector.DrawFilledRect(screen, 0, 33, float32(config.VirtualWidth), 2, color.RGBA{255, 230, 0, 255}, false)

	// 2. Estrelas de Reputação
	stars := p.GetStarRating()
	starStr := ""
	for i := 0; i < stars; i++ {
		starStr += "[*] "
	}
	for i := stars; i < 5; i++ {
		starStr += "[ ] "
	}

	badge := "OURO"
	if stars == 5 {
		badge = "PLATINUM 5★"
	} else if stars <= 2 {
		badge = "ALERTA ⚠️"
	}

	repColor := color.RGBA{255, 230, 0, 255}
	if stars <= 2 {
		repColor = color.RGBA{255, 60, 60, 255}
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("REPUTACAO: %s (%s)", starStr, badge), 12, 10)

	// Barra de Reputação Gráfica
	barWidth := float32(80)
	fillWidth := barWidth * float32(p.Reputation/100.0)
	vector.DrawFilledRect(screen, 240, 10, barWidth, 12, color.RGBA{40, 45, 55, 255}, false)
	vector.DrawFilledRect(screen, 240, 10, fillWidth, 12, repColor, false)

	// 3. Carga de Pacotes no Baú / Mochila
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("PACOTES: %d/%d", p.Cargo, p.MaxCargo), 340, 10)

	// 4. Pontuação e Combo
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("PONTOS: %05d", p.Score), 460, 10)
	if p.Combo > 1 {
		comboTxt := fmt.Sprintf("COMBO %dx!", p.Combo)
		ebitenutil.DebugPrintAt(screen, comboTxt, 570, 10)
	}

	// 5. Mini-Mapa da Rota (Rodapé)
	mapY := float32(config.VirtualHeight - 20)
	mapWidth := float32(360)
	mapStartX := float32(config.VirtualWidth/2.0) - (mapWidth / 2.0)

	vector.DrawFilledRect(screen, mapStartX-4, mapY-3, mapWidth+8, 14, color.RGBA{15, 20, 28, 200}, false)
	vector.DrawFilledRect(screen, mapStartX, mapY+2, mapWidth, 4, color.RGBA{60, 65, 75, 255}, false)

	// Marcadores dos 5 Bosses no minimapa
	bossCheckpoints := []float64{650, 1350, 2100, 2800, 3450}
	for _, bDist := range bossCheckpoints {
		normPos := float32(bDist / totalDistance)
		bx := mapStartX + normPos*mapWidth
		vector.DrawFilledRect(screen, bx-2, mapY-1, 4, 10, color.RGBA{255, 70, 50, 255}, false)
	}

	// Posição Atual do Entregador
	playerNorm := float32(p.Y / totalDistance)
	if playerNorm > 1.0 {
		playerNorm = 1.0
	}
	px := mapStartX + playerNorm*mapWidth
	vector.DrawFilledRect(screen, px-3, mapY-2, 7, 12, color.RGBA{255, 230, 0, 255}, false)

	distKm := p.Y / 1000.0
	totalKm := totalDistance / 1000.0
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("ROTA: %.1f / %.1f km", distKm, totalKm), int(mapStartX)-110, int(mapY)-1)

	// 6. Alerta de Evento de Chefe (Pânico ou Recuperação)
	if activeBoss != nil {
		if activeBoss.State == entities.BossPanic {
			// Banner de Pânico Vermelho Pulsante
			vector.DrawFilledRect(screen, 0, 40, float32(config.VirtualWidth), 38, color.RGBA{220, 30, 30, 220}, false)
			vector.DrawFilledRect(screen, 0, 76, float32(config.VirtualWidth), 2, color.RGBA{255, 255, 255, 255}, false)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("⚠️ %s", activeBoss.Name), 20, 44)
			ebitenutil.DebugPrintAt(screen, activeBoss.PanicDescription, 20, 58)
		} else if activeBoss.State == entities.BossRecover {
			// Banner de Recuperação Heroica Dourado/Azul
			vector.DrawFilledRect(screen, 0, 40, float32(config.VirtualWidth), 38, color.RGBA{25, 95, 215, 230}, false)
			vector.DrawFilledRect(screen, 0, 76, float32(config.VirtualWidth), 2, color.RGBA{255, 230, 0, 255}, false)
			recText := fmt.Sprintf("HEROISMO EXPRESS (%d/%d): %s", activeBoss.RecoverProgress, activeBoss.RecoverTarget, activeBoss.RecoverDescription)
			ebitenutil.DebugPrintAt(screen, recText, 20, 46)
			ebitenutil.DebugPrintAt(screen, "Pressione [ESPACO] para arremessar ou [H] buzinar!", 20, 60)
		}
	}
}

// DrawCustomizerUI desenha a tela da Central de Customização do Entregador.
func DrawCustomizerUI(screen *ebiten.Image, c *customizer.Customization, selCategory int) {
	// Fundo escuro com padrão azul/amarelo
	vector.DrawFilledRect(screen, 0, 0, float32(config.VirtualWidth), float32(config.VirtualHeight), color.RGBA{20, 24, 34, 255}, false)

	// Cabeçalho
	vector.DrawFilledRect(screen, 0, 0, float32(config.VirtualWidth), 40, color.RGBA{255, 230, 0, 255}, false)
	vector.DrawFilledRect(screen, 0, 38, float32(config.VirtualWidth), 3, color.RGBA{30, 95, 220, 255}, false)
	ebitenutil.DebugPrintAt(screen, "BOXBOY EXPRESS - CENTRAL DO ENTREGADOR", 175, 14)

	// Painel de Opções à Direita
	panelX := float32(280)
	panelY := float32(56)
	panelW := float32(340)
	panelH := float32(250)

	vector.DrawFilledRect(screen, panelX, panelY, panelW, panelH, color.RGBA{30, 36, 50, 230}, false)
	vector.DrawFilledRect(screen, panelX, panelY, panelW, 2, color.RGBA{255, 230, 0, 255}, false)

	categories := []string{
		"1. Tom de Pele",
		"2. Estilo de Cabelo",
		"3. Cor do Cabelo",
		"4. Uniforme / Roupa",
		"5. Boné / Capacete",
		"6. Óculos / Acessório",
		"7. Mascote do Cesto",
		"8. Veículo de Entrega",
	}

	for i, cat := range categories {
		itemY := int(panelY) + 14 + i*28

		if i == selCategory {
			// Item selecionado em destaque amarelo
			vector.DrawFilledRect(screen, panelX+4, float32(itemY-3), panelW-8, 22, color.RGBA{255, 230, 0, 190}, false)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("> %s", cat), int(panelX)+12, itemY)
		} else {
			ebitenutil.DebugPrintAt(screen, cat, int(panelX)+12, itemY)
		}

		if i == 0 {
			// Tom de Pele: seleção visual por quadradinhos de cor, sem rótulos.
			drawSwatchRow(screen, int(panelX)+160, itemY, customizer.SkinColors, c.SkinTone)
		} else {
			valTxt := c.GetCurrentValueText(i)
			if i == selCategory {
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("< %s >", valTxt), int(panelX)+160, itemY)
			} else {
				ebitenutil.DebugPrintAt(screen, valTxt, int(panelX)+160, itemY)
			}
		}
	}

	// Painel de Preview à Esquerda
	previewX := float32(30)
	previewY := float32(56)
	previewW := float32(230)
	previewH := float32(250)
	vector.DrawFilledRect(screen, previewX, previewY, previewW, previewH, color.RGBA{25, 30, 42, 230}, false)
	vector.DrawFilledRect(screen, previewX, previewY, previewW, 2, color.RGBA{30, 95, 220, 255}, false)
	ebitenutil.DebugPrintAt(screen, "PREVIEW DO ENTREGADOR", int(previewX)+35, int(previewY)+10)

	// Rodapé com instruções
	vector.DrawFilledRect(screen, 0, float32(config.VirtualHeight-42), float32(config.VirtualWidth), 42, color.RGBA{15, 18, 25, 255}, false)
	ebitenutil.DebugPrintAt(screen, "[W/S ou CIMA/BAIXO]: Navegar Categorias | [A/D ou ESQ/DIR]: Trocar Opcao", 70, config.VirtualHeight-32)
	ebitenutil.DebugPrintAt(screen, "Pressione [ESPACO] ou [ENTER] para CONFIRMAR E INICIAR A ROTA!", 110, config.VirtualHeight-16)
}

// drawSwatchRow desenha uma fileira de quadradinhos de cor selecionáveis,
// destacando com uma borda amarela o índice atualmente selecionado.
func drawSwatchRow(screen *ebiten.Image, x, y int, colors []color.RGBA, selected int) {
	const size, gap = 14, 5
	for i, col := range colors {
		sx := float32(x + i*(size+gap))
		sy := float32(y - 2)

		if i == selected {
			// Borda de destaque um pouco maior por trás do quadrado.
			vector.DrawFilledRect(screen, sx-2, sy-2, size+4, size+4, color.RGBA{255, 230, 0, 255}, false)
		}
		vector.DrawFilledRect(screen, sx, sy, size, size, col, false)
	}
}

// DrawGameOverUI desenha a tela de fim de rota quando o jogador zera a reputação.
func DrawGameOverUI(screen *ebiten.Image, p *entities.Player) {
	vector.DrawFilledRect(screen, 0, 0, float32(config.VirtualWidth), float32(config.VirtualHeight), color.RGBA{15, 10, 12, 235}, false)

	vector.DrawFilledRect(screen, 120, 60, 400, 240, color.RGBA{28, 20, 24, 250}, false)
	vector.DrawFilledRect(screen, 120, 60, 400, 3, color.RGBA{255, 50, 50, 255}, false)

	ebitenutil.DebugPrintAt(screen, "ROTA CANCELADA - REPUTACAO ESGOTADA!", 180, 80)
	ebitenutil.DebugPrintAt(screen, "Os clientes reclamaram de pacotes danificados e atrasos.", 145, 110)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Entregas com Sucesso: %d", p.Successful), 160, 145)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Pontuacao Final: %d", p.Score), 160, 165)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Maior Sequencia Combo: %dx", p.MaxCombo), 160, 185)

	vector.DrawFilledRect(screen, 200, 235, 240, 30, color.RGBA{255, 230, 0, 255}, false)
	ebitenutil.DebugPrintAt(screen, "PRESSIONE [ESPACO] PARA NOVO TURNO", 215, 244)
}

// DrawVictoryUI desenha a celebração de rota concluída com sucesso e 5 estrelas.
func DrawVictoryUI(screen *ebiten.Image, p *entities.Player) {
	vector.DrawFilledRect(screen, 0, 0, float32(config.VirtualWidth), float32(config.VirtualHeight), color.RGBA{10, 18, 25, 235}, false)

	vector.DrawFilledRect(screen, 110, 45, 420, 270, color.RGBA{22, 32, 48, 250}, false)
	vector.DrawFilledRect(screen, 110, 45, 420, 4, color.RGBA{255, 230, 0, 255}, false)

	ebitenutil.DebugPrintAt(screen, "PARABENS! ROTA EXPRESS CONCLUIDA!", 175, 65)
	ebitenutil.DebugPrintAt(screen, "CLASSIFICACAO: ENTREGADOR NOTA 10 (RANK PLATINUM)", 135, 95)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Entregas Perfeitas: %d", p.Successful), 160, 130)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Chefes e Catastrofes Superadas: %d/5", p.BossesBeaten), 160, 150)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Pontuacao Turbo: %d Pts", p.Score), 160, 170)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Reputacao Mantida: %.1f%%", p.Reputation), 160, 190)

	vector.DrawFilledRect(screen, 190, 240, 260, 32, color.RGBA{255, 230, 0, 255}, false)
	ebitenutil.DebugPrintAt(screen, "PRESSIONE [ESPACO] PARA RECOMEÇAR", 210, 250)
}
