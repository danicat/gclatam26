package customizer_test

import (
	"testing"

	"box-boy/internal/customizer"
)

func TestCustomizationDefaults(t *testing.T) {
	c := customizer.NewDefaultCustomization()

	if c.PlayerName != "BoxBoy" {
		t.Errorf("esperava nome 'BoxBoy', obteve '%s'", c.PlayerName)
	}
	if c.VehicleType != 0 {
		t.Errorf("esperava veículo 0 (Bicicleta), obteve %d", c.VehicleType)
	}
	if c.Companion != 0 {
		t.Errorf("esperava mascote 0 (Cão Caramelo), obteve %d", c.Companion)
	}
}

func TestCustomizationCycling(t *testing.T) {
	c := customizer.NewDefaultCustomization()

	// Avançar e retroceder categoria 0 (Tom de Pele)
	initialSkin := c.SkinTone
	c.NextOption(0)
	if c.SkinTone == initialSkin {
		t.Errorf("esperava que o tom de pele mudasse após NextOption")
	}

	c.PrevOption(0)
	if c.SkinTone != initialSkin {
		t.Errorf("esperava voltar ao tom de pele inicial, obteve %d", c.SkinTone)
	}

	// Loop circular
	numSkins := len(customizer.SkinToneNames)
	for i := 0; i < numSkins; i++ {
		c.NextOption(0)
	}
	if c.SkinTone != initialSkin {
		t.Errorf("após %d rotações completas, esperava %d, obteve %d", numSkins, initialSkin, c.SkinTone)
	}
}

func TestOptionNamesAndColors(t *testing.T) {
	c := customizer.NewDefaultCustomization()

	txt := c.GetCurrentValueText(0)
	if txt == "" {
		t.Errorf("GetCurrentValueText(0) não deve ser vazio")
	}

	if len(customizer.SkinColors) != len(customizer.SkinToneNames) {
		t.Errorf("número de cores de pele (%d) diverge dos nomes (%d)", len(customizer.SkinColors), len(customizer.SkinToneNames))
	}
}
