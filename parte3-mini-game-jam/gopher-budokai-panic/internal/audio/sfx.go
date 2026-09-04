package audio

// BuildBlastSFX generates a quick DBZ ki blast sound effect.
func BuildBlastSFX() []byte {
	return SynthesizeNote(NoteDef{
		WaveType:  "square",
		DutyCycle: 0.3,
		Duration:  0.12,
		StartFreq: 850.0,
		EndFreq:   160.0,
		Volume:    0.35,
		Attack:    0.005,
		Decay:     0.05,
		Sustain:   0.3,
		Release:   0.06,
		Pan:       0.0,
	})
}

// BuildHitSFX generates a punch/kick melee impact sound effect.
func BuildHitSFX() []byte {
	noiseNote := NoteDef{
		WaveType:    "noise",
		NoiseFilter: 0.5,
		Duration:    0.10,
		StartFreq:   220.0,
		EndFreq:     60.0,
		Volume:      0.45,
		Attack:      0.003,
		Decay:       0.04,
		Sustain:     0.1,
		Release:     0.05,
		Pan:         0.0,
	}
	thudNote := NoteDef{
		WaveType:  "sine",
		Duration:  0.10,
		StartFreq: 140.0,
		EndFreq:   45.0,
		Volume:    0.50,
		Attack:    0.005,
		Decay:     0.04,
		Sustain:   0.2,
		Release:   0.05,
		Pan:       0.0,
	}
	return MixTracks([]TrackDef{
		{Notes: []NoteDef{noiseNote}},
		{Notes: []NoteDef{thudNote}},
	}, 0.85)
}

// BuildChargeSFX generates the rising ki aura charging hum.
func BuildChargeSFX() []byte {
	return SynthesizeNote(NoteDef{
		WaveType:     "sawtooth",
		Duration:     0.35,
		StartFreq:    110.0,
		EndFreq:      320.0,
		VibratoFreq:  14.0,
		VibratoDepth: 10.0,
		Volume:       0.25,
		Attack:       0.05,
		Decay:        0.10,
		Sustain:      0.8,
		Release:      0.10,
		Pan:          0.0,
	})
}

// BuildBeamSFX generates the massive energy roar of a Super Beam (Kamehameha / Final Flash).
func BuildBeamSFX() []byte {
	beam1 := NoteDef{
		WaveType:     "sawtooth",
		Duration:     0.65,
		StartFreq:    350.0,
		EndFreq:      180.0,
		VibratoFreq:  18.0,
		VibratoDepth: 15.0,
		Volume:       0.40,
		Attack:       0.02,
		Decay:        0.15,
		Sustain:      0.8,
		Release:      0.15,
		Pan:          0.0,
	}
	beam2 := NoteDef{
		WaveType:    "noise",
		NoiseFilter: 0.7,
		Duration:    0.65,
		StartFreq:   500.0,
		EndFreq:     120.0,
		Volume:      0.35,
		Attack:      0.03,
		Decay:       0.20,
		Sustain:     0.7,
		Release:     0.15,
		Pan:         0.0,
	}
	return MixTracks([]TrackDef{
		{Notes: []NoteDef{beam1}},
		{Notes: []NoteDef{beam2}},
	}, 0.85)
}

// BuildVanishSFX generates the Instant Transmission / Vanish "shwing" effect.
func BuildVanishSFX() []byte {
	return SynthesizeNote(NoteDef{
		WaveType:  "sine",
		Duration:  0.14,
		StartFreq: 1400.0,
		EndFreq:   2800.0,
		Volume:    0.40,
		Attack:    0.005,
		Decay:     0.04,
		Sustain:   0.4,
		Release:   0.09,
		Pan:       0.0,
	})
}

// BuildPanicAlertSFX generates the emergency warning chirp when entering PANIC!.
func BuildPanicAlertSFX() []byte {
	n1 := NoteDef{
		WaveType:  "square",
		DutyCycle: 0.5,
		Duration:  0.10,
		StartFreq: 950.0,
		EndFreq:   950.0,
		Volume:    0.35,
		Attack:    0.005,
		Decay:     0.03,
		Sustain:   0.7,
		Release:   0.02,
		Pan:       -0.2,
	}
	n2 := NoteDef{
		WaveType:  "square",
		DutyCycle: 0.5,
		Duration:  0.12,
		StartFreq: 1350.0,
		EndFreq:   1350.0,
		Volume:    0.40,
		Attack:    0.005,
		Decay:     0.03,
		Sustain:   0.7,
		Release:   0.03,
		Pan:       0.2,
	}
	return MixTracks([]TrackDef{
		{Notes: []NoteDef{n1, n2}},
	}, 0.8)
}

// BuildRecoverKiaiSFX generates the explosive Kiai shockwave that clears PANIC! and blows away opponents.
func BuildRecoverKiaiSFX() []byte {
	shockwave := NoteDef{
		WaveType:  "sine",
		Duration:  0.35,
		StartFreq: 260.0,
		EndFreq:   40.0,
		Volume:    0.60,
		Attack:    0.005,
		Decay:     0.08,
		Sustain:   0.3,
		Release:   0.20,
		Pan:       0.0,
	}
	noiseBurst := NoteDef{
		WaveType:    "noise",
		NoiseFilter: 0.8,
		Duration:    0.30,
		StartFreq:   600.0,
		EndFreq:     100.0,
		Volume:      0.45,
		Attack:      0.005,
		Decay:       0.09,
		Sustain:     0.2,
		Release:     0.15,
		Pan:         0.0,
	}
	return MixTracks([]TrackDef{
		{Notes: []NoteDef{shockwave}},
		{Notes: []NoteDef{noiseBurst}},
	}, 0.9)
}

// BuildClashSFX generates an intense energy clash spark.
func BuildClashSFX() []byte {
	return SynthesizeNote(NoteDef{
		WaveType:    "noise",
		NoiseFilter: 0.9,
		Duration:    0.08,
		StartFreq:   1200.0,
		EndFreq:     400.0,
		Volume:      0.30,
		Attack:      0.002,
		Decay:       0.03,
		Sustain:     0.2,
		Release:     0.04,
		Pan:         0.0,
	})
}
