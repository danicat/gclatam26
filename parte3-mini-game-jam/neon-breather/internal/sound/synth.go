package sound

import (
	"encoding/binary"
	"math"
)

const SampleRate = 44100

type Effect string

const (
	EffectPanic       Effect = "panic"
	EffectForcedPanic Effect = "forced-panic"
	EffectElimination Effect = "elimination"
	EffectRecover     Effect = "recover"
	EffectCritical    Effect = "critical"
	EffectVictory     Effect = "victory"
	EffectGameOver    Effect = "game-over"
)

func GenerateEffect(effect Effect) []byte {
	switch effect {
	case EffectPanic:
		return toneSweep(180, 72, 0.24, 0.34)
	case EffectForcedPanic:
		return toneSweep(120, 42, 0.18, 0.34)
	case EffectElimination:
		return toneSweep(420, 1100, 0.09, 0.28)
	case EffectRecover:
		return concat(toneSweep(440, 660, 0.12, 0.25), toneSweep(660, 990, 0.16, 0.25))
	case EffectCritical:
		return noise(0.16, 0.2)
	case EffectVictory:
		return concat(toneSweep(440, 660, 0.12, 0.24), toneSweep(660, 880, 0.16, 0.24), toneSweep(880, 1320, 0.2, 0.24))
	case EffectGameOver:
		return toneSweep(360, 90, 0.42, 0.28)
	default:
		return nil
	}
}

func toneSweep(startFrequency, endFrequency, duration, volume float64) []byte {
	samples := int(math.Round(duration * SampleRate))
	pcm := make([]byte, samples*4)
	phase := 0.0
	for sample := 0; sample < samples; sample++ {
		progress := float64(sample) / float64(samples)
		frequency := startFrequency + (endFrequency-startFrequency)*progress
		phase += 2 * math.Pi * frequency / SampleRate
		envelope := envelope(progress)
		value := int32(math.Round(math.Sin(phase) * 32767 * volume * envelope))
		writeStereoSample(pcm, sample, value)
	}
	return pcm
}

func noise(duration, volume float64) []byte {
	samples := int(math.Round(duration * SampleRate))
	pcm := make([]byte, samples*4)
	state := uint32(0x12345678)
	for sample := 0; sample < samples; sample++ {
		state = state*1664525 + 1013904223
		value := int32(float64(int16(state>>16)) * volume * envelope(float64(sample)/float64(samples)))
		writeStereoSample(pcm, sample, value)
	}
	return pcm
}

func envelope(progress float64) float64 {
	if progress < 0.08 {
		return progress / 0.08
	}
	if progress > 0.82 {
		return (1 - progress) / 0.18
	}
	return 1
}

func writeStereoSample(pcm []byte, sample int, value int32) {
	if value > 32767 {
		value = 32767
	}
	if value < -32768 {
		value = -32768
	}
	encoded := uint16(int16(value))
	offset := sample * 4
	binary.LittleEndian.PutUint16(pcm[offset:offset+2], encoded)
	binary.LittleEndian.PutUint16(pcm[offset+2:offset+4], encoded)
}

func concat(parts ...[]byte) []byte {
	total := 0
	for _, part := range parts {
		total += len(part)
	}
	result := make([]byte, 0, total)
	for _, part := range parts {
		result = append(result, part...)
	}
	return result
}
