package render

import (
	"box-boy/internal/config"
	"math"
)

// IsoProject projeta coordenadas 3D do mundo (X lateral, Y distância da rota, Z altura)
// para a tela 2D no estilo diagonal Paperboy (avanço em diagonal para cima e direita).
func IsoProject(worldX, worldY, worldZ float64, playerX, playerY float64) (float64, float64) {
	relX := worldX - playerX
	relY := worldY - playerY

	// Ângulo pseudo-isométrico diagonal clássico
	// Avanço em Y leva para a direita (+X tela) e para cima (-Y tela)
	screenX := (float64(config.VirtualWidth) * 0.42) + (relX * 1.05) + (relY * 0.72)
	screenY := (float64(config.VirtualHeight) * 0.68) + (relX * 0.42) - (relY * 0.48) - worldZ

	return screenX, screenY
}

// IsOnScreen verifica se um ponto projetado está visível no canvas virtual com margem.
func IsOnScreen(screenX, screenY float64, margin float64) bool {
	return screenX >= -margin &&
		screenX <= float64(config.VirtualWidth)+margin &&
		screenY >= -margin &&
		screenY <= float64(config.VirtualHeight)+margin
}

// Distance2D calcula a distância euclidiana simples entre dois pontos no plano mundo.
func Distance2D(x1, y1, x2, y2 float64) float64 {
	return math.Hypot(x2-x1, y2-y1)
}
