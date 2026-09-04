package audio

import (
	"bytes"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const sampleRate = 44100

type AudioManager struct {
	context   *audio.Context
	laserData []byte
	boomData  []byte
	hitData   []byte
	powerData []byte
}

var GlobalAudio *AudioManager

func InitAudio() {
	ctx := audio.NewContext(sampleRate)
	GlobalAudio = &AudioManager{
		context:   ctx,
		laserData: generateLaser(),
		boomData:  generateExplosion(),
		hitData:   generateHit(),
		powerData: generatePowerup(),
	}
}

func (a *AudioManager) PlayLaser() {
	if a == nil || a.context == nil {
		return
	}
	p := a.context.NewPlayerFromBytes(a.laserData)
	p.Play()
}

func (a *AudioManager) PlayExplosion() {
	if a == nil || a.context == nil {
		return
	}
	p := a.context.NewPlayerFromBytes(a.boomData)
	p.Play()
}

func (a *AudioManager) PlayHit() {
	if a == nil || a.context == nil {
		return
	}
	p := a.context.NewPlayerFromBytes(a.hitData)
	p.Play()
}

func (a *AudioManager) PlayPowerup() {
	if a == nil || a.context == nil {
		return
	}
	p := a.context.NewPlayerFromBytes(a.powerData)
	p.Play()
}

// Generate 16-bit mono PCM bytes for simple retro SFX
func generateLaser() []byte {
	duration := 0.15
	numSamples := int(float64(sampleRate) * duration)
	buf := new(bytes.Buffer)
	startFreq := 880.0
	endFreq := 220.0

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(sampleRate)
		progress := float64(i) / float64(numSamples)
		freq := startFreq + (endFreq-startFreq)*progress
		val := math.Sin(2.0 * math.Pi * freq * t)
		// 16-bit signed integer
		sample := int16(val * 0.25 * float64(math.MaxInt16))
		buf.WriteByte(byte(sample))
		buf.WriteByte(byte(sample >> 8))
	}
	return buf.Bytes()
}

func generateExplosion() []byte {
	duration := 0.3
	numSamples := int(float64(sampleRate) * duration)
	buf := new(bytes.Buffer)

	for i := 0; i < numSamples; i++ {
		decay := 1.0 - float64(i)/float64(numSamples)
		noise := (rand.Float64()*2.0 - 1.0) * decay
		sample := int16(noise * 0.3 * float64(math.MaxInt16))
		buf.WriteByte(byte(sample))
		buf.WriteByte(byte(sample >> 8))
	}
	return buf.Bytes()
}

func generateHit() []byte {
	duration := 0.08
	numSamples := int(float64(sampleRate) * duration)
	buf := new(bytes.Buffer)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(sampleRate)
		val := math.Sin(2.0 * math.Pi * 150.0 * t)
		sample := int16(val * 0.3 * float64(math.MaxInt16))
		buf.WriteByte(byte(sample))
		buf.WriteByte(byte(sample >> 8))
	}
	return buf.Bytes()
}

func generatePowerup() []byte {
	duration := 0.25
	numSamples := int(float64(sampleRate) * duration)
	buf := new(bytes.Buffer)

	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(sampleRate)
		progress := float64(i) / float64(numSamples)
		freq := 400.0 + 800.0*progress
		val := math.Sin(2.0 * math.Pi * freq * t)
		sample := int16(val * 0.25 * float64(math.MaxInt16))
		buf.WriteByte(byte(sample))
		buf.WriteByte(byte(sample >> 8))
	}
	return buf.Bytes()
}
