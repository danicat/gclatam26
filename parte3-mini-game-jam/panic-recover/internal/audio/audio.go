package audio

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/rand"
	"sync"

	ebitenaudio "github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	SampleRate = 44100
)

var (
	ctxOnce sync.Once
	audioCtx *ebitenaudio.Context

	sfxLaser        []byte
	sfxEnemyLaser   []byte
	sfxExplosion    []byte
	sfxPanicSiren   []byte
	sfxRecoverChime []byte
	sfxPickup       []byte
	bgmData         []byte

	bgmPlayer *ebitenaudio.Player
	bgmLoop   *infiniteLoopReader
	bgmMtx    sync.Mutex
)

// Init initializes the audio system and pre-synthesizes all audio buffers.
func Init() {
	ctxOnce.Do(func() {
		audioCtx = ebitenaudio.NewContext(SampleRate)
		sfxLaser = synthLaser()
		sfxEnemyLaser = synthEnemyLaser()
		sfxExplosion = synthExplosion()
		sfxPanicSiren = synthPanicSiren()
		sfxRecoverChime = synthRecoverChime()
		sfxPickup = synthPickup()
		bgmData = synthBGM()

		if audioCtx != nil {
			bgmLoop = &infiniteLoopReader{data: bgmData}
			var err error
			bgmPlayer, err = audioCtx.NewPlayer(bgmLoop)
			if err == nil {
				bgmPlayer.SetVolume(0.4)
				bgmPlayer.Play()
			}
		}
	})
}

// PlayLaser plays the player's blaster sound.
func PlayLaser() {
	playSFX(sfxLaser, 0.25)
}

// PlayEnemyLaser plays enemy firing sound.
func PlayEnemyLaser() {
	playSFX(sfxEnemyLaser, 0.2)
}

// PlayExplosion plays an explosion noise burst.
func PlayExplosion() {
	playSFX(sfxExplosion, 0.45)
}

// PlayPanicSiren plays the panic warning alarm.
func PlayPanicSiren() {
	playSFX(sfxPanicSiren, 0.5)
}

// PlayRecoverChime plays the triumphant recover chime.
func PlayRecoverChime() {
	playSFX(sfxRecoverChime, 0.5)
}

// PlayPickup plays the item pickup sound.
func PlayPickup() {
	playSFX(sfxPickup, 0.35)
}

func playSFX(data []byte, volume float64) {
	if audioCtx == nil || len(data) == 0 {
		return
	}
	p := audioCtx.NewPlayerFromBytes(data)
	if p != nil {
		p.SetVolume(volume)
		p.Play()
	}
}

// infiniteLoopReader streams a byte slice repeatedly without end.
type infiniteLoopReader struct {
	data   []byte
	offset int
}

func (r *infiniteLoopReader) Read(p []byte) (n int, err error) {
	if len(r.data) == 0 {
		return 0, nil
	}
	for n < len(p) {
		remainingData := len(r.data) - r.offset
		needed := len(p) - n
		toCopy := remainingData
		if toCopy > needed {
			toCopy = needed
		}
		copy(p[n:n+toCopy], r.data[r.offset:r.offset+toCopy])
		n += toCopy
		r.offset += toCopy
		if r.offset >= len(r.data) {
			r.offset = 0
		}
	}
	return n, nil
}

// ============================================================================
// DSP SYNTHESIS FUNCTIONS
// ============================================================================

func synthLaser() []byte {
	duration := 0.12
	samples := int(float64(SampleRate) * duration)
	buf := new(bytes.Buffer)
	phase := 0.0

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(SampleRate)
		prog := t / duration
		freq := 950.0 * math.Pow(0.15, prog)
		phase += freq / float64(SampleRate)
		for phase >= 1.0 {
			phase -= 1.0
		}
		var val float64
		if phase < 0.5 {
			val = 0.8
		} else {
			val = -0.8
		}
		env := 1.0 - prog
		sample := int16(val * env * 28000)
		_ = binary.Write(buf, binary.LittleEndian, sample) // Left
		_ = binary.Write(buf, binary.LittleEndian, sample) // Right
	}
	return buf.Bytes()
}

func synthEnemyLaser() []byte {
	duration := 0.14
	samples := int(float64(SampleRate) * duration)
	buf := new(bytes.Buffer)
	phase := 0.0

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(SampleRate)
		prog := t / duration
		freq := 420.0 * math.Pow(0.3, prog)
		phase += freq / float64(SampleRate)
		for phase >= 1.0 {
			phase -= 1.0
		}
		val := math.Sin(2.0 * math.Pi * phase)
		env := math.Sin(prog * math.Pi)
		sample := int16(val * env * 26000)
		_ = binary.Write(buf, binary.LittleEndian, sample)
		_ = binary.Write(buf, binary.LittleEndian, sample)
	}
	return buf.Bytes()
}

func synthExplosion() []byte {
	duration := 0.38
	samples := int(float64(SampleRate) * duration)
	buf := new(bytes.Buffer)
	filter := 0.0

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(SampleRate)
		prog := t / duration
		raw := rand.Float64()*2.0 - 1.0
		filterCoeff := 0.25 * (1.0 - prog*0.7)
		filter += filterCoeff * (raw - filter)
		env := math.Exp(-6.0 * prog)
		sample := int16(filter * env * 31000)
		_ = binary.Write(buf, binary.LittleEndian, sample)
		_ = binary.Write(buf, binary.LittleEndian, sample)
	}
	return buf.Bytes()
}

func synthPanicSiren() []byte {
	duration := 0.55
	samples := int(float64(SampleRate) * duration)
	buf := new(bytes.Buffer)
	phase := 0.0

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(SampleRate)
		prog := t / duration
		// Fast oscillating warble between 550Hz and 880Hz
		freq := 700.0 + 180.0*math.Sin(2.0*math.Pi*12.0*t)
		phase += freq / float64(SampleRate)
		for phase >= 1.0 {
			phase -= 1.0
		}
		var val float64
		if phase < 0.4 {
			val = 0.9
		} else {
			val = -0.9
		}
		env := 1.0
		if prog > 0.85 {
			env = (1.0 - prog) / 0.15
		}
		sample := int16(val * env * 24000)
		_ = binary.Write(buf, binary.LittleEndian, sample)
		_ = binary.Write(buf, binary.LittleEndian, sample)
	}
	return buf.Bytes()
}

func synthRecoverChime() []byte {
	// Uplifting arpeggio: C5 (523.25), E5 (659.25), G5 (783.99), C6 (1046.5)
	notes := []float64{523.25, 659.25, 783.99, 1046.50}
	noteDur := 0.12
	totalDur := noteDur * float64(len(notes)) + 0.25
	samples := int(float64(SampleRate) * totalDur)
	buf := new(bytes.Buffer)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(SampleRate)
		val := 0.0
		for idx, freq := range notes {
			noteStart := float64(idx) * noteDur
			if t >= noteStart {
				localT := t - noteStart
				decay := math.Exp(-4.0 * localT)
				// Harmonic bell tone
				tone := math.Sin(2.0*math.Pi*freq*localT) + 0.4*math.Sin(4.0*math.Pi*freq*localT)
				val += tone * decay
			}
		}
		val = math.Max(-1.0, math.Min(1.0, val*0.45))
		sample := int16(val * 30000)
		_ = binary.Write(buf, binary.LittleEndian, sample)
		_ = binary.Write(buf, binary.LittleEndian, sample)
	}
	return buf.Bytes()
}

func synthPickup() []byte {
	notes := []float64{659.25, 987.77} // E5, B5
	noteDur := 0.07
	totalDur := noteDur * 2
	samples := int(float64(SampleRate) * totalDur)
	buf := new(bytes.Buffer)
	phase := 0.0

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(SampleRate)
		freq := notes[0]
		localT := t
		if t >= noteDur {
			freq = notes[1]
			localT = t - noteDur
		}
		phase += freq / float64(SampleRate)
		for phase >= 1.0 {
			phase -= 1.0
		}
		val := math.Sin(2.0 * math.Pi * phase)
		env := 1.0 - (localT / noteDur)
		sample := int16(val * env * 25000)
		_ = binary.Write(buf, binary.LittleEndian, sample)
		_ = binary.Write(buf, binary.LittleEndian, sample)
	}
	return buf.Bytes()
}

func synthBGM() []byte {
	// A driving 4-bar cyber chiptune loop at 135 BPM
	bpm := 135.0
	secondsPerBeat := 60.0 / bpm
	// 4 bars of 4/4 = 16 beats
	totalBeats := 16.0
	totalDuration := totalBeats * secondsPerBeat
	totalSamples := int(float64(SampleRate) * totalDuration)

	// Bassline sequence in A minor (A1, C2, D2, F2)
	// Frequencies: A1=55, C2=65.4, D2=73.4, E2=82.4, F2=87.3, G2=98
	bassNotes := []float64{
		55.0, 55.0, 110.0, 55.0,
		65.4, 65.4, 130.8, 65.4,
		73.4, 73.4, 146.8, 73.4,
		87.3, 87.3, 98.0,  82.4,
	}

	// Arpeggio melody notes
	// A minor scale chords: A4=440, C5=523.25, E5=659.25, G5=783.99
	leadNotes := []float64{
		440.0, 523.25, 659.25, 523.25,
		440.0, 659.25, 783.99, 659.25,
		523.25, 659.25, 880.0,  659.25,
		783.99, 659.25, 523.25, 440.0,
	}

	buf := new(bytes.Buffer)
	bassPhase := 0.0
	leadPhase := 0.0
	noiseFilter := 0.0

	for i := 0; i < totalSamples; i++ {
		t := float64(i) / float64(SampleRate)
		beatProgress := t / secondsPerBeat
		stepIndex := int(beatProgress) % 16
		stepSubTime := math.Mod(t, secondsPerBeat)

		// 1. Bassline (Sawtooth with punchy decay)
		bassFreq := bassNotes[stepIndex]
		bassPhase += bassFreq / float64(SampleRate)
		for bassPhase >= 1.0 {
			bassPhase -= 1.0
		}
		bassVal := (2.0*bassPhase - 1.0) * math.Exp(-4.0*(stepSubTime/secondsPerBeat))

		// 2. Lead Arpeggio (Pulse wave with slight vibrato)
		leadFreq := leadNotes[stepIndex]
		vibrato := 3.0 * math.Sin(2.0*math.Pi*6.0*t)
		leadPhase += (leadFreq + vibrato) / float64(SampleRate)
		for leadPhase >= 1.0 {
			leadPhase -= 1.0
		}
		var leadVal float64
		if leadPhase < 0.35 {
			leadVal = 0.5
		} else {
			leadVal = -0.5
		}
		leadEnv := math.Exp(-2.5 * (stepSubTime / secondsPerBeat))

		// 3. Drum Percussion (Kick on beat 0, 1, 2, 3; Snare on beat 1, 3)
		subBeatTime := math.Mod(t, secondsPerBeat)
		subBeatProg := subBeatTime / secondsPerBeat

		// Kick: pitch drop sine
		kickVal := 0.0
		if subBeatProg < 0.25 {
			kickFreq := 140.0 * math.Exp(-16.0*subBeatProg)
			kickVal = math.Sin(2.0*math.Pi*kickFreq*subBeatTime) * math.Exp(-8.0*subBeatProg)
		}

		// Snare: noise burst on beats 1 and 3 (0-indexed: 1, 3, 5, 7, 9, 11, 13, 15)
		snareVal := 0.0
		if stepIndex%2 == 1 && subBeatProg < 0.2 {
			rawNoise := rand.Float64()*2.0 - 1.0
			noiseFilter += 0.4 * (rawNoise - noiseFilter)
			snareVal = noiseFilter * (1.0 - subBeatProg/0.2)
		}

		// Hi-hat: every 16th sub-beat
		hihatTime := math.Mod(t, secondsPerBeat/2.0)
		hihatVal := 0.0
		if hihatTime < 0.04 {
			rawNoise := rand.Float64()*2.0 - 1.0
			hihatVal = rawNoise * (1.0 - hihatTime/0.04) * 0.3
		}

		// Master Mix
		mix := 0.32*bassVal + 0.20*(leadVal*leadEnv) + 0.30*kickVal + 0.22*snareVal + 0.12*hihatVal
		mix = math.Max(-1.0, math.Min(1.0, mix))
		sample := int16(mix * 22000)

		_ = binary.Write(buf, binary.LittleEndian, sample) // Left
		_ = binary.Write(buf, binary.LittleEndian, sample) // Right
	}
	return buf.Bytes()
}
