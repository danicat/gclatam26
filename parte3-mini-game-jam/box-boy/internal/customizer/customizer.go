package customizer

import "image/color"

// Customization armazena todas as opções selecionadas pelo jogador para o entregador e veículo.
type Customization struct {
	SkinTone    int
	HairStyle   int
	HairColor   int
	Outfit      int
	Headgear    int
	Glasses     int
	Companion   int
	VehicleType int
	PlayerName  string
}

// Opções e labels para a interface de customização
var (
	SkinToneNames = []string{
		"Tom 1",
		"Tom 2",
		"Tom 3",
		"Tom 4",
		"Tom 5",
	}

	SkinColors = []color.RGBA{
		{255, 218, 185, 255},
		{238, 175, 125, 255},
		{184, 125, 75, 255},
		{110, 65, 38, 255},
		{225, 175, 95, 255},
	}

	HairStyleNames = []string{
		"Degradê Moderno",
		"Coque Samurai",
		"Black Power Afro",
		"Moicano Estiloso",
		"Dreadlocks",
		"Cabelo Longo Solto",
	}

	HairColorNames = []string{
		"Castanho Escuro",
		"Loiro Dourado",
		"Ruivo Fogo",
		"Azul Neon",
		"Rosa Choque Cyber",
		"Prata Futurista",
	}

	HairColors = []color.RGBA{
		{35, 25, 20, 255},    // Castanho escuro
		{235, 195, 80, 255},  // Loiro
		{195, 55, 30, 255},   // Ruivo
		{30, 110, 240, 255},  // Azul Neon
		{245, 60, 160, 255},  // Rosa Choque
		{220, 225, 235, 255}, // Prata
	}

	OutfitNames = []string{
		"Jaqueta Corta-Vento (Amarela)",
		"Colete Refletivo Fluorescente",
		"Moletom Streetwear Express",
		"Camisa Polo Express Amarela",
	}

	HeadgearNames = []string{
		"Boné Aba Reta (Virado)",
		"Viseira Retrô Amarela",
		"Capacete Aerodinâmico",
		"Bandana Esportiva Azul",
		"Nenhum (Sem Acessório)",
	}

	GlassesNames = []string{
		"Nenhum",
		"Óculos Ciclista Neon",
		"Óculos de Grau Estiloso",
		"Óculos Escuros Aviador",
	}

	CompanionNames = []string{
		"Cão Caramelo (Mascote Oficial)",
		"Capivara Zen",
		"Mini Drone Rastreador",
	}

	VehicleNames = []string{
		"Bicicleta Urbana Amarela (Bunny-Hop Ágil)",
		"Scooter Elétrica Express (Arrancada Rápida)",
		"Furgão Compacto Elétrico (Alta Carga)",
	}
)

// NewDefaultCustomization retorna o perfil inicial padrão do BoxBoy
func NewDefaultCustomization() Customization {
	return Customization{
		SkinTone:    1,
		HairStyle:   0, // Degradê
		HairColor:   0, // Castanho
		Outfit:      0, // Jaqueta Corta-Vento
		Headgear:    0, // Boné Aba Reta
		Glasses:     1, // Óculos Neon
		Companion:   0, // Cão Caramelo
		VehicleType: 0, // Bicicleta Amarela
		PlayerName:  "BoxBoy",
	}
}

// NextOption avança o índice da categoria especificada com loop circular
func (c *Customization) NextOption(categoryIndex int) {
	switch categoryIndex {
	case 0:
		c.SkinTone = (c.SkinTone + 1) % len(SkinToneNames)
	case 1:
		c.HairStyle = (c.HairStyle + 1) % len(HairStyleNames)
	case 2:
		c.HairColor = (c.HairColor + 1) % len(HairColorNames)
	case 3:
		c.Outfit = (c.Outfit + 1) % len(OutfitNames)
	case 4:
		c.Headgear = (c.Headgear + 1) % len(HeadgearNames)
	case 5:
		c.Glasses = (c.Glasses + 1) % len(GlassesNames)
	case 6:
		c.Companion = (c.Companion + 1) % len(CompanionNames)
	case 7:
		c.VehicleType = (c.VehicleType + 1) % len(VehicleNames)
	}
}

// PrevOption retrocede o índice da categoria especificada
func (c *Customization) PrevOption(categoryIndex int) {
	switch categoryIndex {
	case 0:
		c.SkinTone = (c.SkinTone - 1 + len(SkinToneNames)) % len(SkinToneNames)
	case 1:
		c.HairStyle = (c.HairStyle - 1 + len(HairStyleNames)) % len(HairStyleNames)
	case 2:
		c.HairColor = (c.HairColor - 1 + len(HairColorNames)) % len(HairColorNames)
	case 3:
		c.Outfit = (c.Outfit - 1 + len(OutfitNames)) % len(OutfitNames)
	case 4:
		c.Headgear = (c.Headgear - 1 + len(HeadgearNames)) % len(HeadgearNames)
	case 5:
		c.Glasses = (c.Glasses - 1 + len(GlassesNames)) % len(GlassesNames)
	case 6:
		c.Companion = (c.Companion - 1 + len(CompanionNames)) % len(CompanionNames)
	case 7:
		c.VehicleType = (c.VehicleType - 1 + len(VehicleNames)) % len(VehicleNames)
	}
}

// GetCurrentValueText retorna a descrição do valor atual da categoria
func (c *Customization) GetCurrentValueText(categoryIndex int) string {
	switch categoryIndex {
	case 0:
		return SkinToneNames[c.SkinTone]
	case 1:
		return HairStyleNames[c.HairStyle]
	case 2:
		return HairColorNames[c.HairColor]
	case 3:
		return OutfitNames[c.Outfit]
	case 4:
		return HeadgearNames[c.Headgear]
	case 5:
		return GlassesNames[c.Glasses]
	case 6:
		return CompanionNames[c.Companion]
	case 7:
		return VehicleNames[c.VehicleType]
	default:
		return ""
	}
}
