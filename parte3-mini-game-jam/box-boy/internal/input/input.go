package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// InputState armazena as intenções lógicas de entrada do jogador no frame atual.
type InputState struct {
	MoveX         float64 // -1.0 (esquerda) a +1.0 (direita)
	Accelerate    bool
	Brake         bool
	JustThrew     bool // Arremesso de pacote
	JustJumped    bool // Bunny-Hop
	JustHorn      bool // Buzina para cães/pedestres
	JustRecovered bool // Ação heroica de recuperação de chefe
	JustSelected  bool // Seleção em menus
	JustBack      bool // Voltar / Pausar
	NavX          int  // Navegação discreta de menu (-1, 0, 1)
	NavY          int  // Navegação vertical (-1, 0, 1)
}

// PollInputs lê os estados de teclado e gamepad mapeando para ações lógicas.
func PollInputs() InputState {
	var in InputState

	// 1. Movimento Lateral
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		in.MoveX -= 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		in.MoveX += 1.0
	}

	// 2. Acelerar e Frear
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		in.Accelerate = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		in.Brake = true
	}

	// 3. Arremesso de Pacote (Espaço, J ou Clique Esquerdo)
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsKeyJustPressed(ebiten.KeyJ) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		in.JustThrew = true
	}

	// 4. Bunny-Hop (Shift Esquerdo, K ou Espaço enquanto pula)
	if inpututil.IsKeyJustPressed(ebiten.KeyShift) ||
		inpututil.IsKeyJustPressed(ebiten.KeyShiftLeft) ||
		inpututil.IsKeyJustPressed(ebiten.KeyK) {
		in.JustJumped = true
	}

	// 5. Buzina (H, E ou B)
	if inpututil.IsKeyJustPressed(ebiten.KeyH) ||
		inpututil.IsKeyJustPressed(ebiten.KeyE) ||
		inpututil.IsKeyJustPressed(ebiten.KeyB) {
		in.JustHorn = true
	}

	// 6. Ação Heroica de Recuperação de Chefe (R, Enter ou Espaço)
	if inpututil.IsKeyJustPressed(ebiten.KeyR) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		in.JustRecovered = true
	}

	// 7. Navegação de Menu
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		in.NavX = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		in.NavX = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		in.NavY = -1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		in.NavY = 1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		in.JustSelected = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyP) {
		in.JustBack = true
	}

	// Suporte Gamepad (se conectado)
	gamepadIDs := ebiten.AppendGamepadIDs(nil)
	if len(gamepadIDs) > 0 {
		gid := gamepadIDs[0]
		// Eixos analógicos
		axisX := ebiten.GamepadAxisValue(gid, 0)
		if axisX < -0.2 {
			in.MoveX = -1.0
		} else if axisX > 0.2 {
			in.MoveX = 1.0
		}

		if ebiten.IsGamepadButtonPressed(gid, ebiten.GamepadButton1) { // RT / Acelerar
			in.Accelerate = true
		}
		if ebiten.IsGamepadButtonPressed(gid, ebiten.GamepadButton2) { // LT / Freio
			in.Brake = true
		}
		if inpututil.IsGamepadButtonJustPressed(gid, ebiten.GamepadButton0) { // Botão A (Sul)
			in.JustThrew = true
			in.JustSelected = true
		}
		if inpututil.IsGamepadButtonJustPressed(gid, ebiten.GamepadButton3) { // Botão X (Oeste)
			in.JustJumped = true
		}
		if inpututil.IsGamepadButtonJustPressed(gid, ebiten.GamepadButton1) { // Botão B (Leste)
			in.JustHorn = true
		}
		if inpututil.IsGamepadButtonJustPressed(gid, ebiten.GamepadButton2) { // Botão Y (Norte)
			in.JustRecovered = true
		}
	}

	return in
}
