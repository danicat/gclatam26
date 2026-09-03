package audio

// Note frequencies in Hz
const (
	FreqD2  = 73.42
	FreqF2  = 87.31
	FreqG2  = 98.00
	FreqA2  = 110.00
	FreqBb2 = 116.54
	FreqC3  = 130.81
	FreqD3  = 146.83
	FreqE3  = 164.81
	FreqF3  = 174.61
	FreqG3  = 196.00
	FreqA3  = 220.00
	FreqBb3 = 233.08
	FreqC4  = 261.63
	FreqD4  = 293.66
	FreqE4  = 329.63
	FreqF4  = 349.23
	FreqG4  = 392.00
	FreqA4  = 440.00
	FreqBb4 = 466.16
	FreqC5  = 523.25
	FreqD5  = 587.33
	FreqE5  = 659.25
	FreqF5  = 698.46
	FreqG5  = 783.99
	FreqA5  = 880.00
	FreqBb5 = 932.32
)

// BuildBattleSoundtrack synthesizes the full 6-channel DBZ Budokai Tenkaichi inspired battle theme.
// 145 BPM, 16 bars looped composition in D Minor.
func BuildBattleSoundtrack() []byte {
	// Tempo: 145 BPM => 1 beat = ~0.4138s
	b := 60.0 / 145.0 // ~0.414s
	h := b / 2.0      // 1/8 note = ~0.207s
	q := b / 4.0      // 1/16 note = ~0.103s

	// 1. Lead Melody (Sawtooth with vibrato)
	var leadNotes []NoteDef
	// Phrase 1 (4 bars)
	leadSeq1 := []struct {
		freq float64
		dur  float64
	}{
		{FreqD4, h}, {FreqD4, q}, {FreqF4, q}, {FreqG4, h}, {FreqA4, h},
		{FreqBb4, h}, {FreqA4, h}, {FreqG4, h}, {FreqF4, h},
		{FreqD4, h}, {FreqF4, q}, {FreqG4, q}, {FreqA4, b},
		{FreqG4, h}, {FreqF4, h}, {FreqE4, h}, {FreqD4, h},
	}
	// Phrase 2 (4 bars, climax higher octave)
	leadSeq2 := []struct {
		freq float64
		dur  float64
	}{
		{FreqD5, h}, {FreqD5, q}, {FreqF5, q}, {FreqG5, h}, {FreqA5, h},
		{FreqBb5, h}, {FreqA5, h}, {FreqG5, h}, {FreqF5, h},
		{FreqG5, h}, {FreqA5, q}, {FreqBb5, q}, {FreqA5, b},
		{FreqF5, h}, {FreqE5, h}, {FreqD5, b},
	}

	for _, p := range leadSeq1 {
		leadNotes = append(leadNotes, NoteDef{
			WaveType:     "sawtooth",
			Duration:     p.dur,
			StartFreq:    p.freq,
			EndFreq:      p.freq,
			VibratoFreq:  6.5,
			VibratoDepth: 4.0,
			Volume:       0.28,
			Attack:       0.02,
			Decay:        0.05,
			Sustain:      0.7,
			Release:      0.04,
			Pan:          -0.2,
		})
	}
	for _, p := range leadSeq2 {
		leadNotes = append(leadNotes, NoteDef{
			WaveType:     "sawtooth",
			Duration:     p.dur,
			StartFreq:    p.freq,
			EndFreq:      p.freq,
			VibratoFreq:  7.0,
			VibratoDepth: 5.0,
			Volume:       0.30,
			Attack:       0.02,
			Decay:        0.05,
			Sustain:      0.7,
			Release:      0.04,
			Pan:          -0.2,
		})
	}

	// 2. Counter Melody (Square wave, duty 0.25)
	var counterNotes []NoteDef
	counterSeq := []struct {
		freq float64
		dur  float64
	}{
		{FreqA3, b}, {FreqD4, b}, {FreqF4, b}, {FreqE4, b},
		{FreqD4, b}, {FreqG3, b}, {FreqA3, b}, {FreqD3, b},
		{FreqF4, b}, {FreqA4, b}, {FreqD5, b}, {FreqC5, b},
		{FreqBb4, b}, {FreqA4, b}, {FreqG4, b}, {FreqA4, b},
	}
	for _, p := range counterSeq {
		counterNotes = append(counterNotes, NoteDef{
			WaveType:  "square",
			DutyCycle: 0.25,
			Duration:  p.dur,
			StartFreq: p.freq,
			EndFreq:   p.freq,
			Volume:    0.16,
			Attack:    0.03,
			Decay:     0.06,
			Sustain:   0.6,
			Release:   0.05,
			Pan:       0.3,
		})
	}

	// 3. Power Chords / Harmony Pad (Triangle + soft Saw)
	var padNotes []NoteDef
	chordSeq := []struct {
		freq float64
		dur  float64
	}{
		{FreqD3, b * 2}, {FreqF3, b * 2}, {FreqG3, b * 2}, {FreqA3, b * 2},
		{FreqBb3, b * 2}, {FreqA3, b * 2}, {FreqG3, b * 2}, {FreqD3, b * 2},
	}
	for _, c := range chordSeq {
		padNotes = append(padNotes, NoteDef{
			WaveType:  "sawtooth",
			Duration:  c.dur,
			StartFreq: c.freq,
			EndFreq:   c.freq,
			Volume:    0.14,
			Attack:    0.08,
			Decay:     0.1,
			Sustain:   0.8,
			Release:   0.1,
			Pan:       0.0,
		})
	}

	// 4. Bassline (16th-note galloping D minor bass)
	var bassNotes []NoteDef
	bassRoots := []float64{FreqD2, FreqD2, FreqF2, FreqG2, FreqBb2, FreqA2, FreqG2, FreqD2}
	for _, root := range bassRoots {
		// 1 bar = 4 beats = 16 sixteenths
		for s := 0; s < 16; s++ {
			f := root
			if s == 6 || s == 14 {
				f = root * 1.5 // fifth
			} else if s == 10 {
				f = root * 2.0 // octave
			}
			bassNotes = append(bassNotes, NoteDef{
				WaveType:  "square",
				DutyCycle: 0.5,
				Duration:  q,
				StartFreq: f,
				EndFreq:   f,
				Volume:    0.26,
				Attack:    0.01,
				Decay:     0.04,
				Sustain:   0.5,
				Release:   0.02,
				Pan:       0.0,
			})
		}
	}

	// 5. Arpeggiator (Rapid 16th-note anime sweeps)
	var arpNotes []NoteDef
	arpScale := []float64{FreqD4, FreqF4, FreqA4, FreqD5, FreqF5, FreqD5, FreqA4, FreqF4}
	totalArpNotes := len(bassRoots) * 16
	for i := 0; i < totalArpNotes; i++ {
		f := arpScale[i%len(arpScale)]
		arpNotes = append(arpNotes, NoteDef{
			WaveType:  "triangle",
			Duration:  q,
			StartFreq: f,
			EndFreq:   f,
			Volume:    0.12,
			Attack:    0.01,
			Decay:     0.03,
			Sustain:   0.4,
			Release:   0.02,
			Pan:       0.4,
		})
	}

	// 6. Drums & Percussion (Kick, Snare, Hi-hats)
	var drumNotes []NoteDef
	for bar := 0; bar < len(bassRoots); bar++ {
		// 1 bar = 4 beats = 8 eighth notes
		for beat := 0; beat < 8; beat++ {
			if beat == 0 || beat == 4 {
				// Kick (pitch sweep from 150Hz down to 40Hz)
				drumNotes = append(drumNotes, NoteDef{
					WaveType:  "sine",
					Duration:  h,
					StartFreq: 150.0,
					EndFreq:   40.0,
					Volume:    0.45,
					Attack:    0.005,
					Decay:     0.08,
					Sustain:   0.1,
					Release:   0.05,
					Pan:       0.0,
				})
			} else if beat == 2 || beat == 6 {
				// Snare (filtered noise burst + mid-tone snap)
				drumNotes = append(drumNotes, NoteDef{
					WaveType:    "noise",
					NoiseFilter: 0.45,
					Duration:    h,
					StartFreq:   240.0,
					EndFreq:     80.0,
					Volume:      0.35,
					Attack:      0.005,
					Decay:       0.09,
					Sustain:     0.05,
					Release:     0.04,
					Pan:         0.0,
				})
			} else {
				// Hi-hat (short crisp noise tick)
				drumNotes = append(drumNotes, NoteDef{
					WaveType:    "noise",
					NoiseFilter: 0.85,
					Duration:    h,
					StartFreq:   8000.0,
					EndFreq:     6000.0,
					Volume:      0.15,
					Attack:      0.002,
					Decay:       0.03,
					Sustain:     0.0,
					Release:     0.02,
					Pan:         0.25,
				})
			}
		}
	}

	tracks := []TrackDef{
		{Name: "Lead", Notes: leadNotes, Pan: -0.2},
		{Name: "Counter", Notes: counterNotes, Pan: 0.3},
		{Name: "Pad", Notes: padNotes, Pan: 0.0},
		{Name: "Bass", Notes: bassNotes, Pan: 0.0},
		{Name: "Arp", Notes: arpNotes, Pan: 0.35},
		{Name: "Drums", Notes: drumNotes, Pan: 0.0},
	}

	return MixTracks(tracks, 0.75)
}

// BuildPanicSoundtrack synthesizes an adrenaline-fueled high tension version of the soundtrack
// featuring rapid staccato pulses, siren warbles, and driving heavy percussion.
func BuildPanicSoundtrack() []byte {
	// Faster tempo: 165 BPM
	b := 60.0 / 165.0
	h := b / 2.0
	q := b / 4.0

	// 1. Alarm Lead (Alternating emergency siren frequencies)
	var alarmNotes []NoteDef
	numBeats := 32
	for i := 0; i < numBeats; i++ {
		f := FreqA4
		if i%2 == 1 {
			f = FreqD5
		}
		alarmNotes = append(alarmNotes, NoteDef{
			WaveType:     "sawtooth",
			Duration:     h,
			StartFreq:    f,
			EndFreq:      f,
			VibratoFreq:  10.0,
			VibratoDepth: 12.0,
			Volume:       0.32,
			Attack:       0.01,
			Decay:        0.04,
			Sustain:      0.8,
			Release:      0.02,
			Pan:          -0.1,
		})
	}

	// 2. Frantic Arpeggios (Diminished chords for pure panic feeling)
	var arpNotes []NoteDef
	dimScale := []float64{FreqD4, FreqF4, FreqBb4, FreqC5, FreqE5, FreqC5, FreqBb4, FreqF4}
	for i := 0; i < numBeats*2; i++ {
		f := dimScale[i%len(dimScale)]
		arpNotes = append(arpNotes, NoteDef{
			WaveType:  "square",
			DutyCycle: 0.2,
			Duration:  q,
			StartFreq: f,
			EndFreq:   f,
			Volume:    0.18,
			Attack:    0.005,
			Decay:     0.02,
			Sustain:   0.3,
			Release:   0.01,
			Pan:       0.3,
		})
	}

	// 3. Heartbeat / Heavy Pulse Bass
	var bassNotes []NoteDef
	for i := 0; i < numBeats; i++ {
		// Double thud like a racing heart
		bassNotes = append(bassNotes, NoteDef{
			WaveType:  "sine",
			Duration:  q,
			StartFreq: 110.0,
			EndFreq:   35.0,
			Volume:    0.45,
			Attack:    0.01,
			Decay:     0.05,
			Sustain:   0.1,
			Release:   0.02,
			Pan:       0.0,
		})
		bassNotes = append(bassNotes, NoteDef{
			WaveType:  "sine",
			Duration:  q,
			StartFreq: 90.0,
			EndFreq:   30.0,
			Volume:    0.35,
			Attack:    0.01,
			Decay:     0.04,
			Sustain:   0.1,
			Release:   0.02,
			Pan:       0.0,
		})
	}

	// 4. Heavy Driving Drums
	var drumNotes []NoteDef
	for i := 0; i < numBeats; i++ {
		if i%2 == 0 {
			drumNotes = append(drumNotes, NoteDef{
				WaveType:  "sine",
				Duration:  h,
				StartFreq: 180.0,
				EndFreq:   40.0,
				Volume:    0.50,
				Attack:    0.005,
				Decay:     0.07,
				Sustain:   0.1,
				Release:   0.03,
				Pan:       0.0,
			})
		} else {
			drumNotes = append(drumNotes, NoteDef{
				WaveType:    "noise",
				NoiseFilter: 0.6,
				Duration:    h,
				StartFreq:   300.0,
				EndFreq:     100.0,
				Volume:      0.40,
				Attack:      0.005,
				Decay:       0.08,
				Sustain:     0.05,
				Release:     0.03,
				Pan:         0.0,
			})
		}
	}

	tracks := []TrackDef{
		{Name: "AlarmLead", Notes: alarmNotes, Pan: -0.1},
		{Name: "PanicArp", Notes: arpNotes, Pan: 0.3},
		{Name: "Heartbeat", Notes: bassNotes, Pan: 0.0},
		{Name: "PanicDrums", Notes: drumNotes, Pan: 0.0},
	}

	return MixTracks(tracks, 0.8)
}

// BuildTitleSoundtrack synthesizes an anime title screen melody (120 BPM, triumphant).
func BuildTitleSoundtrack() []byte {
	b := 60.0 / 120.0
	h := b / 2.0

	var leadNotes []NoteDef
	seq := []struct {
		freq float64
		dur  float64
	}{
		{FreqD4, b}, {FreqF4, b}, {FreqG4, b}, {FreqA4, b * 2},
		{FreqBb4, b}, {FreqA4, b}, {FreqF4, b}, {FreqG4, b * 2},
		{FreqA4, b}, {FreqD5, b * 2}, {FreqC5, b}, {FreqBb4, b},
		{FreqA4, b * 2}, {FreqG4, b}, {FreqF4, b}, {FreqD4, b * 3},
	}
	for _, p := range seq {
		leadNotes = append(leadNotes, NoteDef{
			WaveType:     "sawtooth",
			Duration:     p.dur,
			StartFreq:    p.freq,
			EndFreq:      p.freq,
			VibratoFreq:  6.0,
			VibratoDepth: 3.5,
			Volume:       0.30,
			Attack:       0.03,
			Decay:        0.06,
			Sustain:      0.75,
			Release:      0.08,
			Pan:          0.0,
		})
	}

	var bassNotes []NoteDef
	bassRoots := []float64{FreqD3, FreqF3, FreqG3, FreqA3, FreqBb3, FreqA3, FreqG3, FreqD3}
	for _, r := range bassRoots {
		bassNotes = append(bassNotes, NoteDef{
			WaveType:  "square",
			DutyCycle: 0.4,
			Duration:  h,
			StartFreq: r,
			EndFreq:   r,
			Volume:    0.22,
			Attack:    0.01,
			Decay:     0.05,
			Sustain:   0.6,
			Release:   0.03,
			Pan:       0.0,
		})
		bassNotes = append(bassNotes, NoteDef{
			WaveType:  "square",
			DutyCycle: 0.4,
			Duration:  h,
			StartFreq: r,
			EndFreq:   r,
			Volume:    0.18,
			Attack:    0.01,
			Decay:     0.05,
			Sustain:   0.6,
			Release:   0.03,
			Pan:       0.0,
		})
	}

	tracks := []TrackDef{
		{Name: "TitleLead", Notes: leadNotes, Pan: 0.0},
		{Name: "TitleBass", Notes: bassNotes, Pan: 0.0},
	}

	return MixTracks(tracks, 0.7)
}
