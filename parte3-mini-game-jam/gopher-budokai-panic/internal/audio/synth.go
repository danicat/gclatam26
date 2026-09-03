package audio

import (
	"encoding/binary"
	"math"
	"math/rand"
)

const (
	SampleRate     = 44100
	NumChannels    = 2
	BytesPerSample = 2
)

// NoteDef defines parameters for a synthesized note or sound sample.
type NoteDef struct {
	WaveType     string  // "sawtooth", "square", "sine", "triangle", "noise"
	Duration     float64 // in seconds
	StartFreq    float64 // in Hz
	EndFreq      float64 // in Hz
	DutyCycle    float64 // for square wave (0.1 to 0.9, default 0.5)
	VibratoFreq  float64 // in Hz
	VibratoDepth float64 // in Hz
	NoiseFilter  float64 // low-pass filter coefficient (0.0 to 1.0)
	Volume       float64 // 0.0 to 1.0
	Attack       float64 // in seconds
	Decay        float64 // in seconds
	Sustain      float64 // 0.0 to 1.0 amplitude
	Release      float64 // in seconds
	Pan          float64 // -1.0 (left) to 1.0 (right)
}

// TrackDef represents an instrument track containing a sequence of notes.
type TrackDef struct {
	Name  string
	Pan   float64
	Notes []NoteDef
}

// clamp16 clamps an int32 value to the 16-bit signed range [-32768, 32767].
func clamp16(v int32) int16 {
	if v > 32767 {
		return 32767
	}
	if v < -32768 {
		return -32768
	}
	return int16(v)
}

// SynthesizeNote generates 16-bit 44.1kHz stereo PCM for a single note definition.
func SynthesizeNote(n NoteDef) []byte {
	numSamples := int(n.Duration * float64(SampleRate))
	if numSamples <= 0 {
		return nil
	}

	buf := make([]byte, numSamples*NumChannels*BytesPerSample)

	phase := 0.0
	vibratoPhase := 0.0
	filterState := 0.0
	rng := rand.New(rand.NewSource(1337))

	duty := n.DutyCycle
	if duty <= 0 || duty >= 1 {
		duty = 0.5
	}

	leftGain := math.Cos((n.Pan + 1.0) * math.Pi / 4.0)
	rightGain := math.Sin((n.Pan + 1.0) * math.Pi / 4.0)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(SampleRate)
		progress := t / n.Duration

		// Pitch interpolation with vibrato
		baseFreq := n.StartFreq + (n.EndFreq-n.StartFreq)*progress
		if n.VibratoDepth > 0 && n.VibratoFreq > 0 {
			vibratoPhase += 2.0 * math.Pi * n.VibratoFreq / float64(SampleRate)
			baseFreq += math.Sin(vibratoPhase) * n.VibratoDepth
		}
		if baseFreq < 10.0 {
			baseFreq = 10.0
		}

		// Oscillator
		phase += 2.0 * math.Pi * baseFreq / float64(SampleRate)
		for phase >= 2.0*math.Pi {
			phase -= 2.0 * math.Pi
		}

		var sample float64
		switch n.WaveType {
		case "sawtooth":
			sample = (phase / math.Pi) - 1.0
		case "square":
			if phase < 2.0*math.Pi*duty {
				sample = 0.8
			} else {
				sample = -0.8
			}
		case "triangle":
			if phase < math.Pi {
				sample = -1.0 + (2.0 * phase / math.Pi)
			} else {
				sample = 3.0 - (2.0 * phase / math.Pi)
			}
		case "noise":
			white := rng.Float64()*2.0 - 1.0
			coeff := n.NoiseFilter
			if coeff <= 0 {
				coeff = 0.3
			}
			filterState += coeff * (white - filterState)
			sample = filterState
		case "sine":
			fallthrough
		default:
			sample = math.Sin(phase)
		}

		// ADSR Envelope
		env := 0.0
		attackTime := n.Attack
		decayTime := n.Decay
		releaseTime := n.Release
		sustainLevel := n.Sustain

		// Ensure envelope fits within note duration
		if attackTime+decayTime+releaseTime > n.Duration {
			scale := n.Duration / (attackTime + decayTime + releaseTime + 0.0001)
			attackTime *= scale
			decayTime *= scale
			releaseTime *= scale
		}

		sustainDuration := n.Duration - attackTime - decayTime - releaseTime
		if sustainDuration < 0 {
			sustainDuration = 0
		}

		if t < attackTime {
			if attackTime > 0 {
				env = t / attackTime
			} else {
				env = 1.0
			}
		} else if t < attackTime+decayTime {
			dt := t - attackTime
			env = 1.0 - (1.0-sustainLevel)*(dt/decayTime)
		} else if t < attackTime+decayTime+sustainDuration {
			env = sustainLevel
		} else {
			rt := t - (attackTime + decayTime + sustainDuration)
			if releaseTime > 0 {
				env = sustainLevel * (1.0 - rt/releaseTime)
			} else {
				env = 0.0
			}
		}
		if env < 0 {
			env = 0
		}

		sampleVal := sample * env * n.Volume
		valL := int32(sampleVal * leftGain * 32767.0)
		valR := int32(sampleVal * rightGain * 32767.0)

		s16L := clamp16(valL)
		s16R := clamp16(valR)

		idx := i * NumChannels * BytesPerSample
		binary.LittleEndian.PutUint16(buf[idx:idx+2], uint16(s16L))
		binary.LittleEndian.PutUint16(buf[idx+2:idx+4], uint16(s16R))
	}

	return buf
}

// MixTracks mixes multiple TrackDefs into a single continuous 16-bit 44.1kHz stereo PCM buffer.
// Mixes in int32 to prevent overflow and applies master volume and hard clamping.
func MixTracks(tracks []TrackDef, masterVolume float64) []byte {
	type renderedTrack struct {
		data []byte
	}

	var rendered []renderedTrack
	maxLen := 0

	for _, track := range tracks {
		var trackBuf []byte
		for _, note := range track.Notes {
			if track.Pan != 0 && note.Pan == 0 {
				note.Pan = track.Pan
			}
			pcm := SynthesizeNote(note)
			trackBuf = append(trackBuf, pcm...)
		}
		if len(trackBuf) > maxLen {
			maxLen = len(trackBuf)
		}
		rendered = append(rendered, renderedTrack{data: trackBuf})
	}

	if maxLen == 0 {
		return nil
	}

	totalSamples := maxLen / (NumChannels * BytesPerSample)
	out := make([]byte, maxLen)

	for s := 0; s < totalSamples; s++ {
		offset := s * NumChannels * BytesPerSample
		var sumL int32 = 0
		var sumR int32 = 0

		for _, r := range rendered {
			if offset+4 <= len(r.data) {
				sL := int32(int16(binary.LittleEndian.Uint16(r.data[offset : offset+2])))
				sR := int32(int16(binary.LittleEndian.Uint16(r.data[offset+2 : offset+4])))
				sumL += sL
				sumR += sR
			}
		}

		finalL := clamp16(int32(float64(sumL) * masterVolume))
		finalR := clamp16(int32(float64(sumR) * masterVolume))

		binary.LittleEndian.PutUint16(out[offset:offset+2], uint16(finalL))
		binary.LittleEndian.PutUint16(out[offset+2:offset+4], uint16(finalR))
	}

	return out
}
