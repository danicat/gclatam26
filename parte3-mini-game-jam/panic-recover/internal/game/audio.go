package game

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	sampleRate = 44100
)

var (
	audioCtx *audio.Context

	sndStep      []byte
	sndPush      []byte
	sndFillHole  []byte
	sndRecover   []byte
	sndHeartbeat []byte
	sndWin       []byte
	sndFaint     []byte
	sndTick      []byte

	bgmPlayer *audio.Player
)

func initAudio() {
	if audioCtx == nil {
		audioCtx = audio.NewContext(sampleRate)
	}

	sndStep = synthStep()
	sndPush = synthPush()
	sndFillHole = synthFillHole()
	sndRecover = synthRecover()
	sndHeartbeat = synthHeartbeat()
	sndWin = synthWin()
	sndFaint = synthFaint()
	sndTick = synthClockTick()

	// Initialize background ominous dungeon soundtrack
	bgmData := synthOminousDungeonTheme()
	if len(bgmData) > 0 {
		loop := audio.NewInfiniteLoop(bytes.NewReader(bgmData), int64(len(bgmData)))
		var err error
		bgmPlayer, err = audioCtx.NewPlayer(loop)
		if err == nil {
			bgmPlayer.SetVolume(0.42)
			bgmPlayer.Play()
		}
	}
}

func playSound(data []byte, vol float64) {
	if audioCtx == nil || len(data) == 0 {
		return
	}
	p := audioCtx.NewPlayerFromBytes(data)
	p.SetVolume(vol)
	p.Play()
}

func pcmBytes(samples []int16) []byte {
	buf := new(bytes.Buffer)
	for _, s := range samples {
		// 16-bit stereo (L, R)
		binary.Write(buf, binary.LittleEndian, s)
		binary.Write(buf, binary.LittleEndian, s)
	}
	return buf.Bytes()
}

// synthOminousDungeonTheme generates a 16-bar ominous, suspenseful Lovecraftian soundtrack
// featuring deep detuned sub-bass, dissonant diminished pads, clock ticking, heartbeat thumps,
// and spectral chimes.
func synthOminousDungeonTheme() []byte {
	bpm := 60.0
	beatDur := 60.0 / bpm   // 1.0s per beat
	barDur := beatDur * 4.0 // 4.0s per bar
	totalBars := 16
	totalDuration := barDur * float64(totalBars) // 64 seconds rich loop
	numSamples := int(sampleRate * totalDuration)
	samples := make([]int16, numSamples)

	// Ominous chord progression in D minor with tritones & diminished suspense
	chords := []struct {
		bassFreq float64
		p1, p2   float64
	}{
		{36.71, 146.83, 220.00}, // Bar 1-2: Dm (D1, D3, A3)
		{36.71, 146.83, 220.00},
		{38.89, 155.56, 207.65}, // Bar 3-4: D dim / Eb tension (Eb1, Eb3, Ab3)
		{38.89, 155.56, 207.65},
		{36.71, 146.83, 233.08}, // Bar 5-6: Bbm/D (D1, D3, Bb3)
		{36.71, 146.83, 233.08},
		{55.00, 138.59, 233.08}, // Bar 7-8: A7b9 (A1, C#3, Bb3)
		{55.00, 138.59, 233.08},
		{36.71, 146.83, 220.00}, // Bar 9-10: Dm return
		{36.71, 146.83, 220.00},
		{38.89, 155.56, 233.08}, // Bar 11-12: Ebm dread (Eb1, Eb3, Bb3)
		{38.89, 155.56, 233.08},
		{51.91, 146.83, 246.94}, // Bar 13-14: G# dim tritone (G#1, D3, B3)
		{51.91, 146.83, 246.94},
		{55.00, 138.59, 220.00}, // Bar 15-16: A dim cadence to D
		{36.71, 146.83, 220.00},
	}

	// Spectral high chime melodies (D minor / Locrian)
	chimeNotes := []float64{
		587.33, 0.0, 554.37, 0.0, 622.25, 0.0, 587.33, 0.0, // D5, C#5, Eb5, D5
		0.0, 698.46, 0.0, 659.25, 587.33, 0.0, 554.37, 0.0, // F5, E5, D5, C#5
	}

	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		barIdx := int(t / barDur)
		if barIdx >= len(chords) {
			barIdx = len(chords) - 1
		}
		chord := chords[barIdx]

		// 1. Channel 1: Detuned Sub-Bass Drone (36.7Hz + 37.1Hz slow acoustic beating)
		drone1 := math.Sin(2 * math.Pi * chord.bassFreq * t)
		drone2 := math.Sin(2 * math.Pi * (chord.bassFreq + 0.4) * t)
		subBass := (drone1 + drone2) * 0.5

		// 2. Channel 2: Dissonant Diminished Pad with subtle slow LFO vibrato
		lfo := math.Sin(2*math.Pi*0.25*t) * 1.5
		p1 := math.Sin(2*math.Pi*(chord.p1+lfo)*t) * 0.35
		p2 := math.Sin(2*math.Pi*(chord.p2-lfo)*t) * 0.30
		barTime := math.Mod(t, barDur)
		padEnv := math.Sin((barTime / barDur) * math.Pi)
		pad := (p1 + p2) * padEnv

		// 3. Channel 3: Clock Ticking Track ("tick... tock... tick... tock..." every 0.5s)
		tickSubT := math.Mod(t, 0.5)
		var tick float64 = 0.0
		if tickSubT < 0.03 {
			// Higher pitch on beat, lower on half-beat
			beatNum := int(t / 0.5)
			tickFreq := 1600.0
			if beatNum%2 == 1 {
				tickFreq = 1100.0 // Tock
			}
			tEnv := math.Exp(-tickSubT * 140.0)
			noise := (rand.Float64()*2 - 1) * 0.2
			tick = (math.Sin(2*math.Pi*tickFreq*tickSubT) + noise) * tEnv * 0.22
		}

		// 4. Channel 4: Deep Ominous Heartbeat / Thud (on beat 1 and 3 of every measure)
		measureTime := math.Mod(t, 2.0)
		var thud float64 = 0.0
		if measureTime < 0.18 {
			thudFreq := 58.0 * math.Exp(-measureTime*12.0)
			thudEnv := math.Sin((measureTime / 0.18) * math.Pi)
			thud = math.Sin(2*math.Pi*thudFreq*measureTime) * thudEnv * 0.45
		}

		// 5. Channel 5: Spectral high chime bells (sparse, ghostly)
		chimeStep := int(t / 4.0)
		chimeFreq := chimeNotes[chimeStep%len(chimeNotes)]
		var chime float64 = 0.0
		if chimeFreq > 0.0 {
			chimeSubT := math.Mod(t, 4.0)
			if chimeSubT < 1.8 {
				cEnv := math.Exp(-chimeSubT * 2.5)
				chime = math.Sin(2*math.Pi*chimeFreq*t) * cEnv * 0.25
			}
		}

		// 6. Channel 6: Cosmic wind noise breathing swell
		windEnv := (math.Sin(2*math.Pi*0.125*t) + 1.0) * 0.5
		wind := (rand.Float64()*2 - 1) * 0.06 * windEnv

		// 32-bit mixing with hard-clamping to prevent clipping
		mix := (subBass*0.35 + pad*0.35 + tick + thud + chime + wind) * 17500.0
		if mix > 32767 {
			mix = 32767
		} else if mix < -32768 {
			mix = -32768
		}

		samples[i] = int16(mix)
	}

	return pcmBytes(samples)
}

// Clock tick sound (mechanical pocket-watch click)
func synthClockTick() []byte {
	duration := 0.03
	numSamples := int(sampleRate * duration)
	samples := make([]int16, numSamples)
	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		freq := 1800.0 - (t / duration * 800.0)
		env := math.Exp(-t * 160.0)
		noise := (rand.Float64()*2 - 1) * 0.3
		val := (math.Sin(2*math.Pi*freq*t)*0.7 + noise) * env
		samples[i] = int16(val * 14000)
	}
	return pcmBytes(samples)
}

// Footstep sound
func synthStep() []byte {
	duration := 0.04
	numSamples := int(sampleRate * duration)
	samples := make([]int16, numSamples)
	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		freq := 480.0 - (t / duration * 200.0)
		env := 1.0 - (float64(i) / float64(numSamples))
		val := math.Sin(2*math.Pi*freq*t) * env
		samples[i] = int16(val * 8000)
	}
	return pcmBytes(samples)
}

// Boulder push rumble
func synthPush() []byte {
	duration := 0.12
	numSamples := int(sampleRate * duration)
	samples := make([]int16, numSamples)
	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		freq := 140.0 - (t / duration * 60.0)
		noise := (rand.Float64()*2.0 - 1.0) * 0.4
		env := math.Sin((float64(i) / float64(numSamples)) * math.Pi)
		val := (math.Sin(2*math.Pi*freq*t)*0.6 + noise) * env
		samples[i] = int16(val * 12000)
	}
	return pcmBytes(samples)
}

// Boulder falling into hole
func synthFillHole() []byte {
	duration := 0.22
	numSamples := int(sampleRate * duration)
	samples := make([]int16, numSamples)
	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		freq := 220.0*math.Exp(-t*15.0) + 70.0
		noise := (rand.Float64()*2.0 - 1.0) * 0.3
		env := (1.0 - float64(i)/float64(numSamples))
		val := (math.Sin(2*math.Pi*freq*t)*0.7 + noise) * env
		samples[i] = int16(val * 16000)
	}
	return pcmBytes(samples)
}

// Recovery chime (Eldritch sanity soothing)
func synthRecover() []byte {
	duration := 0.35
	numSamples := int(sampleRate * duration)
	samples := make([]int16, numSamples)
	notes := []float64{523.25, 659.25, 783.99, 1046.50}
	noteDur := duration / float64(len(notes))

	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		noteIdx := int(t / noteDur)
		if noteIdx >= len(notes) {
			noteIdx = len(notes) - 1
		}
		freq := notes[noteIdx]
		subT := math.Mod(t, noteDur)
		env := math.Exp(-subT * 12.0)
		val := math.Sin(2*math.Pi*freq*t) * env
		samples[i] = int16(val * 14000)
	}
	return pcmBytes(samples)
}

// Panic Heartbeat pulse
func synthHeartbeat() []byte {
	duration := 0.12
	numSamples := int(sampleRate * duration)
	samples := make([]int16, numSamples)
	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		freq := 65.0 - (t / duration * 25.0)
		env := math.Sin((float64(i) / float64(numSamples)) * math.Pi)
		val := math.Sin(2*math.Pi*freq*t) * env
		samples[i] = int16(val * 20000)
	}
	return pcmBytes(samples)
}

// Room victory fanfare
func synthWin() []byte {
	duration := 0.45
	numSamples := int(sampleRate * duration)
	samples := make([]int16, numSamples)
	notes := []float64{440.0, 554.37, 659.25, 880.0}
	noteDur := duration / float64(len(notes))

	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		noteIdx := int(t / noteDur)
		if noteIdx >= len(notes) {
			noteIdx = len(notes) - 1
		}
		freq := notes[noteIdx]
		env := 1.0 - (float64(i) / float64(numSamples))
		val := (math.Sin(2*math.Pi*freq*t)*0.7 + math.Sin(4*math.Pi*freq*t)*0.3) * env
		samples[i] = int16(val * 15000)
	}
	return pcmBytes(samples)
}

// Panic fainting sound
func synthFaint() []byte {
	duration := 0.4
	numSamples := int(sampleRate * duration)
	samples := make([]int16, numSamples)
	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		freq := 350.0 * (1.0 - t/duration)
		env := 1.0 - (float64(i) / float64(numSamples))
		val := (math.Sin(2*math.Pi*freq*t) + (rand.Float64()*2-1)*0.2) * env
		samples[i] = int16(val * 15000)
	}
	return pcmBytes(samples)
}
