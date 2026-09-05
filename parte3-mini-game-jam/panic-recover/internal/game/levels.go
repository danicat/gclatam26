package game

func GetAllLevels() []LevelDef {
	return []LevelDef{
		// Level 1: "The Whispers" (Small 5x5 room)
		// Teaches basic movement and turn counter.
		{
			Name:              "Level 1: The Whispers",
			SubTitle:          "Reach the Eldritch Artifact before panic sets in.",
			Width:             7,
			Height:            7,
			MaxTurns:          8,
			ClockRecoverTurns: 4,
			PlayerStartX:      1,
			PlayerStartY:      1,
			Layout: []string{
				"#######",
				"#P....#",
				"#.###.#",
				"#.#A#.#",
				"#.#.#.#",
				"#.....#",
				"#######",
			},
		},

		// Level 2: "Bridge the Void" (7x7 room)
		// Teaches boulder pushing into holes to bridge gaps.
		{
			Name:              "Level 2: Bridge the Void",
			SubTitle:          "Push the ancient stone into the abyss to forge a path.",
			Width:             7,
			Height:            7,
			MaxTurns:          10,
			ClockRecoverTurns: 4,
			PlayerStartX:      1,
			PlayerStartY:      2,
			Layout: []string{
				"#######",
				"#.###.#",
				"#P.B.O#",
				"#.#.#.#",
				"#...#A#",
				"#...###",
				"#######",
			},
		},

		// Level 3: "Tick-Tock Panic" (7x7 room)
		// Introduces 80% Chromatic Aberration & Recovery Clock.
		{
			Name:              "Level 3: Tick-Tock Panic",
			SubTitle:          "When panic reaches 80%, reality distorts. Recover your mind!",
			Width:             7,
			Height:            7,
			MaxTurns:          7,
			ClockRecoverTurns: 7,
			PlayerStartX:      5,
			PlayerStartY:      1,
			Layout: []string{
				"#######",
				"#C...P#",
				"#.###.#",
				"#.#.#.#",
				"#.#...#",
				"#.###A#",
				"#######",
			},
		},

		// Level 4: "Crypt of the Elder Gopher" (8x8 room)
		// Multi-boulder + hole puzzle with tight turn budget.
		{
			Name:              "Level 4: Crypt of the Elder",
			SubTitle:          "Two stones, two chasms. Plan your route carefully.",
			Width:             8,
			Height:            8,
			MaxTurns:          16,
			ClockRecoverTurns: 8,
			PlayerStartX:      1,
			PlayerStartY:      1,
			Layout: []string{
				"########",
				"#P.B...#",
				"#.####.#",
				"#...O..#",
				"#.##B#.#",
				"#.#.O#.#",
				"#C...#A#",
				"########",
			},
		},

		// Level 5: "The Eldritch Chamber" (9x9 room)
		// The grand finale: tight panic clock, sanity recovery, and double-bridge.
		{
			Name:              "Level 5: The Eldritch Chamber",
			SubTitle:          "The ultimate relic awaits. Conquer the void, master the panic!",
			Width:             9,
			Height:            9,
			MaxTurns:          20,
			ClockRecoverTurns: 10,
			PlayerStartX:      1,
			PlayerStartY:      7,
			Layout: []string{
				"#########",
				"#A..#..C#",
				"#O#.#.#.#",
				"#.B...B.#",
				"###.###.#",
				"#...#O..#",
				"#.###.###",
				"#P......#",
				"#########",
			},
		},
	}
}
