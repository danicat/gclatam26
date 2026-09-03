# Game Design Document (GDD)

> **Game Title**: Panic! At The Disco: Saturday Night Flee-ver  
> **Target Genre**: Top-Down 2D Action Arcade / Rhythm Evader  
> **Target Platform**: Desktop (macOS / Linux / Windows) & WebAssembly (Browser / Cloud Run)  
> **Target Aspect Ratio & Resolution**: 16:9 Widescreen (`640x360` Virtual Resolution, integer scaled to window)  
> **Author / Lead Designer**: Antigravity Game Designer & Matheus Arthus  

---

## 1. Executive Summary & Elevator Pitch

- **Elevator Pitch**: *A high-voltage top-down retro arcade escape game where an anxious disco dancer must flee a crumbling 1970s nightclub whose ceiling is collapsing in sync with a 120-BPM four-on-the-floor beat.*
- **Core Inspiration**: *Pac-Man* evasion + *Smash TV* arena energy + *Crypt of the NecroDancer* rhythmic anticipation + *Canabalt* urgency.
- **Target Audience & Mood**: Fast-paced, comedic, panic-inducing arcade action with dazzling neon disco aesthetics and funky groovy tension.

---

## 2. Core Gameplay Loop & Mechanics

### 2.1 Primary Gameplay Loop
```text
  [Hear Beat & Notice Ceiling Shadow / Flashing Tile]
                     │
                     ▼
  [Dodge / Move in Rhythm Toward Exit Door]
                     │
                     ▼
  [Avoid Panicked Clubbers & Slippery Spills]
                     │
                     ▼
  [Charge Panic Meter & Trigger 'Disco Dash']
                     │
                     ▼
  [Survive the Zone Countdown & Reach the Emergency Exit]
```

1. **Observe & Predict**: The 120-BPM disco track drives the club's environment. Floor tiles pulse neon colors, while telegraph shadows and laser warnings signal where chandeliers, mirror balls, and trusses will crash down on the next beat.
2. **Evade & Maneuver**: The player steers their flailing, wide-eyed disco dancer across dynamic obstacles (spilled cocktail slicks, panicked wandering dancers, collapsing velvet ropes).
3. **Beat-Panic Mechanic**:
   - **Panic Meter (0–100%)**: Starts at 20%. Rises when near-misses occur, staying in hazard zones, or when the zone timer dwindles.
   - **High Panic Effects**: Heartbeat thuds, camera jitter, slight control slip.
   - **Calming & Charging**: Moving cleanly or picking up discarded glowsticks/cocktails reduces panic and builds the **Groove Bar**.
   - **Disco Dash (Action)**: When Groove is charged (or Panic is low), pressing `Space` triggers an invulnerable, glittery slipstream dash forward, leaving neon afterimages and smashing through light debris.
4. **Zone Completion**: Reach the illuminated `EXIT` door before the Zone Collapse Timer reaches 0:00 and the roof caves in completely.

### 2.2 Win & Loss Conditions
- **Victory Condition**: Successfully escape all 3 progressive club zones (The Dance Floor $\rightarrow$ The VIP Lounge $\rightarrow$ The Backstage Alley) and leap into the getaway cab.
- **Defeat Condition**:
  - Crushed by falling ceiling debris (disco ball, speaker truss, concrete chunks).
  - Zone Collapse Timer hits `00:00` (catastrophic roof collapse).
  - Health/Stamina depleted by hazards and electrified cables.

### 2.3 Stage Breakdown (3-Stage Gauntlet)
| Stage | Name | Key Hazards & Obstacles | Aesthetic Palette |
| :--- | :--- | :--- | :--- |
| **Zone 1** | **The Main Dance Floor** | Giant falling disco balls (creates radial shockwave), rhythm floor tiles that turn red/burn, spinning spotlight dazzlers. | Hot Magenta, Electric Cyan, Gold |
| **Zone 2** | **The VIP Lounge & Bar** | Collapsing champagne towers, slippery alcohol puddles (loss of traction), toppling velvet stanchions, stampeding VIPs. | Deep Purple, Velvet Red, Amber |
| **Zone 3** | **Backstage & Fire Alley** | Falling lighting trusses, sparking electric cables (rhythmic arcs), steam pipes, falling heavy roadie amp cases. | Neon Green, Industrial Gray, HazMat Yellow |

---

## 3. Controls & Input Mapping Scheme

- **Primary Input Devices**: Keyboard (WASD / Arrows) & Gamepad (DualShock / Xbox standard)

| Logical Action | Keyboard / Mouse | Gamepad Button | Description |
| :--- | :--- | :--- | :--- |
| `ActionMoveUp` | `W` / `Up Arrow` | D-Pad Up / Left Stick Up | Move upward |
| `ActionMoveDown` | `S` / `Down Arrow` | D-Pad Down / Left Stick Down | Move downward |
| `ActionMoveLeft` | `A` / `Left Arrow` | D-Pad Left / Left Stick Left | Move left |
| `ActionMoveRight` | `D` / `Right Arrow` | D-Pad Right / Left Stick Right | Move right |
| `ActionDiscoDash` | `Space` / `J` | Button South (`A` / `Cross`) | Invulnerable sprint burst (on cooldown/groove) |
| `ActionPause` | `Escape` / `P` | Start / Options | Pause game & view menu |
| `ActionRestart` | `R` | Button West (`X` / `Square`) | Quick restart after Game Over |

---

## 4. Visual Style & Asset Strategy

- **Aesthetic Direction**: **100% Pure-Code Procedural Retro Arcade (Zero External Asset Dependencies)**.
  - Crisp, vibrant 1970s synth/disco palette rendered via Go procedural vector primitives and color blending.
  - Multi-layered procedural lighting: illuminated 8x8 disco floor tiles, pulsing neon spotlights, dynamic radial glow around the player, shadow telegraph circles for falling debris.
- **Pure-Code Procedural Implementation (`procedural-art`)**:
  - **Player Character**: Procedurally drawn retro character with animated afro, bell-bottom pants that swing with 8-direction running gait, and flailing panic arm motion.
  - **Environment**: Procedural tile grid with cycling RGB color palettes, reflective mirror ball sprite synthesized with mathematical ray shading and sparkle bursts.
  - **Particle Systems**:
    - *Mirror Shards*: Silver/cyan glittering debris that burst and bounce upon disco ball impact.
    - *Spark Emitters*: Yellow/white sparks crackling from fallen trusses and severed wires.
    - *Sweat Drops*: Comedic panic water drop particles spurting from the player when sprinting or in danger.
    - *Disco Afterimages*: Fading translucent neon silhouettes left behind during the Disco Dash.

---

## 5. Audio & Soundscape Strategy

- **Audio Pipeline**: **100% Pure-Code DSP Synthesis (`procedural-composer`)**.
  - All music and sound effects synthesized directly in memory at runtime into 16-bit 44.1 kHz PCM stereo streams—no external MP3/WAV files required.
- **Background Music (BGM)**:
  - **Track**: *"Stayin' Uncrushed"* (Procedural 120-BPM 70s Disco Funk).
  - **Instrumentation**:
    - Four-on-the-floor synthesized punchy kick drum (`55 Hz` pitch drop).
    - Off-beat 16th-note synthesized hi-hats (`White noise` high-pass bursts).
    - Funky slapping sawtooth bassline running active minor-pentatonic arpeggios.
    - 70s synth brass stabs on beats 2 and 4.
  - **Dynamic Audio Modulation**: As the timer drops below 15 seconds, the tempo increases to 135 BPM and an emergency siren oscillator joins the mix.
- **Sound Effects (SFX)**:
  - `SFX_DISCO_FALL`: Downward pitch whistle FM synth followed by explosive bass impact and glass shatter.
  - `SFX_DASH`: Upward disco glissando sweep with glitter sparkle noise.
  - `SFX_STEP_BEAT`: Subtle sync click rewarding moves timed exactly to the beat.
  - `SFX_PANIC_GASP`: Comedic vocaloid-style high pitch squeak.
  - `SFX_EXIT_CHEER`: Ascending brass fanfare when crossing the emergency threshold.

---

## 6. Game State Sequence & HUD Layout

### 6.1 State Flow Diagram
```text
  ┌──────────────┐
  │  Boot Scene  │
  └──────┬───────┘
         ▼
  ┌──────────────┐       [Game Start]       ┌──────────────┐
  │ Title Screen ├─────────────────────────►│  Game Scene  │◄────────────┐
  └──────────────┘                          └──────┬───────┘             │
                                                   │                     │
                    ┌──────────────────────────────┴──────────────┐      │
                    ▼                                             ▼      │
            [Time/Crush Loss]                              [Door Reached]│
                    ▼                                             ▼      │
            ┌──────────────┐                              ┌──────────────┴──────────┐
            │  Game Over   │                              │ Stage Complete / Inter. │
            └──────┬───────┘                              └──────────────┬──────────┘
                   │                                                     │
                   │ [Press R]                                  [Zone 3 Complete]
                   └─────────────────────────────────────────────────────▼
                                                          ┌──────────────┐
                                                          │ Victory Win  │
                                                          └──────────────┘
```

### 6.2 HUD & UI Overlay Layout (640x360 Canvas)
- **Top-Left**:
  - **Panic Gauge**: Heart icon with a retro bar that fills with vibrating color (Green $\rightarrow$ Yellow $\rightarrow$ Crimson Flash).
  - **Lives / Stamina**: 3 Disco Heart icons.
- **Top-Center**:
  - **Ceiling Collapse Countdown Timer**: Retro digital LED font displaying `00:45.00` in glowing amber/red.
  - **Stage Indicator**: `ZONE 1: THE DANCE FLOOR`.
- **Top-Right**:
  - **Score**: Points accumulated from time survived, near-misses, and on-beat steps.
  - **Groove Meter**: Flashing stars indicating Disco Dash availability (`READY [SPACE]`).
- **Screen Borders**:
  - Subtle flashing neon vignette pulsing to the 120-BPM tempo; turns red and shakes when debris is about to hit.

---

## 7. Technical Scope & Architecture Notes

### 7.1 Ebitengine Subsystem Plan (Go 1.26+)
- **Package Structure**:
  ```text
  panic-at-the-disco/
  ├── main.go               # Window configuration, entrypoint, virtual canvas scale
  ├── go.mod                # Module definitions and Ebitengine v2 dependencies
  ├── GDD.md                # Authoritative Game Design Document
  └── internal/
      ├── game/             # Core Game struct, scene manager, state dispatch
      ├── entities/         # Player, falling debris, obstacles, NPCs
      ├── scenes/           # TitleScene, PlayScene, GameOverScene, WinScene
      ├── audio/            # Procedural DSP disco engine & SFX synth (procedural-composer)
      ├── gfx/              # Procedural vector drawing, disco floor shader/grid, particles
      └── input/            # Unified keyboard and gamepad mapping
  ```
- **Physics & Collision Detection**:
  - Axis-Aligned Bounding Boxes (AABB) for player-to-obstacle and player-to-exit collisions.
  - Radial circles for falling disco ball impact zones with telegraph radius expansions.
  - Low-friction velocity physics for slippery alcohol spills.
- **Zero-Allocation Render Loops**:
  - Reusable slice buffers for particles and debris entities.
  - Pre-allocated vertex buffers and off-screen render targets for disco floor tile matrix.
  - Strict separation of `Update()` (fixed 60 Hz tick) and `Draw()` (interpolation).

---

## 8. Milestone Implementation Roadmap

- **Phase 1: Foundation & Core Loop Prototype**
  - Initialize Ebitengine v2 project structure, virtual resolution (640x360), input handling, and 8-direction player movement with basic AABB collisions.
- **Phase 2: Collapsing Ceiling & Hazard System**
  - Implement rhythmic telegraphing indicators, falling disco ball and debris spawner, collision damage, and collapse timer.
- **Phase 3: Procedural Visuals & Disco Aesthetics**
  - Implement animated neon dance floor grid, procedural disco character animation (afro, bell-bottoms, panic expressions), dynamic shadows, and particle sparks.
- **Phase 4: Procedural Audio & Beat-Sync Engine**
  - Implement procedural 120-BPM disco music track and SFX synthesizer in pure Go using `oto` / Ebitengine audio player.
- **Phase 5: Multi-Stage Gauntlet & Polish**
  - Build Zone 1 (Dance Floor), Zone 2 (VIP Lounge), Zone 3 (Backstage Alley).
  - Add screen shake, Panic Meter mechanics, Disco Dash afterimages, HUD, Title, Game Over, and Victory screens.
