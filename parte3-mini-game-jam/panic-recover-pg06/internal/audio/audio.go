package audio

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/rand"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	SampleRate = 44100
)

// SoundSystem manages all procedural DSP sound synthesis and playback.
type SoundSystem struct {
	mu           sync.Mutex
	audioCtx     *audio.Context
	bgmPlayer    *audio.Player
	bgmLoop      *audio.InfiniteLoop
	muted        bool
	masterVolume float64

	keyClickPCM   []byte
	selectPCM     []byte
	editPCM       []byte
	warningPCM    []byte
	recoverPCM    []byte
	panicPCM      []byte
	bgmPCM        []byte
}

var (
	globalSoundSys *SoundSystem
	once           sync.Once
)

// GetSoundSystem returns the singleton SoundSystem instance.
func GetSoundSystem() *SoundSystem {
	once.Do(func() {
		globalSoundSys = NewSoundSystem()
	})
	return globalSoundSys
}

// NewSoundSystem creates and pre-synthesizes all procedural audio into memory.
func NewSoundSystem() *SoundSystem {
	ctx := audio.NewContext(SampleRate)
	ss := &SoundSystem{
		audioCtx:     ctx,
		masterVolume: 0.8,
	}

	// Pre-synthesize all sound buffers
	ss.keyClickPCM = synthKeyClick()
	ss.selectPCM = synthTone(720, 0.04, "sine", 0.3)
	ss.editPCM = synthDoubleTone(523.25, 659.25, 0.08, 0.4)
	ss.warningPCM = synthSquareSiren(880, 440, 0.18, 0.5)
	ss.recoverPCM = synthRecoverFanfare()
	ss.panicPCM = synthPanicCrash()
	ss.bgmPCM = synthCyberBGM()

	return ss
}

// ToggleMute toggles mute on/off.
func (s *SoundSystem) ToggleMute() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.muted = !s.muted
	if s.bgmPlayer != nil {
		if s.muted {
			s.bgmPlayer.SetVolume(0)
		} else {
			s.bgmPlayer.SetVolume(0.4 * s.masterVolume)
		}
	}
	return s.muted
}

// IsMuted returns current mute state.
func (s *SoundSystem) IsMuted() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.muted
}

func (s *SoundSystem) playBytes(pcm []byte, volume float64) {
	if s.muted || len(pcm) == 0 {
		return
	}
	p := s.audioCtx.NewPlayerFromBytes(pcm)
	if p != nil {
		p.SetVolume(volume * s.masterVolume)
		p.Play()
	}
}

// PlayKeyClick plays a mechanical keyboard click.
func (s *SoundSystem) PlayKeyClick() {
	s.playBytes(s.keyClickPCM, 0.35)
}

// PlaySelect plays a navigation beep.
func (s *SoundSystem) PlaySelect() {
	s.playBytes(s.selectPCM, 0.3)
}

// PlayEditMode plays an edit enter chime.
func (s *SoundSystem) PlayEditMode() {
	s.playBytes(s.editPCM, 0.4)
}

// PlayWarning plays an urgent warning alert tone.
func (s *SoundSystem) PlayWarning() {
	s.playBytes(s.warningPCM, 0.5)
}

// PlayRecover plays the triumphant recovery fanfare.
func (s *SoundSystem) PlayRecover() {
	s.playBytes(s.recoverPCM, 0.7)
}

// PlayPanic plays the fatal crash explosion.
func (s *SoundSystem) PlayPanic() {
	s.playBytes(s.panicPCM, 0.8)
}

// StartBGM starts the looping ambient cyberpunk background music.
func (s *SoundSystem) StartBGM() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.bgmPlayer != nil && s.bgmPlayer.IsPlaying() {
		return
	}
	if len(s.bgmPCM) == 0 {
		return
	}

	r := bytes.NewReader(s.bgmPCM)
	s.bgmLoop = audio.NewInfiniteLoop(r, int64(len(s.bgmPCM)))
	var err error
	s.bgmPlayer, err = s.audioCtx.NewPlayer(s.bgmLoop)
	if err == nil {
		vol := 0.35 * s.masterVolume
		if s.muted {
			vol = 0
		}
		s.bgmPlayer.SetVolume(vol)
		s.bgmPlayer.Play()
	}
}

// StopBGM pauses/stops the background music.
func (s *SoundSystem) StopBGM() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.bgmPlayer != nil {
		s.bgmPlayer.Pause()
	}
}

// ============================================================================
// DSP Synthesis Algorithms
// ============================================================================

// writeStereoSample writes a 16-bit stereo sample to buf.
func writeStereoSample(buf *bytes.Buffer, left, right float64) {
	if left > 1.0 {
		left = 1.0
	} else if left < -1.0 {
		left = -1.0
	}
	if right > 1.0 {
		right = 1.0
	} else if right < -1.0 {
		right = -1.0
	}
	l16 := int16(left * 32767.0)
	r16 := int16(right * 32767.0)
	_ = binary.Write(buf, binary.LittleEndian, l16)
	_ = binary.Write(buf, binary.LittleEndian, r16)
}

// synthKeyClick generates a crisp mechanical typing click.
func synthKeyClick() []byte {
	duration := 0.025
	totalSamples := int(duration * SampleRate)
	buf := new(bytes.Buffer)
	r := rand.New(rand.NewSource(42))

	for i := 0; i < totalSamples; i++ {
		t := float64(i) / float64(SampleRate)
		env := math.Exp(-t * 220.0) // Fast exponential decay
		// Blend high-frequency click with short noise
		tone := math.Sin(2.0*math.Pi*2400.0*t) * 0.4
		noise := (r.Float64()*2.0 - 1.0) * 0.6
		sample := (tone + noise) * env
		writeStereoSample(buf, sample, sample)
	}
	return buf.Bytes()
}

// synthTone creates a single tone with specified waveform.
func synthTone(freq, duration float64, wave string, volume float64) []byte {
	totalSamples := int(duration * SampleRate)
	buf := new(bytes.Buffer)

	for i := 0; i < totalSamples; i++ {
		t := float64(i) / float64(SampleRate)
		env := 1.0 - (float64(i) / float64(totalSamples)) // Linear fade
		var sample float64
		switch wave {
		case "square":
			if math.Sin(2.0*math.Pi*freq*t) >= 0 {
				sample = 1.0
			} else {
				sample = -1.0
			}
		case "saw":
			sample = 2.0*(t*freq-math.Floor(0.5+t*freq))
		default: // sine
			sample = math.Sin(2.0 * math.Pi * freq * t)
		}
		sample *= env * volume
		writeStereoSample(buf, sample, sample)
	}
	return buf.Bytes()
}

// synthDoubleTone creates a sequential dual frequency chime.
func synthDoubleTone(f1, f2, duration float64, volume float64) []byte {
	halfSamples := int(duration * 0.5 * SampleRate)
	buf := new(bytes.Buffer)

	// Note 1
	for i := 0; i < halfSamples; i++ {
		t := float64(i) / float64(SampleRate)
		env := 1.0 - (float64(i) / float64(halfSamples))
		s := math.Sin(2.0*math.Pi*f1*t) * env * volume
		writeStereoSample(buf, s, s)
	}
	// Note 2
	for i := 0; i < halfSamples; i++ {
		t := float64(i) / float64(SampleRate)
		env := 1.0 - (float64(i) / float64(halfSamples))
		s := math.Sin(2.0*math.Pi*f2*t) * env * volume
		writeStereoSample(buf, s, s)
	}
	return buf.Bytes()
}

// synthSquareSiren creates an alternating square wave siren.
func synthSquareSiren(f1, f2, duration float64, volume float64) []byte {
	totalSamples := int(duration * SampleRate)
	buf := new(bytes.Buffer)

	for i := 0; i < totalSamples; i++ {
		t := float64(i) / float64(SampleRate)
		freq := f1
		if math.Mod(t, 0.08) > 0.04 {
			freq = f2
		}
		var s float64
		if math.Sin(2.0*math.Pi*freq*t) >= 0 {
			s = 0.8
		} else {
			s = -0.8
		}
		writeStereoSample(buf, s*volume, s*volume)
	}
	return buf.Bytes()
}

// synthRecoverFanfare synthesizes a rising major arpeggio C5 -> E5 -> G5 -> C6.
func synthRecoverFanfare() []byte {
	notes := []float64{523.25, 659.25, 783.99, 1046.50} // C5, E5, G5, C6
	noteDur := 0.12
	buf := new(bytes.Buffer)

	for _, freq := range notes {
		samples := int(noteDur * SampleRate)
		for i := 0; i < samples; i++ {
			t := float64(i) / float64(SampleRate)
			// Smooth ADSR
			env := math.Min(1.0, float64(i)/(float64(samples)*0.1)) * (1.0 - float64(i)/float64(samples))
			// Sine fundamental + soft harmonic
			s := (math.Sin(2*math.Pi*freq*t) + 0.3*math.Sin(4*math.Pi*freq*t)) * env * 0.6
			writeStereoSample(buf, s, s)
		}
	}
	return buf.Bytes()
}

// synthPanicCrash synthesizes a dramatic descending noise crash.
func synthPanicCrash() []byte {
	duration := 0.6
	totalSamples := int(duration * SampleRate)
	buf := new(bytes.Buffer)
	r := rand.New(rand.NewSource(99))

	for i := 0; i < totalSamples; i++ {
		t := float64(i) / float64(SampleRate)
		progress := float64(i) / float64(totalSamples)
		env := math.Exp(-progress * 4.0)

		// Descending pitch
		freq := 180.0 * (1.0 - progress*0.8)
		saw := 2.0 * (t*freq - math.Floor(0.5+t*freq))
		noise := (r.Float64()*2.0 - 1.0) * 0.7

		sample := (saw*0.5 + noise*0.5) * env * 0.8
		writeStereoSample(buf, sample, sample)
	}
	return buf.Bytes()
}

// synthCyberBGM synthesizes a seamless 4-measure looping cyberpunk bass/arpeggio track.
func synthCyberBGM() []byte {
	bpm := 120.0
	beatSec := 60.0 / bpm
	totalMeasures := 2
	totalDuration := float64(totalMeasures) * 4.0 * beatSec // ~4.0 seconds seamless loop
	totalSamples := int(totalDuration * SampleRate)
	buf := new(bytes.Buffer)

	// Bass notes sequence (A1, C2, D2, F1)
	bassFreqs := []float64{110.0, 130.81, 146.83, 87.31}
	// Arp notes (A3, C4, E4, A4)
	arpFreqs := []float64{220.0, 261.63, 329.63, 440.0}

	for i := 0; i < totalSamples; i++ {
		t := float64(i) / float64(SampleRate)

		// Bass line (changes every half measure)
		bassIdx := int(t/(beatSec*2.0)) % len(bassFreqs)
		bFreq := bassFreqs[bassIdx]
		// Sawtooth bass
		bassSample := (2.0*(t*bFreq-math.Floor(0.5+t*bFreq))) * 0.25

		// 16th note Arpeggiator (16 notes per measure = 4 notes per beat)
		arpStep := int(t / (beatSec / 4.0))
		aFreq := arpFreqs[arpStep%len(arpFreqs)]
		stepTime := math.Mod(t, beatSec/4.0)
		arpEnv := math.Exp(-stepTime * 18.0)
		arpSample := math.Sin(2.0*math.Pi*aFreq*t) * arpEnv * 0.2

		// Subtle pulse kick on every beat
		beatProgress := math.Mod(t, beatSec)
		kickEnv := math.Exp(-beatProgress * 25.0)
		kickFreq := 120.0 * (1.0 - beatProgress*3.0)
		if kickFreq < 40.0 {
			kickFreq = 40.0
		}
		kickSample := math.Sin(2.0*math.Pi*kickFreq*beatProgress) * kickEnv * 0.35

		mix := (bassSample + arpSample + kickSample) * 0.7
		writeStereoSample(buf, mix, mix)
	}

	return buf.Bytes()
}
