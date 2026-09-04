# Panic! At The Disco: Saturday Night Flee-ver 🪩

> **"Stayin' Alive is no longer just a song — it's an evacuation plan!"**

**Panic! At The Disco: Saturday Night Flee-ver** is a 2D top-down retro arcade rhythm evader built in **Go** using [Ebitengine v2](https://ebitengine.org/). 

Featuring **100% pure-code procedural graphics**, custom retro bitmap typography, and a real-time **44.1 kHz PCM DSP disco synthesizer**, the game runs completely standalone with **zero external image or audio files**.

---

## 📖 The Premise

It is Saturday night, 1978. The discotheque's subwoofers have overloaded the venue's vintage architecture, and the ceiling is collapsing in 4/4 disco time! 

You play as **Tony**, the club's grooviest patron. As giant mirrored disco balls plummet, steel trusses snap, and panicked clubbers stumble across the flashing dance floor, you must race against the clock to navigate 3 escalating club zones and reach the emergency fire exits before the entire venue caves in.

---

## ✨ Features

- **100% Zero-Asset Procedural Graphics**: Glowing retro neon vector art, animated afro dancer with dynamic panic jitters, 3D-angled facet disco balls, overhead trusses with hazard stripes, and spilled cocktail puddles rendered in real-time.
- **Custom 5x7 Retro Pixel Typography**: Built-in procedural bitmap font engine for clean arcade HUDs, timer alerts, and menus.
- **Real-Time DSP Audio Synthesizer**: 16-bit 44.1 kHz stereo audio engine generating a funky 120-BPM disco chiptune groove (*"Stayin' Uncrushed"*) with slap bass, kick, open hi-hats, brass stabs, and procedural SFX (whistles, crashes, dashes, slips, and sirens).
- **Dynamic Musical Acceleration**: When the collapse timer dips below 12 seconds, the music dynamically transitions to a 144-BPM fast panic tempo accompanied by emergency sirens.
- **Beat-Panic Rhythm Mechanics**: Hazards drop and floor tiles surge on the musical beat. Evading cleanly charges your **Groove Meter**, enabling an invulnerable **Disco Dash** sprint.
- **3-Stage Escape Gauntlet**:
  1. **Zone 1: The Main Dance Floor** (45s collapse timer, falling disco balls, electrified strobe tiles, wandering clubbers).
  2. **Zone 2: The VIP Lounge & Bar** (40s collapse timer, slippery cocktail puddles with drift physics, tighter pathways).
  3. **Zone 3: Backstage & Fire Alley** (35s collapse timer, falling steel lighting trusses, sparking floor hazards, final getaway).

---

## 🚀 How to Run the Game

### Prerequisites
- [Go](https://go.dev/dl/) (version **1.22** or newer recommended).
- A working display and audio output (macOS, Windows, or Linux).

### 1. Run Directly with `go run`
Navigate to the project root directory and execute:
```bash
cd panic-at-the-disco
go run .
```

### 2. Build and Run Standalone Binary
To compile an optimized, self-contained binary:
```bash
# Compile native binary
go build -o panic-at-the-disco .

# Run the executable:
./panic-at-the-disco
```

---

## 🕹️ Controls

| Action | Keyboard | Gamepad |
| :--- | :--- | :--- |
| **Move Dancer** | `W, A, S, D` or `Arrow Keys` | Left Analog Stick / D-Pad |
| **Disco Dash** (Invulnerable Sprint) | `Space` or `J` | Button South / Right Trigger |
| **Confirm / Start / Continue** | `Space` or `Enter` | Button South / Start |
| **Quick Restart** | `R` | Button Select / Back |
| **Pause / Return to Title** | `Escape` | Start / Menu |
| **Toggle Fullscreen** | `F11` or `Alt + Enter` | — |

---

## 🕺 How to Play & Survival Guide

### Core Objective
Reach the illuminated green **EMERGENCY EXIT** before the **COLLAPSE TIMER** reaches zero.

### Meters & Mechanics
- **Lives (`<3 <3 <3`)**: You start with 3 lives. Taking a hit from a falling disco ball, steel truss, or electrical tile cost 1 life and grants 1.2s of invulnerability.
- **Panic Meter (`0% - 100%`)**: As the ceiling rumbles, your panic rises. High panic causes Tony's eyes to widen, arms to flail, and sweat drops to burst.
- **Groove Meter & Disco Dash**: Moving cleanly builds your Groove Meter. When it reaches 30%, press `[SPACE]` to trigger a **Disco Dash**—a high-speed sprint leaving neon ghost afterimages that ignores hazards and reduces your panic!
- **Spilled Cocktails (Drink Puddles)**: Walking into spilled drinks (Amber Whiskey, Red Cosmopolitan, Green Midori) triggers low-friction drift physics. Use dashes to correct your slide.
- **Telegraph Warnings**: Watch the floor! Expanding red circular shadows telegraph incoming disco balls, and amber flashing tiles warn of impending electrical floor bursts.

---

## 🧪 Running Automated Tests

The codebase includes comprehensive unit tests covering player physics, hazard lifecycles, level generation, particle pooling, and procedural DSP waveform synthesis.

To run the full test suite:
```bash
go test -v ./...
```

---

## 📁 Project Architecture

```text
panic-at-the-disco/
├── main.go                     # Application entrypoint & window configuration (1280x720, 640x360 virtual)
├── GDD.md                      # Complete Game Design Document
├── README.md                   # You are here!
├── go.mod / go.sum             # Ebitengine v2 dependencies
└── internal/
    ├── audio/
    │   ├── synth.go            # Procedural 44.1 kHz PCM synthesizer & SFX generators
    │   └── synth_test.go       # Audio buffer & waveform unit tests
    ├── entities/
    │   ├── player.go           # Player movement, friction drift, panic/groove meters, dash logic
    │   ├── player_test.go      # Player physics & invulnerability unit tests
    │   ├── hazards.go          # Falling disco balls, trusses, drink puddles, clubbers, exit portal
    │   └── hazards_test.go     # Hazard states & collision unit tests
    ├── gfx/
    │   ├── disco_floor.go      # Animated 70s dance floor grid & hazard tile management
    │   ├── sprites.go          # Pure-code procedural vector renderers
    │   ├── font.go             # 5x7 retro bitmap pixel font renderer
    │   ├── particles.go        # Pre-allocated zero-allocation particle pool
    │   └── particles_test.go   # Particle pool management unit tests
    ├── levels/
    │   ├── level.go            # 3-Zone configurations and stage entity generator
    │   └── level_test.go       # Zone layout & parameter unit tests
    ├── input/
    │   └── input.go            # Unified keyboard and gamepad input polling
    ├── scenes/
    │   ├── scene.go            # Scene FSM interface and action definitions
    │   ├── title_scene.go      # Attract title screen & control instructions
    │   ├── play_scene.go       # Core gameplay loop, camera shake, spawner, and HUD
    │   ├── clear_scene.go      # Zone evacuation interstitial & score tally
    │   ├── gameover_scene.go   # Death screen with retry option
    │   ├── victory_scene.go    # Sunrise getaway celebration with confetti
    │   └── scenes_test.go      # Scene transition & lifecycle tests
    └── game/
        └── game.go             # Ebitengine Game implementation (Update, Draw, Layout)
```

---

## 🏆 Game Jam & Tech Stack

- **Engine**: [Ebitengine v2](https://github.com/hajimehoshi/ebiten) (v2.9.11)
- **Language**: Go 1.22+
- **Platform**: Cross-platform desktop (macOS, Linux, Windows) & WebAssembly compatible.
- **Audio & Visuals**: 100% Procedurally synthesized in Go without external audio or image files.
