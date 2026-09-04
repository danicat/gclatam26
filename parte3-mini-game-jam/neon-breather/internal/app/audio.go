package app

import (
	"bytes"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"

	"panic-recover/internal/game"
	"panic-recover/internal/sound"
)

var gameplayEffects = [...]sound.Effect{
	sound.EffectPanic,
	sound.EffectForcedPanic,
	sound.EffectElimination,
	sound.EffectRecover,
	sound.EffectCritical,
	sound.EffectVictory,
	sound.EffectGameOver,
}

type soundSystem struct {
	players [len(gameplayEffects)]*audio.Player
	bgm     *audio.Player
}

func newSoundSystem(assets spriteAssets) (*soundSystem, error) {
	context := audio.NewContext(sound.SampleRate)
	system := &soundSystem{}
	for i, effect := range gameplayEffects {
		player, err := context.NewPlayer(bytes.NewReader(sound.GenerateEffect(effect)))
		if err != nil {
			return nil, fmt.Errorf("create %s player: %w", effect, err)
		}
		system.players[i] = player
	}
	if len(assets.music) > 0 {
		stream, err := mp3.Decode(context, bytes.NewReader(assets.music))
		if err == nil {
			loop := audio.NewInfiniteLoop(stream, stream.Length())
			if player, playerErr := context.NewPlayer(loop); playerErr == nil {
				player.SetVolume(0.32)
				system.bgm = player
				player.Play()
			}
		}
	}
	return system, nil
}

func (s *soundSystem) play(effect sound.Effect) error {
	for i, candidate := range gameplayEffects {
		if candidate != effect {
			continue
		}
		player := s.players[i]
		if player == nil {
			return fmt.Errorf("player for %s is nil", effect)
		}
		if err := player.Rewind(); err != nil {
			return fmt.Errorf("rewind %s player: %w", effect, err)
		}
		player.Play()
		return nil
	}
	return fmt.Errorf("unknown sound effect %q", effect)
}

func effectForPanicTransition(forced bool) sound.Effect {
	if forced {
		return sound.EffectForcedPanic
	}
	return sound.EffectPanic
}

func effectForPhase(phase game.Phase) sound.Effect {
	if phase == game.PhaseRecoverAvailable {
		return sound.EffectRecover
	}
	return sound.EffectPanic
}
