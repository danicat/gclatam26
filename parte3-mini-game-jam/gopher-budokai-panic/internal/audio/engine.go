package audio

import (
	"bytes"
	"sync"
	"time"

	ebitenaudio "github.com/hajimehoshi/ebiten/v2/audio"
)

type Engine struct {
	mu           sync.Mutex
	audioCtx     *ebitenaudio.Context
	bgmPlayer    *ebitenaudio.Player
	currentBGM   string
	masterVolume float64

	// Pre-rendered BGM buffers
	battleBGM []byte
	panicBGM  []byte
	titleBGM  []byte

	// Pre-rendered SFX buffers
	blastSFX   []byte
	hitSFX     []byte
	chargeSFX  []byte
	beamSFX    []byte
	vanishSFX  []byte
	panicSFX   []byte
	recoverSFX []byte
	clashSFX   []byte

	// Rate limiting for high-frequency SFX
	lastBlastTime  time.Time
	lastHitTime    time.Time
	lastChargeTime time.Time
}

var (
	globalEngine *Engine
	engineOnce   sync.Once
)

// Get returns the singleton audio Engine instance.
func Get() *Engine {
	engineOnce.Do(func() {
		globalEngine = newEngine()
	})
	return globalEngine
}

func newEngine() *Engine {
	ctx := ebitenaudio.NewContext(SampleRate)
	e := &Engine{
		audioCtx:     ctx,
		masterVolume: 1.0,
	}
	e.preloadAssets()
	return e
}

// preloadAssets pre-synthesizes all audio into memory buffers on startup.
func (e *Engine) preloadAssets() {
	// Synthesize BGMs
	e.battleBGM = BuildBattleSoundtrack()
	e.panicBGM = BuildPanicSoundtrack()
	e.titleBGM = BuildTitleSoundtrack()

	// Synthesize SFXs
	e.blastSFX = BuildBlastSFX()
	e.hitSFX = BuildHitSFX()
	e.chargeSFX = BuildChargeSFX()
	e.beamSFX = BuildBeamSFX()
	e.vanishSFX = BuildVanishSFX()
	e.panicSFX = BuildPanicAlertSFX()
	e.recoverSFX = BuildRecoverKiaiSFX()
	e.clashSFX = BuildClashSFX()
}

// PlayBGM plays the specified music loop ("title", "battle", "panic").
func (e *Engine) PlayBGM(name string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.currentBGM == name && e.bgmPlayer != nil && e.bgmPlayer.IsPlaying() {
		return
	}

	if e.bgmPlayer != nil {
		e.bgmPlayer.Close()
		e.bgmPlayer = nil
	}

	var pcm []byte
	switch name {
	case "title":
		pcm = e.titleBGM
	case "panic":
		pcm = e.panicBGM
	case "battle":
		fallthrough
	default:
		pcm = e.battleBGM
	}

	if len(pcm) == 0 {
		return
	}

	loop := ebitenaudio.NewInfiniteLoop(bytes.NewReader(pcm), int64(len(pcm)))
	p, err := e.audioCtx.NewPlayer(loop)
	if err != nil {
		return
	}

	p.SetVolume(0.85 * e.masterVolume)
	p.Play()
	e.bgmPlayer = p
	e.currentBGM = name
}

// StopBGM stops current background music.
func (e *Engine) StopBGM() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.bgmPlayer != nil {
		e.bgmPlayer.Close()
		e.bgmPlayer = nil
	}
	e.currentBGM = ""
}

// playSFX plays a one-off pre-rendered PCM buffer.
func (e *Engine) playSFX(pcm []byte, volume float64) {
	if len(pcm) == 0 {
		return
	}
	p, err := e.audioCtx.NewPlayer(bytes.NewReader(pcm))
	if err != nil {
		return
	}
	p.SetVolume(volume * e.masterVolume)
	p.Play()
}

func (e *Engine) PlayBlast() {
	now := time.Now()
	if now.Sub(e.lastBlastTime) < 50*time.Millisecond {
		return
	}
	e.lastBlastTime = now
	e.playSFX(e.blastSFX, 0.7)
}

func (e *Engine) PlayHit() {
	now := time.Now()
	if now.Sub(e.lastHitTime) < 40*time.Millisecond {
		return
	}
	e.lastHitTime = now
	e.playSFX(e.hitSFX, 0.8)
}

func (e *Engine) PlayCharge() {
	now := time.Now()
	if now.Sub(e.lastChargeTime) < 250*time.Millisecond {
		return
	}
	e.lastChargeTime = now
	e.playSFX(e.chargeSFX, 0.6)
}

func (e *Engine) PlayBeam() {
	e.playSFX(e.beamSFX, 0.9)
}

func (e *Engine) PlayVanish() {
	e.playSFX(e.vanishSFX, 0.8)
}

func (e *Engine) PlayPanicAlert() {
	e.playSFX(e.panicSFX, 0.9)
}

func (e *Engine) PlayRecoverKiai() {
	e.playSFX(e.recoverSFX, 1.0)
}

func (e *Engine) PlayClash() {
	e.playSFX(e.clashSFX, 0.7)
}
