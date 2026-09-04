package system

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	SampleRate     = 44100
	BytesPerSample = 4 // 16-bit stereo = 2 bytes * 2 channels
)

// AudioManager synthesizes 16-bit PCM sound effects and background music entirely in code.
type AudioManager struct {
	context     *audio.Context
	bgmPlayer   *audio.Player
	sfxCooldown map[string]time.Time
	mu          sync.Mutex

	// Pre-synthesized SFX PCM buffers
	sfxWarning    []byte
	sfxAlarm      []byte
	sfxSteam      []byte
	sfxExtinguish []byte
	sfxClutch     []byte
	sfxExplosion  []byte
	sfxCat        []byte
	sfxPickup     []byte

	bgmNormal []byte
	bgmPanic  []byte
	isPanic   bool
	muted     bool
}

// NewAudioManager initializes the audio context and pre-synthesizes all audio buffers at startup.
func NewAudioManager() *AudioManager {
	ctx := audio.NewContext(SampleRate)
	am := &AudioManager{
		context:     ctx,
		sfxCooldown: make(map[string]time.Time),
	}

	am.synthesizeAll()
	return am
}

func (am *AudioManager) synthesizeAll() {
	// 1. Warning Beep (Square wave beep)
	am.sfxWarning = synthesizeTone(0.08, 880, 880, 0.5, "square", 0.3)

	// 2. Alarm Klaxon (Two-tone alternating siren)
	b1 := synthesizeTone(0.12, 600, 600, 0.5, "sawtooth", 0.35)
	b2 := synthesizeTone(0.12, 900, 900, 0.5, "sawtooth", 0.35)
	am.sfxAlarm = append(b1, b2...)

	// 3. Steam Hiss (Filtered noise burst)
	am.sfxSteam = synthesizeNoise(0.25, 0.25, 0.8)

	// 4. Extinguisher (Sustained gentle foam hiss)
	am.sfxExtinguish = synthesizeNoise(0.2, 0.2, 0.5)

	// 5. Clutch Recovery (Triumphant ascending arpeggio C5 -> E5 -> G5 -> C6)
	c5 := synthesizeTone(0.06, 523.25, 523.25, 0.5, "triangle", 0.4)
	e5 := synthesizeTone(0.06, 659.25, 659.25, 0.5, "triangle", 0.4)
	g5 := synthesizeTone(0.06, 783.99, 783.99, 0.5, "triangle", 0.4)
	c6 := synthesizeTone(0.18, 1046.50, 1046.50, 0.5, "sine", 0.5)
	am.sfxClutch = append(append(append(c5, e5...), g5...), c6...)

	// 6. Explosion (Heavy noise burst with steep pitch drop)
	am.sfxExplosion = synthesizeExplosion(0.6)

	// 7. Cat Meow (Chirpy frequency glide)
	am.sfxCat = synthesizeTone(0.25, 450, 800, 0.5, "sine", 0.35)

	// 8. Tool Pickup (Quick crisp pop)
	am.sfxPickup = synthesizeTone(0.05, 300, 600, 0.5, "square", 0.25)

	// 9. Synthesize Multi-track BGM (Normal & Panic)
	am.bgmNormal = synthesizeBGM(false)
	am.bgmPanic = synthesizeBGM(true)
}

// StartBGM initiates looping background music.
func (am *AudioManager) StartBGM() {
	if am.context == nil || am.muted {
		return
	}
	am.mu.Lock()
	defer am.mu.Unlock()

	if am.bgmPlayer != nil {
		am.bgmPlayer.Close()
	}

	buf := am.bgmNormal
	if am.isPanic {
		buf = am.bgmPanic
	}

	loop := audio.NewInfiniteLoop(bytes.NewReader(buf), int64(len(buf)))
	p, err := am.context.NewPlayer(loop)
	if err == nil {
		am.bgmPlayer = p
		am.bgmPlayer.SetVolume(0.5)
		am.bgmPlayer.Play()
	}
}

// SetPanicMode switches between normal groovy chiptune and high-tempo panic BGM.
func (am *AudioManager) SetPanicMode(panicState bool) {
	if am.isPanic == panicState {
		return
	}
	am.isPanic = panicState
	am.StartBGM()
}

// PlaySFX plays a sound effect buffer with rate limiting to prevent choking.
func (am *AudioManager) playBuffer(name string, buf []byte, cooldown time.Duration) {
	if am.context == nil || len(buf) == 0 || am.muted {
		return
	}
	am.mu.Lock()
	defer am.mu.Unlock()

	now := time.Now()
	if last, ok := am.sfxCooldown[name]; ok {
		if now.Sub(last) < cooldown {
			return
		}
	}
	am.sfxCooldown[name] = now

	p, err := am.context.NewPlayer(bytes.NewReader(buf))
	if err == nil {
		p.Play()
	}
}

func (am *AudioManager) PlayWarning()    { am.playBuffer("warning", am.sfxWarning, 300*time.Millisecond) }
func (am *AudioManager) PlayAlarm()      { am.playBuffer("alarm", am.sfxAlarm, 250*time.Millisecond) }
func (am *AudioManager) PlaySteam()      { am.playBuffer("steam", am.sfxSteam, 150*time.Millisecond) }
func (am *AudioManager) PlayExtinguish() { am.playBuffer("extinguish", am.sfxExtinguish, 120*time.Millisecond) }
func (am *AudioManager) PlayClutch()     { am.playBuffer("clutch", am.sfxClutch, 200*time.Millisecond) }
func (am *AudioManager) PlayExplosion()  { am.playBuffer("explosion", am.sfxExplosion, 200*time.Millisecond) }
func (am *AudioManager) PlayCat()        { am.playBuffer("cat", am.sfxCat, 400*time.Millisecond) }
func (am *AudioManager) PlayPickup()     { am.playBuffer("pickup", am.sfxPickup, 80*time.Millisecond) }

// ============================================================================
// DSP SYNTHESIS HELPERS
// ============================================================================

func synthesizeTone(duration, startFreq, endFreq, duty float64, waveType string, vol float64) []byte {
	numSamples := int(duration * SampleRate)
	buf := make([]byte, numSamples*BytesPerSample)

	var phase float64
	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(numSamples)
		curFreq := startFreq + (endFreq-startFreq)*t
		phase += curFreq / SampleRate
		if phase >= 1.0 {
			phase -= math.Floor(phase)
		}

		var s float64
		switch waveType {
		case "square":
			if phase < duty {
				s = 1.0
			} else {
				s = -1.0
			}
		case "sawtooth":
			s = 2.0*phase - 1.0
		case "triangle":
			s = 4.0*math.Abs(phase-0.5) - 1.0
		default: // sine
			s = math.Sin(2.0 * math.Pi * phase)
		}

		// Enveloping (Attack + Release ramp)
		env := 1.0
		ramp := 0.01 * SampleRate
		if float64(i) < ramp {
			env = float64(i) / ramp
		} else if float64(numSamples-i) < ramp*2 {
			env = float64(numSamples-i) / (ramp * 2)
		}

		val := int16(s * vol * env * 32767.0)
		offset := i * BytesPerSample
		binary.LittleEndian.PutUint16(buf[offset:], uint16(val))
		binary.LittleEndian.PutUint16(buf[offset+2:], uint16(val))
	}
	return buf
}

func synthesizeNoise(duration, vol, filter float64) []byte {
	numSamples := int(duration * SampleRate)
	buf := make([]byte, numSamples*BytesPerSample)

	var lastVal float64
	for i := 0; i < numSamples; i++ {
		raw := (rand.Float64()*2.0 - 1.0)
		filtered := lastVal + filter*(raw-lastVal)
		lastVal = filtered

		// Fade out envelope
		fade := 1.0 - (float64(i) / float64(numSamples))
		val := int16(filtered * vol * fade * 32767.0)

		offset := i * BytesPerSample
		binary.LittleEndian.PutUint16(buf[offset:], uint16(val))
		binary.LittleEndian.PutUint16(buf[offset+2:], uint16(val))
	}
	return buf
}

func synthesizeExplosion(duration float64) []byte {
	numSamples := int(duration * SampleRate)
	buf := make([]byte, numSamples*BytesPerSample)

	var filterVal float64
	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(numSamples)
		raw := rand.Float64()*2.0 - 1.0
		// Low-pass filter that gets progressively heavier to simulate sub-bass rumble
		cutoff := 0.4 * (1.0 - t*0.8)
		filterVal += cutoff * (raw - filterVal)

		// Steep exponential decay
		decay := math.Exp(-t * 5.0)
		val := int16(filterVal * decay * 32767.0)

		offset := i * BytesPerSample
		binary.LittleEndian.PutUint16(buf[offset:], uint16(val))
		binary.LittleEndian.PutUint16(buf[offset+2:], uint16(val))
	}
	return buf
}

func synthesizeBGM(panicState bool) []byte {
	bpm := 125.0
	if panicState {
		bpm = 155.0
	}
	beatSec := 60.0 / bpm
	barSec := beatSec * 4
	totalSec := barSec * 4 // 4-bar loop
	numSamples := int(totalSec * SampleRate)
	buf := make([]byte, numSamples*BytesPerSample)

	// Scale notes: C, E, G, A, Bb
	bassNotes := []float64{130.81, 164.81, 196.00, 220.00} // C3, E3, G3, A3
	leadNotes := []float64{261.63, 329.63, 392.00, 440.00, 523.25}

	for i := 0; i < numSamples; i++ {
		sec := float64(i) / SampleRate
		curBeat := math.Mod(sec, beatSec) / beatSec

		// 1. Driving Bassline (Every quarter note)
		beatIdx := int(sec/beatSec) % len(bassNotes)
		freqBass := bassNotes[beatIdx]
		phaseBass := math.Mod(sec*freqBass, 1.0)
		bass := 0.0
		if phaseBass < 0.5 {
			bass = 1.0
		} else {
			bass = -1.0
		}
		bassEnv := 1.0 - curBeat*0.8

		// 2. Chiptune Arp/Lead
		leadIdx := int(sec/(beatSec*0.5)) % len(leadNotes)
		freqLead := leadNotes[leadIdx]
		if panicState {
			freqLead *= 1.5 // higher pitch in panic
		}
		phaseLead := math.Mod(sec*freqLead, 1.0)
		lead := 2.0*phaseLead - 1.0 // sawtooth
		leadEnv := 0.6

		// 3. Hi-Hat noise pulse on 8th notes
		hatSec := math.Mod(sec, beatSec*0.5)
		hatEnv := math.Max(0, 1.0-hatSec*30.0)
		hat := (rand.Float64()*2.0 - 1.0) * hatEnv * 0.2

		mix := (bass*0.25*bassEnv + lead*0.15*leadEnv + hat)
		if panicState {
			mix *= 1.2
		}

		// Clamp
		if mix > 1.0 {
			mix = 1.0
		} else if mix < -1.0 {
			mix = -1.0
		}

		val := int16(mix * 22000.0)
		offset := i * BytesPerSample
		binary.LittleEndian.PutUint16(buf[offset:], uint16(val))
		binary.LittleEndian.PutUint16(buf[offset+2:], uint16(val))
	}
	return buf
}
