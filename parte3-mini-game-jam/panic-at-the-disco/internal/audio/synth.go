package audio

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	SampleRate = 44100
)

type AudioEngine struct {
	ctx           *audio.Context
	bgmPlayer     *audio.Player
	bgmFastPlayer *audio.Player
	isFastBGM     bool

	sfxWhistleBuf []byte
	sfxCrashBuf   []byte
	sfxDashBuf    []byte
	sfxSlipBuf    []byte
	sfxWinBuf     []byte
	sfxAlarmBuf   []byte
}

var globalEngine *AudioEngine

// InitAudioEngine initializes the audio context and procedural sound buffers.
func InitAudioEngine() *AudioEngine {
	if globalEngine != nil {
		return globalEngine
	}

	ctx := audio.NewContext(SampleRate)
	ae := &AudioEngine{
		ctx: ctx,
	}

	// 1. Synthesize SFX buffers
	ae.sfxWhistleBuf = synthesizeWhistle()
	ae.sfxCrashBuf = synthesizeCrash()
	ae.sfxDashBuf = synthesizeDash()
	ae.sfxSlipBuf = synthesizeSlip()
	ae.sfxWinBuf = synthesizeWin()
	ae.sfxAlarmBuf = synthesizeAlarm()

	// 2. Synthesize BGM loops (120 BPM normal, 144 BPM fast panic)
	bgmNormalBuf := synthesizeDiscoBGM(120.0, false)
	loopNormal := audio.NewInfiniteLoop(bytes.NewReader(bgmNormalBuf), int64(len(bgmNormalBuf)))
	pNormal, err := ctx.NewPlayer(loopNormal)
	if err == nil {
		pNormal.SetVolume(0.55)
		ae.bgmPlayer = pNormal
	}

	bgmFastBuf := synthesizeDiscoBGM(144.0, true)
	loopFast := audio.NewInfiniteLoop(bytes.NewReader(bgmFastBuf), int64(len(bgmFastBuf)))
	pFast, err := ctx.NewPlayer(loopFast)
	if err == nil {
		pFast.SetVolume(0.60)
		ae.bgmFastPlayer = pFast
	}

	globalEngine = ae
	return ae
}

func GetAudioEngine() *AudioEngine {
	if globalEngine == nil {
		return InitAudioEngine()
	}
	return globalEngine
}

func (ae *AudioEngine) PlayBGM(fastMode bool) {
	if ae == nil || ae.ctx == nil {
		return
	}
	if fastMode {
		if ae.bgmPlayer != nil && ae.bgmPlayer.IsPlaying() {
			ae.bgmPlayer.Pause()
		}
		if ae.bgmFastPlayer != nil && !ae.bgmFastPlayer.IsPlaying() {
			ae.bgmFastPlayer.Play()
		}
		ae.isFastBGM = true
	} else {
		if ae.bgmFastPlayer != nil && ae.bgmFastPlayer.IsPlaying() {
			ae.bgmFastPlayer.Pause()
		}
		if ae.bgmPlayer != nil && !ae.bgmPlayer.IsPlaying() {
			ae.bgmPlayer.Play()
		}
		ae.isFastBGM = false
	}
}

func (ae *AudioEngine) PauseBGM() {
	if ae == nil {
		return
	}
	if ae.bgmPlayer != nil && ae.bgmPlayer.IsPlaying() {
		ae.bgmPlayer.Pause()
	}
	if ae.bgmFastPlayer != nil && ae.bgmFastPlayer.IsPlaying() {
		ae.bgmFastPlayer.Pause()
	}
}

func (ae *AudioEngine) StopBGM() {
	if ae == nil {
		return
	}
	ae.PauseBGM()
	if ae.bgmPlayer != nil {
		_ = ae.bgmPlayer.Rewind()
	}
	if ae.bgmFastPlayer != nil {
		_ = ae.bgmFastPlayer.Rewind()
	}
}

func (ae *AudioEngine) PlaySFXWhistle() {
	if ae == nil || ae.ctx == nil {
		return
	}
	p := ae.ctx.NewPlayerFromBytes(ae.sfxWhistleBuf)
	p.SetVolume(0.7)
	p.Play()
}

func (ae *AudioEngine) PlaySFXCrash() {
	if ae == nil || ae.ctx == nil {
		return
	}
	p := ae.ctx.NewPlayerFromBytes(ae.sfxCrashBuf)
	p.SetVolume(0.85)
	p.Play()
}

func (ae *AudioEngine) PlaySFXDash() {
	if ae == nil || ae.ctx == nil {
		return
	}
	p := ae.ctx.NewPlayerFromBytes(ae.sfxDashBuf)
	p.SetVolume(0.75)
	p.Play()
}

func (ae *AudioEngine) PlaySFXSlip() {
	if ae == nil || ae.ctx == nil {
		return
	}
	p := ae.ctx.NewPlayerFromBytes(ae.sfxSlipBuf)
	p.SetVolume(0.7)
	p.Play()
}

func (ae *AudioEngine) PlaySFXWin() {
	if ae == nil || ae.ctx == nil {
		return
	}
	p := ae.ctx.NewPlayerFromBytes(ae.sfxWinBuf)
	p.SetVolume(0.8)
	p.Play()
}

func (ae *AudioEngine) PlaySFXAlarm() {
	if ae == nil || ae.ctx == nil {
		return
	}
	p := ae.ctx.NewPlayerFromBytes(ae.sfxAlarmBuf)
	p.SetVolume(0.65)
	p.Play()
}

// -------------------------------------------------------------
// Procedural Waveform Synthesis Helpers (16-bit Stereo PCM)
// -------------------------------------------------------------

func encodePCM16(left, right []float64) []byte {
	buf := make([]byte, len(left)*4)
	for i := range left {
		l := math.Max(-1.0, math.Min(1.0, left[i]))
		r := math.Max(-1.0, math.Min(1.0, right[i]))

		sampleL := int16(l * 32767.0)
		sampleR := int16(r * 32767.0)

		binary.LittleEndian.PutUint16(buf[i*4:], uint16(sampleL))
		binary.LittleEndian.PutUint16(buf[i*4+2:], uint16(sampleR))
	}
	return buf
}

// synthesizeWhistle generates a downward whistling FM chirp of a falling object.
func synthesizeWhistle() []byte {
	duration := 0.45
	samples := int(float64(SampleRate) * duration)
	left := make([]float64, samples)
	right := make([]float64, samples)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(SampleRate)
		progress := t / duration
		// Pitch drops from 1100 Hz down to 220 Hz
		freq := 1100.0 - 880.0*math.Pow(progress, 1.4)
		phase := 2.0 * math.Pi * freq * t
		val := math.Sin(phase) * (1.0 - progress*0.7)
		left[i] = val * 0.4
		right[i] = val * 0.4
	}
	return encodePCM16(left, right)
}

// synthesizeCrash generates a massive low-end thump + shattering noise burst.
func synthesizeCrash() []byte {
	duration := 0.6
	samples := int(float64(SampleRate) * duration)
	left := make([]float64, samples)
	right := make([]float64, samples)

	rnd := rand.New(rand.NewSource(42))
	for i := 0; i < samples; i++ {
		t := float64(i) / float64(SampleRate)
		env := math.Exp(-t * 8.0) // Fast exponential decay

		// Low-end bass thump (dropping from 120Hz to 35Hz)
		bassFreq := 120.0 - 85.0*(t/duration)
		thump := math.Sin(2.0*math.Pi*bassFreq*t) * env * 0.7

		// High-pitch glass shatter noise
		noise := (rnd.Float64()*2.0 - 1.0) * math.Exp(-t*14.0) * 0.5
		shatterTink := math.Sin(2.0*math.Pi*2800.0*t) * math.Exp(-t*20.0) * 0.3

		sig := thump + noise + shatterTink
		left[i] = sig
		right[i] = sig
	}
	return encodePCM16(left, right)
}

// synthesizeDash generates a funky glissando sweep with sparkle.
func synthesizeDash() []byte {
	duration := 0.25
	samples := int(float64(SampleRate) * duration)
	left := make([]float64, samples)
	right := make([]float64, samples)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(SampleRate)
		prog := t / duration
		freq := 300.0 + 900.0*math.Pow(prog, 2.0)
		env := math.Sin(prog * math.Pi)
		sig := (math.Sin(2.0*math.Pi*freq*t) + 0.3*math.Sin(2.0*math.Pi*freq*2.0*t)) * env * 0.5
		left[i] = sig * (1.0 - prog*0.3)
		right[i] = sig * (0.7 + prog*0.3) // Stereo pan sweep
	}
	return encodePCM16(left, right)
}

// synthesizeSlip generates a comic cartoon slide whistle.
func synthesizeSlip() []byte {
	duration := 0.3
	samples := int(float64(SampleRate) * duration)
	left := make([]float64, samples)
	right := make([]float64, samples)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(SampleRate)
		prog := t / duration
		freq := 500.0 + 350.0*math.Sin(prog*math.Pi*2.0)
		env := 1.0 - prog
		sig := math.Sin(2.0*math.Pi*freq*t) * env * 0.4
		left[i] = sig
		right[i] = sig
	}
	return encodePCM16(left, right)
}

// synthesizeWin generates an upbeat ascending major triad fanfare.
func synthesizeWin() []byte {
	notes := []float64{523.25, 659.25, 783.99, 1046.50} // C5, E5, G5, C6
	noteDuration := 0.12
	totalDuration := noteDuration * float64(len(notes))
	samples := int(float64(SampleRate) * totalDuration)
	left := make([]float64, samples)
	right := make([]float64, samples)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(SampleRate)
		noteIdx := int(t / noteDuration)
		if noteIdx >= len(notes) {
			noteIdx = len(notes) - 1
		}
		freq := notes[noteIdx]
		subT := math.Mod(t, noteDuration)
		env := math.Exp(-subT * 7.0)
		sig := (math.Sin(2.0*math.Pi*freq*t) + 0.4*math.Sin(2.0*math.Pi*freq*2.0*t)) * env * 0.45
		left[i] = sig
		right[i] = sig
	}
	return encodePCM16(left, right)
}

// synthesizeAlarm generates a pulsing two-tone danger siren.
func synthesizeAlarm() []byte {
	duration := 0.4
	samples := int(float64(SampleRate) * duration)
	left := make([]float64, samples)
	right := make([]float64, samples)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(SampleRate)
		freq := 700.0
		if math.Mod(t, 0.2) > 0.1 {
			freq = 950.0
		}
		sig := math.Sin(2.0*math.Pi*freq*t) * 0.35
		left[i] = sig
		right[i] = sig
	}
	return encodePCM16(left, right)
}

// synthesizeDiscoBGM generates a seamless looping 70s Disco Funk soundtrack in Go memory.
func synthesizeDiscoBGM(bpm float64, includeSiren bool) []byte {
	// 4 bars of 4/4 = 16 beats
	secondsPerBeat := 60.0 / bpm
	totalDuration := secondsPerBeat * 16.0
	samples := int(float64(SampleRate) * totalDuration)
	left := make([]float64, samples)
	right := make([]float64, samples)

	rnd := rand.New(rand.NewSource(777))

	// Bassline notes for 4 bars (A minor disco groove: A, C, D, E, G)
	bassFrequencies := []float64{
		110.0, 110.0, 130.81, 146.83, // Bar 1: A2, A2, C3, D3
		110.0, 220.0, 164.81, 146.83, // Bar 2: A2, A3 (octave jump), E3, D3
		110.0, 110.0, 130.81, 146.83, // Bar 3: A2, A2, C3, D3
		196.0, 164.81, 146.83, 130.81, // Bar 4: G3, E3, D3, C3
	}

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(SampleRate)
		beatIndex := int(t / secondsPerBeat)
		beatFrac := math.Mod(t, secondsPerBeat) / secondsPerBeat
		subBeatT := math.Mod(t, secondsPerBeat)

		// 1. Four-on-the-floor kick drum (punchy 55Hz pitch drop on every beat)
		kickEnv := math.Exp(-subBeatT * 22.0)
		kickFreq := 130.0 - 85.0*beatFrac
		kick := math.Sin(2.0*math.Pi*kickFreq*t) * kickEnv * 0.55

		// 2. Disco Hi-Hat: 16th-note sizzle with accented open hi-hat on the "&" of each beat (subBeatT around 0.5)
		sixteenthT := math.Mod(t, secondsPerBeat/4.0)
		isAndOfBeat := beatFrac >= 0.45 && beatFrac <= 0.65
		hatDecay := 35.0
		if isAndOfBeat {
			hatDecay = 14.0 // Open hi-hat sustain
		}
		hatEnv := math.Exp(-sixteenthT * hatDecay)
		noise := (rnd.Float64()*2.0 - 1.0)
		hihat := noise * hatEnv * 0.22

		// 3. Funky 70s Slap Bass (Octave Sawtooth)
		bassNote := bassFrequencies[beatIndex%len(bassFrequencies)]
		// 8th note bass articulation
		eighthT := math.Mod(t, secondsPerBeat/2.0)
		bassEnv := math.Exp(-eighthT * 6.5)
		// Sawtooth synthesis with 3 harmonics
		bassPhase := math.Mod(t*bassNote, 1.0)
		saw := (2.0*bassPhase - 1.0) + 0.5*(2.0*math.Mod(t*bassNote*2.0, 1.0)-1.0)
		bass := saw * bassEnv * 0.30

		// 4. Synth Brass Stabs on beats 2 and 4
		synthStab := 0.0
		if beatIndex%2 == 1 && beatFrac < 0.3 {
			stabEnv := math.Exp(-beatFrac * 10.0)
			stabFreq := 440.0 // A4 chord
			synthStab = (math.Sin(2.0*math.Pi*stabFreq*t) + 0.5*math.Sin(2.0*math.Pi*stabFreq*1.5*t)) * stabEnv * 0.20
		}

		// 5. Emergency Panic Siren (if enabled)
		siren := 0.0
		if includeSiren {
			sirenFreq := 800.0 + 350.0*math.Sin(t*math.Pi*3.0)
			siren = math.Sin(2.0*math.Pi*sirenFreq*t) * 0.15
		}

		sigL := kick + hihat*0.9 + bass + synthStab*0.8 + siren
		sigR := kick + hihat*1.1 + bass + synthStab*1.2 + siren

		left[i] = sigL * 0.75
		right[i] = sigR * 0.75
	}

	return encodePCM16(left, right)
}
