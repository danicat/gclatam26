package audio

import (
	"bytes"
	"math"
	"math/rand"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// NoteDef define parâmetros acústicos para síntese de nota por software (DSP).
type NoteDef struct {
	WaveType     string  // "sine", "square", "triangle", "sawtooth", "noise"
	Duration     float64 // Duração em segundos
	StartFreq    float64 // Frequência inicial em Hz
	EndFreq      float64 // Frequência final em Hz
	VibratoFreq  float64 // Frequência de vibrato LFO em Hz
	VibratoDepth float64 // Profundidade do vibrato em Hz
	NoiseFilter  float64 // Filtro passa-baixa para ruído (0.0 a 1.0)
	DutyCycle    float64 // Ciclo de trabalho da onda quadrada (0.1 a 0.9)
	Pan          float64 // Pan estéreo (-1.0 esq a +1.0 dir)
	Volume       float64 // 0.0 a 1.0
	Attack       float64
	Decay        float64
	Sustain      float64
	Release      float64
}

// TrackDef é uma faixa com uma sequência de notas.
type TrackDef struct {
	Name  string
	Pan   float64
	Notes []NoteDef
}

// AudioSystem gerencia a síntese e reprodução de som do BoxBoy.
type AudioSystem struct {
	ctx *audio.Context
	mu  sync.Mutex

	// BGM Players
	bgmPlayer *audio.Player
	bgmGroove []byte
	bgmPanic  []byte
	bgmWin    []byte
	currTrack string

	// SFX Pre-sintetizados
	sfxThrow         []byte
	sfxDeliver       []byte
	sfxCombo         []byte
	sfxCrash         []byte
	sfxBunnyHop      []byte
	sfxBell          []byte
	sfxHorn          []byte
	sfxBark          []byte
	sfxPanicSiren    []byte
	sfxRecoverAction []byte
	sfxTornado       []byte
	sfxClick         []byte

	volume float64
}

// NewAudioSystem cria o sistema de áudio e sintetiza as faixas e efeitos em memória.
func NewAudioSystem() *AudioSystem {
	ctx := audio.NewContext(44100)
	sys := &AudioSystem{
		ctx:    ctx,
		volume: 0.85,
	}

	sys.buildSFX()
	sys.buildBGM()

	return sys
}

// GenerateNotePCM sintetiza uma nota pura em PCM estéreo de 16 bits a 44.100 Hz.
func GenerateNotePCM(n NoteDef) []byte {
	const sampleRate = 44100
	numSamples := int(n.Duration * sampleRate)
	if numSamples <= 0 {
		return nil
	}

	buf := make([]byte, numSamples*4)
	var phase float64
	var filterVal float64

	pan := math.Max(-1.0, math.Min(1.0, n.Pan))
	leftPan := math.Cos((pan + 1.0) * math.Pi / 4.0)
	rightPan := math.Sin((pan + 1.0) * math.Pi / 4.0)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		progress := t / n.Duration

		var currentFreq float64
		if n.StartFreq == n.EndFreq {
			currentFreq = n.StartFreq
		} else if n.StartFreq > 0 && n.EndFreq > 0 {
			currentFreq = n.StartFreq * math.Pow(n.EndFreq/n.StartFreq, progress)
		} else {
			currentFreq = n.StartFreq + progress*(n.EndFreq-n.StartFreq)
		}

		if n.VibratoFreq > 0 && n.VibratoDepth > 0 {
			currentFreq += math.Sin(2.0*math.Pi*n.VibratoFreq*t) * n.VibratoDepth
		}

		phase += currentFreq / sampleRate
		for phase >= 1.0 {
			phase -= 1.0
		}

		var oscAmp float64
		switch n.WaveType {
		case "sine":
			oscAmp = math.Sin(2.0 * math.Pi * phase)
		case "square":
			duty := n.DutyCycle
			if duty <= 0 || duty >= 1.0 {
				duty = 0.5
			}
			if phase < duty {
				oscAmp = 0.85
			} else {
				oscAmp = -0.85
			}
		case "triangle":
			if phase < 0.5 {
				oscAmp = 4.0*phase - 1.0
			} else {
				oscAmp = 3.0 - 4.0*phase
			}
		case "sawtooth":
			oscAmp = 2.0*phase - 1.0
		case "noise":
			white := (rand.Float64()*2.0 - 1.0) * 0.9
			alpha := n.NoiseFilter
			if alpha <= 0 || alpha > 1.0 {
				alpha = 0.35
			}
			filterVal = filterVal + alpha*(white-filterVal)
			oscAmp = filterVal
		default:
			oscAmp = math.Sin(2.0 * math.Pi * phase)
		}

		// Envelope ADSR
		var envAmp float64
		attTime := n.Attack
		decTime := n.Decay
		susLevel := n.Sustain
		relTime := n.Release

		totalTime := n.Duration
		susTime := totalTime - attTime - decTime - relTime
		if susTime < 0 {
			ratio := totalTime / (attTime + decTime + relTime)
			attTime *= ratio
			decTime *= ratio
			relTime *= ratio
			susTime = 0
		}

		if t < attTime {
			if attTime > 0 {
				envAmp = t / attTime
			} else {
				envAmp = 1.0
			}
		} else if t < attTime+decTime {
			if decTime > 0 {
				decayProg := (t - attTime) / decTime
				envAmp = 1.0 - decayProg*(1.0-susLevel)
			} else {
				envAmp = susLevel
			}
		} else if t < attTime+decTime+susTime {
			envAmp = susLevel
		} else {
			relProg := (t - (attTime + decTime + susTime)) / relTime
			envAmp = math.Max(0.0, susLevel*(1.0-relProg))
		}

		sampleVal := oscAmp * envAmp * n.Volume
		if sampleVal > 1.0 {
			sampleVal = 1.0
		} else if sampleVal < -1.0 {
			sampleVal = -1.0
		}

		leftSample := int16(sampleVal * leftPan * 32767.0)
		rightSample := int16(sampleVal * rightPan * 32767.0)

		idx := i * 4
		buf[idx] = byte(leftSample)
		buf[idx+1] = byte(leftSample >> 8)
		buf[idx+2] = byte(rightSample)
		buf[idx+3] = byte(rightSample >> 8)
	}

	return buf
}

// MixTracksPCM mixa múltiplas faixas PCM em um único buffer estéreo.
func MixTracksPCM(tracks [][]byte) []byte {
	maxLen := 0
	for _, t := range tracks {
		if len(t) > maxLen {
			maxLen = len(t)
		}
	}
	result := make([]byte, maxLen)

	for i := 0; i < maxLen; i += 2 {
		var sum int32
		for _, t := range tracks {
			if i+1 < len(t) {
				sample := int16(uint16(t[i]) | (uint16(t[i+1]) << 8))
				sum += int32(sample)
			}
		}
		if sum > 32767 {
			sum = 32767
		} else if sum < -32768 {
			sum = -32768
		}
		result[i] = byte(sum)
		result[i+1] = byte(sum >> 8)
	}
	return result
}

// buildSFX sintetiza todos os efeitos sonoros temáticos do jogo.
func (s *AudioSystem) buildSFX() {
	// 1. Arremesso de Pacote: Swoosh rápido + impacto macio de papelão
	w := GenerateNotePCM(NoteDef{
		WaveType: "noise", Duration: 0.18, NoiseFilter: 0.7,
		Volume: 0.28, Attack: 0.01, Decay: 0.08, Sustain: 0.3, Release: 0.09, Pan: -0.2,
	})
	s.sfxThrow = w

	// 2. Entrega com Sucesso: Campainha Ding-Dong harmônica (Mi -> Dó) + Chime cintilante
	d1 := GenerateNotePCM(NoteDef{
		WaveType: "sine", Duration: 0.22, StartFreq: 659.25, EndFreq: 659.25,
		Volume: 0.32, Attack: 0.005, Decay: 0.08, Sustain: 0.4, Release: 0.13, Pan: -0.1,
	})
	d2 := GenerateNotePCM(NoteDef{
		WaveType: "sine", Duration: 0.35, StartFreq: 523.25, EndFreq: 523.25,
		Volume: 0.35, Attack: 0.005, Decay: 0.12, Sustain: 0.5, Release: 0.22, Pan: 0.1,
	})
	sparkle := GenerateNotePCM(NoteDef{
		WaveType: "triangle", Duration: 0.4, StartFreq: 1046.5, EndFreq: 1318.5,
		VibratoFreq: 12.0, VibratoDepth: 18.0, Volume: 0.18, Attack: 0.01, Decay: 0.1, Sustain: 0.4, Release: 0.2, Pan: 0.3,
	})
	dingDong := append(d1, d2...)
	s.sfxDeliver = MixTracksPCM([][]byte{dingDong, sparkle})

	// 3. Combo Perfeito: Sequência rápida ascendente de acordes dourados
	c1 := GenerateNotePCM(NoteDef{WaveType: "triangle", Duration: 0.08, StartFreq: 587.33, EndFreq: 587.33, Volume: 0.25, Attack: 0.005, Decay: 0.03, Sustain: 0.5, Release: 0.04})
	c2 := GenerateNotePCM(NoteDef{WaveType: "triangle", Duration: 0.08, StartFreq: 739.99, EndFreq: 739.99, Volume: 0.25, Attack: 0.005, Decay: 0.03, Sustain: 0.5, Release: 0.04})
	c3 := GenerateNotePCM(NoteDef{WaveType: "triangle", Duration: 0.18, StartFreq: 880.00, EndFreq: 880.00, Volume: 0.3, Attack: 0.005, Decay: 0.05, Sustain: 0.6, Release: 0.12})
	s.sfxCombo = append(c1, append(c2, c3...)...)

	// 4. Quebra de Vidro / Colisão Cômica
	glassNoise := GenerateNotePCM(NoteDef{WaveType: "noise", Duration: 0.28, NoiseFilter: 0.85, Volume: 0.35, Attack: 0.002, Decay: 0.1, Sustain: 0.3, Release: 0.17})
	glassCrack := GenerateNotePCM(NoteDef{WaveType: "sawtooth", Duration: 0.32, StartFreq: 980, EndFreq: 110, Volume: 0.28, Attack: 0.005, Decay: 0.12, Sustain: 0.3, Release: 0.19})
	s.sfxCrash = MixTracksPCM([][]byte{glassNoise, glassCrack})

	// 5. Bunny-Hop: Som de mola/salto acrobático da bicicleta
	s.sfxBunnyHop = GenerateNotePCM(NoteDef{
		WaveType: "sine", Duration: 0.26, StartFreq: 180, EndFreq: 460,
		Volume: 0.25, Attack: 0.01, Decay: 0.08, Sustain: 0.5, Release: 0.17, Pan: 0.0,
	})

	// 6. Campainha de Bicicleta: Triiim-Triiim clássico
	b1 := GenerateNotePCM(NoteDef{WaveType: "sine", Duration: 0.09, StartFreq: 1760, EndFreq: 1760, Volume: 0.26, Attack: 0.002, Decay: 0.03, Sustain: 0.4, Release: 0.05})
	b2 := GenerateNotePCM(NoteDef{WaveType: "sine", Duration: 0.14, StartFreq: 2093, EndFreq: 2093, Volume: 0.28, Attack: 0.002, Decay: 0.04, Sustain: 0.4, Release: 0.09})
	s.sfxBell = append(b1, b2...)

	// 7. Buzina de Scooter: Bi-Bi urbano estridente
	h1 := GenerateNotePCM(NoteDef{WaveType: "square", Duration: 0.08, StartFreq: 520, EndFreq: 520, DutyCycle: 0.3, Volume: 0.22, Attack: 0.005, Decay: 0.02, Sustain: 0.6, Release: 0.05})
	h2 := GenerateNotePCM(NoteDef{WaveType: "square", Duration: 0.12, StartFreq: 620, EndFreq: 620, DutyCycle: 0.3, Volume: 0.24, Attack: 0.005, Decay: 0.02, Sustain: 0.6, Release: 0.09})
	s.sfxHorn = append(h1, h2...)

	// 8. Latido do Cão de Portão: Duplo au-au cômico
	bark1 := GenerateNotePCM(NoteDef{WaveType: "sawtooth", Duration: 0.11, StartFreq: 240, EndFreq: 140, Volume: 0.32, Attack: 0.01, Decay: 0.04, Sustain: 0.5, Release: 0.06})
	bark2 := GenerateNotePCM(NoteDef{WaveType: "sawtooth", Duration: 0.15, StartFreq: 270, EndFreq: 130, Volume: 0.34, Attack: 0.01, Decay: 0.05, Sustain: 0.5, Release: 0.09})
	s.sfxBark = append(bark1, bark2...)

	// 9. Sirene de Pânico (Catástrofe de Boss): Alarme oscilante agudo
	s.sfxPanicSiren = GenerateNotePCM(NoteDef{
		WaveType: "sawtooth", Duration: 0.65, StartFreq: 750, EndFreq: 1150,
		VibratoFreq: 7.0, VibratoDepth: 180.0, Volume: 0.3, Attack: 0.05, Decay: 0.1, Sustain: 0.7, Release: 0.15,
	})

	// 10. Ação de Recuperação Heroica: Acorde triunfante com power-up
	rec1 := GenerateNotePCM(NoteDef{WaveType: "triangle", Duration: 0.35, StartFreq: 330, EndFreq: 660, Volume: 0.28, Attack: 0.01, Decay: 0.08, Sustain: 0.6, Release: 0.2})
	rec2 := GenerateNotePCM(NoteDef{WaveType: "sine", Duration: 0.45, StartFreq: 523, EndFreq: 1046, VibratoFreq: 14, VibratoDepth: 25, Volume: 0.25, Attack: 0.02, Decay: 0.1, Sustain: 0.6, Release: 0.25})
	s.sfxRecoverAction = MixTracksPCM([][]byte{rec1, rec2})

	// 11. Rugido de Vento de Tornado: Ruído turbilhonante com modulação
	s.sfxTornado = GenerateNotePCM(NoteDef{
		WaveType: "noise", Duration: 0.8, NoiseFilter: 0.45,
		VibratoFreq: 5.0, VibratoDepth: 0.2, Volume: 0.35, Attack: 0.1, Decay: 0.2, Sustain: 0.7, Release: 0.3,
	})

	// 12. Clique de Interface
	s.sfxClick = GenerateNotePCM(NoteDef{
		WaveType: "triangle", Duration: 0.04, StartFreq: 900, EndFreq: 450,
		Volume: 0.2, Attack: 0.002, Decay: 0.01, Sustain: 0.3, Release: 0.02,
	})
}

// buildBGM compõe as faixas de música completas em loop.
func (s *AudioSystem) buildBGM() {
	// =========================================================================
	// TEMA 1: "Turbo Delivery Groove" (Funk/Synthwave Enérgico - 128 BPM)
	// =========================================================================
	beat := 60.0 / 128.0 // ~0.46875s
	halfBeat := beat / 2.0
	quarterBeat := beat / 4.0

	// Faixa 1: Slap Bassline Dançante (Notas em D, G, A, C)
	bassFreqs := []float64{
		146.83, 146.83, 174.61, 196.00, // D, D, F, G
		196.00, 220.00, 196.00, 146.83, // G, A, G, D
		130.81, 130.81, 164.81, 196.00, // C, C, E, G
		220.00, 261.63, 220.00, 146.83, // A, C, A, D
	}
	var bassTrack []byte
	for _, f := range bassFreqs {
		note := GenerateNotePCM(NoteDef{
			WaveType: "square", Duration: halfBeat, StartFreq: f, EndFreq: f * 0.96,
			DutyCycle: 0.35, Volume: 0.24, Attack: 0.005, Decay: 0.08, Sustain: 0.4, Release: 0.08, Pan: -0.1,
		})
		bassTrack = append(bassTrack, note...)
	}

	// Faixa 2: Lead Brass Alegre de Entrega
	leadNotes := []struct {
		freq float64
		dur  float64
	}{
		{587.33, halfBeat}, {659.25, halfBeat}, {783.99, beat},
		{880.00, halfBeat}, {783.99, halfBeat}, {587.33, beat},
		{523.25, halfBeat}, {659.25, halfBeat}, {783.99, halfBeat}, {880.00, halfBeat},
		{1046.50, beat}, {880.00, beat},
	}
	var leadTrack []byte
	for _, ln := range leadNotes {
		note := GenerateNotePCM(NoteDef{
			WaveType: "sawtooth", Duration: ln.dur, StartFreq: ln.freq, EndFreq: ln.freq,
			VibratoFreq: 6.0, VibratoDepth: 4.0, Volume: 0.16, Attack: 0.01, Decay: 0.08, Sustain: 0.6, Release: 0.1, Pan: 0.2,
		})
		leadTrack = append(leadTrack, note...)
	}

	// Faixa 3: Bateria Rítmica (Kick, Snare, Hi-Hats em 16th)
	var drumTrack []byte
	numMeasures := 4
	for m := 0; m < numMeasures; m++ {
		for b := 0; b < 4; b++ {
			// Bumbo no tempo 0 e 2, Caixa no 1 e 3
			var mainDrum []byte
			if b%2 == 0 {
				mainDrum = GenerateNotePCM(NoteDef{
					WaveType: "sine", Duration: halfBeat, StartFreq: 130, EndFreq: 45,
					Volume: 0.32, Attack: 0.002, Decay: 0.08, Sustain: 0.2, Release: 0.08,
				})
			} else {
				snareTone := GenerateNotePCM(NoteDef{
					WaveType: "triangle", Duration: halfBeat, StartFreq: 220, EndFreq: 90,
					Volume: 0.18, Attack: 0.002, Decay: 0.06, Sustain: 0.2, Release: 0.06,
				})
				snareNoise := GenerateNotePCM(NoteDef{
					WaveType: "noise", Duration: halfBeat, NoiseFilter: 0.8,
					Volume: 0.25, Attack: 0.002, Decay: 0.08, Sustain: 0.3, Release: 0.08,
				})
				mainDrum = MixTracksPCM([][]byte{snareTone, snareNoise})
			}

			// 2 Hi-Hats rápidos por batida
			hihat1 := GenerateNotePCM(NoteDef{
				WaveType: "noise", Duration: quarterBeat, NoiseFilter: 0.95,
				Volume: 0.12, Attack: 0.001, Decay: 0.02, Sustain: 0.1, Release: 0.02, Pan: 0.3,
			})
			hihat2 := GenerateNotePCM(NoteDef{
				WaveType: "noise", Duration: quarterBeat, NoiseFilter: 0.95,
				Volume: 0.08, Attack: 0.001, Decay: 0.02, Sustain: 0.1, Release: 0.02, Pan: -0.3,
			})
			hihats := append(hihat1, hihat2...)

			mixedBeat := MixTracksPCM([][]byte{mainDrum, hihats})
			drumTrack = append(drumTrack, mixedBeat...)
		}
	}

	s.bgmGroove = MixTracksPCM([][]byte{bassTrack, leadTrack, drumTrack})

	// =========================================================================
	// TEMA 2: "Panic Siren Theme" (152 BPM - Aceleração Máxima e Tensão)
	// =========================================================================
	pBeat := 60.0 / 152.0
	var panicBass []byte
	panicNotes := []float64{110, 116.5, 110, 123.47, 110, 130.81, 110, 146.83}
	for r := 0; r < 4; r++ {
		for _, pf := range panicNotes {
			pn := GenerateNotePCM(NoteDef{
				WaveType: "sawtooth", Duration: pBeat / 2.0, StartFreq: pf, EndFreq: pf,
				Volume: 0.26, Attack: 0.005, Decay: 0.05, Sustain: 0.5, Release: 0.05,
			})
			panicBass = append(panicBass, pn...)
		}
	}

	// Sirene de urgência sobreposta
	var panicSirenTrack []byte
	sirenDur := pBeat * 4.0
	for k := 0; k < 4; k++ {
		sNote := GenerateNotePCM(NoteDef{
			WaveType: "square", Duration: sirenDur, StartFreq: 880, EndFreq: 1320,
			DutyCycle: 0.4, VibratoFreq: 9.0, VibratoDepth: 120.0, Volume: 0.18,
			Attack: 0.05, Decay: 0.2, Sustain: 0.6, Release: 0.2, Pan: 0.0,
		})
		panicSirenTrack = append(panicSirenTrack, sNote...)
	}

	s.bgmPanic = MixTracksPCM([][]byte{panicBass, panicSirenTrack})

	// =========================================================================
	// TEMA 3: "Victory Fanfare"
	// =========================================================================
	vNotes := []struct {
		f float64
		d float64
	}{
		{523.25, 0.15}, {659.25, 0.15}, {783.99, 0.15}, {1046.50, 0.45},
		{880.00, 0.18}, {1046.50, 0.65},
	}
	for _, vn := range vNotes {
		part := GenerateNotePCM(NoteDef{
			WaveType: "triangle", Duration: vn.d, StartFreq: vn.f, EndFreq: vn.f,
			Volume: 0.35, Attack: 0.01, Decay: 0.08, Sustain: 0.6, Release: 0.15,
		})
		s.bgmWin = append(s.bgmWin, part...)
	}
}

// PlayBGM toca a música de fundo com loop infinito.
func (s *AudioSystem) PlayBGM(trackName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.currTrack == trackName && s.bgmPlayer != nil && s.bgmPlayer.IsPlaying() {
		return
	}

	if s.bgmPlayer != nil {
		s.bgmPlayer.Close()
		s.bgmPlayer = nil
	}

	var data []byte
	switch trackName {
	case "panic":
		data = s.bgmPanic
	case "victory":
		data = s.bgmWin
	case "groove":
		fallthrough
	default:
		data = s.bgmGroove
	}

	if len(data) == 0 {
		return
	}

	s.currTrack = trackName
	loop := audio.NewInfiniteLoop(bytes.NewReader(data), int64(len(data)))
	p, err := s.ctx.NewPlayer(loop)
	if err == nil {
		p.SetVolume(s.volume * 0.75)
		p.Play()
		s.bgmPlayer = p
	}
}

// StopBGM pausa e fecha o tocador da música de fundo.
func (s *AudioSystem) StopBGM() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.bgmPlayer != nil {
		s.bgmPlayer.Close()
		s.bgmPlayer = nil
		s.currTrack = ""
	}
}

// PlaySFX reproduz um efeito sonoro instantaneamente.
func (s *AudioSystem) PlaySFX(sfxType string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var data []byte
	vol := s.volume

	switch sfxType {
	case "throw":
		data = s.sfxThrow
	case "deliver":
		data = s.sfxDeliver
		vol *= 1.2
	case "combo":
		data = s.sfxCombo
		vol *= 1.3
	case "crash":
		data = s.sfxCrash
	case "bunnyhop":
		data = s.sfxBunnyHop
	case "bell":
		data = s.sfxBell
	case "horn":
		data = s.sfxHorn
	case "bark":
		data = s.sfxBark
	case "panic":
		data = s.sfxPanicSiren
	case "recover":
		data = s.sfxRecoverAction
	case "tornado":
		data = s.sfxTornado
	case "click":
		data = s.sfxClick
	default:
		return
	}

	if len(data) == 0 {
		return
	}

	p := s.ctx.NewPlayerFromBytes(data)
	if p != nil {
		p.SetVolume(vol)
		p.Play()
	}
}
