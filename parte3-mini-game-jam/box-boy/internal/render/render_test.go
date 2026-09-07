package render_test

import (
	"testing"

	"box-boy/internal/render"
)

func TestIsoProject(t *testing.T) {
	// Ponto onde o jogador está na origem
	playerX, playerY := 0.0, 0.0

	// 1. Posição na origem (0, 0, 0)
	sx1, sy1 := render.IsoProject(0, 0, 0, playerX, playerY)

	// 2. Avanço para frente (+Y no mundo): deve ir para cima (-Y na tela) e para a direita (+X na tela)
	sx2, sy2 := render.IsoProject(0, 100.0, 0, playerX, playerY)
	if sx2 <= sx1 {
		t.Errorf("avanço na rota deveria aumentar X na tela (diagonal), sx1=%f, sx2=%f", sx1, sx2)
	}
	if sy2 >= sy1 {
		t.Errorf("avanço na rota deveria diminuir Y na tela (para cima), sy1=%f, sy2=%f", sy1, sy2)
	}

	// 3. Salto (Bunny-Hop / Z positivo): deve mover para cima (-Y na tela) sem alterar X
	sx3, sy3 := render.IsoProject(0, 0, 50.0, playerX, playerY)
	if sx3 != sx1 {
		t.Errorf("altura Z não deve afetar X da tela: sx1=%f, sx3=%f", sx1, sx3)
	}
	if sy3 >= sy1 {
		t.Errorf("altura Z deve deslocar para cima (-Y): sy1=%f, sy3=%f", sy1, sy3)
	}
}

func TestIsOnScreen(t *testing.T) {
	if !render.IsOnScreen(320, 180, 0) {
		t.Errorf("centro da tela (320, 180) deve estar visível")
	}
	if render.IsOnScreen(-100, -100, 20) {
		t.Errorf("(-100, -100) não deve estar visível com margem 20")
	}
}
