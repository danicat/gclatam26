# Game Project Plan: Exploding Kitchens

> **Event**: GopherCon LATAM 2026 Mini Game Jam  
> **Theme**: Panic and Recover  
> **Inspiration**: *Exploding Kittens* × *Overcooked* / Kitchen Meltdown  
> **Target Framework**: Go 1.26+ with [Ebitengine v2](https://ebitengine.org/) (`ebitengineer`)  
> **Platform**: Desktop (macOS/Linux/Windows) + WebAssembly (WASM)  
> **Canvas**: 16:9 Virtual Resolution (`320x180` retro pixel canvas scaled to `1280x720`)  

---

## 1. Executive Summary & Pitch

In **Exploding Kitchens**, players take on the role of a frantic chef managing a high-stress kitchen where appliances and mischievous cats conspire to trigger catastrophic detonations. 

### Core Theme: Panic & Recover
- **The Panic**: Cooking stations tick up in heat and pressure. Unattended stations transition into critical overload with flashing sirens, smoke particles, and violent camera shakes. Mischievous cats roam the countertops, bumping knobs to maximum blast.
- **The Recover**: Players must sprint across the kitchen, deploy recovery tools (fire extinguishers, ice buckets, venting wrenches), and defuse stations in the nick of time. Pulling off a **"Clutch Recovery"** (defusing within the final 15% of the panic fuse) triggers an Adrenaline Rush combo multiplier and cools the chaos gauge.

---

## 2. Gameplay Loop & Mechanics

### 2.1 The Core Cycle
```text
[Prep / Cooking Stage]
        │
        ▼ (Time accrues / Cat interference)
[WARNING: Beeps + Steam]
        │
        ▼ (Ignored)
[PANIC: Red Strobe + Klaxon Siren + Violent Smoke]
   ├─── Action: Defuse with Extinguisher/Ice ──► [RECOVERY: Score Boost + Adrenaline Combo]
   └─── Timer Reaches 0 ──────────────────────► [EXPLOSION: Station Disabled + Hazard Spread]
                                                        │
                                                        ▼ (Wrench repair)
                                                 [STATION RESTORED]
```

### 2.2 Kitchen Stations & Elements
1. **Cooking Stations (4–6 Interactive Appliances)**:
   - **Pressure Cooker**: Builds pressure quickly; requires valve venting before detonation.
   - **Deep Fryer**: Bubbles with boiling oil; catches fire if left unattended; extinguished with wet towels/extinguisher.
   - **Microwave Bomb**: Countdown timer with fast ticking; requires defuse interaction.
   - **Stove Top**: Basic burners prone to boiling over; cooled with ice packs.
2. **Tools & Pickups**:
   - **Fire Extinguisher**: Smothers active grease fires and cools critical stations.
   - **Ice Pack / Water Bucket**: Resets overheating temperature gauges to zero.
   - **Wrench**: Repairs exploded or wrecked stations back to operational status.
3. **The Chaos Agent ("The Cat")**:
   - An unpredictable kitten roaming the kitchen counter.
   - Randomly activates burners, speeds up appliance fuses, or knocks over items unless shooed away or distracted.

### 2.3 Win & Loss Conditions
- **Format**: 2-minute arcade rush / survival mode.
- **Loss Condition**: Total Kitchen Meltdown (Chaos Gauge reaches 100% when 3+ stations detonate simultaneously).
- **Victory Condition**: Survive the full 2-minute shift with the kitchen standing, maximizing reputation score and clutch recoveries.

---

## 3. Controls & Input Mapping

Standard action mapping decoupled from specific keys via an input abstraction layer:

| Logical Action | Keyboard | Gamepad | Touch / Mouse |
| :--- | :--- | :--- | :--- |
| **Move Up / Down / Left / Right** | `W`/`S`/`A`/`D` or Arrow Keys | Left Stick / D-Pad | Virtual Joystick / Tap-to-move |
| **Interact / Defuse / Vent** | `Space` or `E` | Button South (`A` / `Cross`) | Tap Station |
| **Drop / Swap Tool** | `Q` or `Shift` | Button West (`X` / `Square`) | Tool Icon Tap |
| **Pause Game** | `Escape` or `P` | Start / Options | Pause Button |
| **Toggle Fullscreen** | `F11` or `Alt+Enter` | - | Fullscreen Icon |

---

## 4. Visual & Audio Pipeline

To maintain 100% reliability, zero external file bugs, and instant WebAssembly compilation during a 90-minute game jam:

- **Graphics (`procedural-art`)**:
  - Pure Go code rendering vector shapes and retro pixel sprites.
  - Pre-allocated zero-allocation particle systems for steam, smoke, embers, fire, and explosion shockwaves.
  - Screen shake and emergency red border vignette when the kitchen chaos gauge exceeds 75%.
- **Sound & Music (`procedural-composer`)**:
  - Zero-dependency procedural 16-bit PCM sound synthesis generated in Go.
  - Distinct SFX: Station tick, warning beeps, steam hiss, alarm klaxon, explosion boom, defuse chime, and fanfare.
  - Fast-tempo chiptune BGM that dynamically accelerates its BPM when Panic states are triggered.

---

## 5. Engine Architecture (`ebitengineer` Guidelines)

Organized under `parte3-mini-game-jam/exploding-kitchens/` adhering to modular Go standards:

```text
parte3-mini-game-jam/exploding-kitchens/
├── cmd/
│   └── game/
│       └── main.go              # Entry point: SetTPS(60), SetWindowSize, ebiten.RunGame
├── internal/
│   ├── game/
│   │   ├── game.go              # Master ebiten.Game implementation & Virtual Resolution Layout
│   │   └── state.go             # Scene State Machine (Boot -> Title -> Gameplay -> GameOver)
│   ├── scene/
│   │   ├── title.go             # Title screen with instructions & Attract Demo mode
│   │   ├── play.go              # Core gameplay loop (Stations, Chaos meter, Defuse actions)
│   │   └── gameover.go          # End-of-shift report, Clutch counter, High Score
│   ├── entity/
│   │   ├── chef.go              # Player entity, movement delta physics, collision bounding
│   │   ├── station.go           # Appliance FSM (Idle -> Cooking -> Warning -> Panic -> Exploded)
│   │   ├── tool.go              # Tools (Extinguisher, Ice, Wrench)
│   │   └── cat.go               # Autonomous roaming chaos entity
│   ├── system/
│   │   ├── input.go             # Rebindable input manager (Keyboard / Gamepad)
│   │   ├── particle.go          # Pre-allocated zero-alloc smoke/fire/steam particle pool
│   │   └── audio.go             # Procedural sound synthesis & dynamic BGM
│   └── ui/
│       ├── hud.go               # Chaos meter, Score, Timer, Station warning icons
│       └── font.go              # Embedded pixel typography
├── go.mod
└── README.md
```

### Critical Performance Guardrails
1. **Zero Allocations in `Draw()`**: No `ebiten.NewImage()`, slice expansions, or format strings inside draw loops.
2. **Fixed Virtual Resolution Canvas**: Internal fixed canvas of $320 \times 180$ automatically scaled to fit any window without blur.
3. **Decoupled Delta-Time Logic**: Updates run at 60 TPS with delta-time integration to ensure identical gameplay on high-refresh monitors.

---

## 6. 90-Minute Game Jam Phased Milestones

```text
[00:00 - 00:20] Milestone 1: Engine Skeleton & Player
  ├── Initialize Go module & Ebitengine v2 dependencies
  ├── Virtual canvas (320x180) and window setup
  ├── Scene State Machine (Title -> Play -> GameOver)
  └── Chef entity with top-down WASD movement & kitchen boundary collisions

[00:20 - 00:45] Milestone 2: Cooking Stations & The Panic State Machine
  ├── Deploy interactive stations (Stove, Fryer, Pressure Cooker, Microwave)
  ├── Implement appliance tick timers (Idle -> Warning -> Panic -> Explode)
  ├── Tool holding & pickup system (Extinguisher, Ice, Wrench)
  └── Visual warning indicators (gauge overlays, color strobing)

[00:45 - 01:10] Milestone 3: Panic & Recovery Mechanics
  ├── Defuse / Vent / Extinguish interaction logic
  ├── "Clutch Recovery" window (bonus score & chaos relief for last-second saves)
  ├── Roaming Cat chaos agent
  └── Chaos gauge & 2-minute countdown timer with win/loss conditions

[01:10 - 01:25] Milestone 4: Audio, VFX Particles & Polish
  ├── Pre-allocated particle system: Steam, fire, smoke, and explosion bursts
  ├── Procedural audio generator: Sirens, defuse chime, explosions, background loop
  └── Screen shake and HUD polish

[01:25 - 01:30] Milestone 5: Verification & Jam Delivery
  ├── Automated testing (`go test ./...`)
  ├── Native build verification (`go run ./cmd/game`)
  ├── WebAssembly build verification (`GOOS=js GOARCH=wasm`)
  └── Documentation: Controls, lore, and execution guide in README.md
```
