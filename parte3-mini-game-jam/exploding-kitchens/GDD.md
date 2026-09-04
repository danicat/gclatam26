# Game Design Document: Exploding Kitchens

> **Game Title**: Exploding Kitchens  
> **Event**: GopherCon LATAM 2026 Mini Game Jam  
> **Theme**: Panic and Recover  
> **Genre**: 2D Top-Down Action / Kitchen Defusal Arcade  
> **Target Platform**: Desktop (macOS / Linux / Windows) & WebAssembly (WASM)  
> **Target Resolution**: 16:9 Widescreen (`320x180` virtual pixel canvas, upscaled to `1280x720`)  
> **Author / Lead Designer**: Antigravity & User  

---

## 1. Executive Summary & Elevator Pitch

- **Elevator Pitch**: In *Exploding Kitchens*, you play as a frantic chef desperately sprinting through a single-screen kitchen to defuse overheating appliances and fend off mischievous counter-jumping cats before the entire kitchen detonates.
- **Core Inspiration**: *Exploding Kittens* meets *Overcooked* and *Bomb Squad*.
- **The "Panic & Recover" Dynamic**:
  - **Panic**: Appliances (stoves, fryers, pressure cookers, microwaves) rapidly heat up, emitting smoke, steam, and blaring sirens. Mischievous cats roam the counters, bumping burners to maximum blast.
  - **Recover**: Grab physical emergency tools (Fire Extinguisher, Ice Bucket, Wrench, Cat Toy) and defuse stations in the nick of time. Performing a **"Clutch Recovery"** (defusing within the final 15% of an appliance's fuse) awards an Adrenaline Rush combo multiplier and cools down the kitchen Chaos Gauge.

---

## 2. Core Gameplay Loop & Mechanics

### 2.1 The Gameplay Loop
```text
┌────────────────────────────────────────────────────────┐
│ 1. MONITOR STATIONS                                    │
│    Appliance timers advance; cats roam the room.       │
└──────────────────────────┬─────────────────────────────┘
                           ▼
┌────────────────────────────────────────────────────────┐
│ 2. PANIC TRIGGERS                                      │
│    Appliance enters Warning (Yellow) -> Panic (Red).   │
│    Cat jumps onto counter, speeding up fuse timer!    │
└──────────────────────────┬─────────────────────────────┘
                           ▼
┌────────────────────────────────────────────────────────┐
│ 3. THE SPRINT & TOOL SELECTION                         │
│    Run to Tool Rack: Pick up Extinguisher / Ice /      │
│    Wrench / Cat Toy.                                   │
└──────────────────────────┬─────────────────────────────┘
                           ▼
┌────────────────────────────────────────────────────────┐
│ 4. EMERGENCY ACTION                                    │
│    Apply tool to station with Spacebar before Detonation!│
└──────────┬─────────────────────────────────┬───────────┘
           │ (Saved in time)                 │ (Timer hits 0)
           ▼                                 ▼
┌───────────────────────────────┐ ┌───────────────────────────────┐
│ RECOVER                       │ │ EXPLOSION                     │
│ - Gauge resets to green       │ │ - Station disabled & wrecked  │
│ - Chaos Meter decreases       │ │ - Chaos Meter spikes +25%     │
│ - Clutch Bonus (+500 pts)     │ │ - Camera shake & fire hazard  │
└───────────────────────────────┘ └───────────────────────────────┘
```

### 2.2 Stations, Hazards & Tools Matrix

| Appliance / Entity | Threat / Panic Behavior | Required Tool / Action | Recovery Effect |
| :--- | :--- | :--- | :--- |
| **Pressure Cooker** | Hisses violently; internal PSI rises to critical blowout | **Ice Bucket / Coolant** | Cools pressure to zero; safely vents steam |
| **Deep Fryer** | Boiling oil splatters; catches fire if unattended | **Fire Extinguisher** | Smothers flames; resets frying timer |
| **Microwave Bomb** | Digital countdown ticks rapidly toward detonation | **Manual Vent / Defuse** | Resets cooking timer; dishes cooked safely |
| **Stove Burner** | Pan boils over; spreads grease heat | **Extinguisher or Ice** | Resets burner flame to simmer |
| **Wrecked Station** | Blown up from previous explosion; unusable | **Wrench** | Repairs station back to operational state |
| **Mischievous Cat** | Roams kitchen, leaps on counters, accelerates timers | **Cat Toy / Bell** (or walk up to shoo) | Distracts/shoos cat back to floor rug for 15s |

### 2.3 Win & Loss Conditions
- **Format**: 2-Minute Rush Shift.
- **Victory Condition**: Survive the full 2:00 timer without the Chaos Meter hitting 100%.
- **Loss Condition**: Kitchen Meltdown! Chaos Meter reaches 100% (each explosion adds +25% chaos; continuous unattended fires steadily increase chaos).
- **Scoring**:
  - Regular Defusal: +100 pts
  - Clutch Recovery (final 15% of timer): +500 pts + "CLUTCH!" pop-up
  - Cat Soothed/Shooed: +50 pts
  - Station Repaired: +150 pts
  - Final Grade: 3 Stars based on score, dishes saved, and zero-detonation shifts.

---

## 3. Controls & Input Mapping Scheme

| Logical Action | Keyboard | Gamepad | Touch / Mouse |
| :--- | :--- | :--- | :--- |
| **Move Up / Down / Left / Right** | `W` / `S` / `A` / `D` or Arrows | Left Analog Stick / D-Pad | Virtual Analog D-Pad |
| **Interact / Defuse / Use Tool** | `Space` or `E` | Button South (`A` / `Cross`) | Tap Active Station |
| **Drop Tool** | `Q` or `Shift` | Button West (`X` / `Square`) | Drop Tool Button |
| **Pause Game** | `Escape` or `P` | Start Button | Pause Icon Button |
| **Toggle Fullscreen** | `F11` or `Alt+Enter` | - | Fullscreen Icon |

---

## 4. Visual Style & Asset Strategy

- **Aesthetic**: Pure-Code Procedural Retro Pixel Art (16-bit arcade aesthetic).
- **Zero-Dependency Asset Guarantee**: 100% generated in Go code via `procedural-art`. No external PNGs or image files required.
- **Entities & Visual Cues**:
  - **Chef**: Animated pixel character with 4-directional walking frames, holding tool overhead.
  - **Appliances**: Distinct silhouettes with dynamic status LEDs and floating warning meters.
  - **Mischievous Cat**: Animated calico/orange cat with walking, counter-jumping, and startled shoo animations.
  - **VFX Particle Pools**:
    - Steam: Rising translucent white puffs.
    - Fire & Embers: Flickering yellow/orange/red rising particles.
    - Smoke: Billowing dark gray clouds during Panic states.
    - Explosion Shockwave: Expanding circular blast ring, debris particles, and camera screen shake.
  - **Emergency Strobe**: Pulsing vignette border when Chaos Meter exceeds 75%.

---

## 5. Audio & Soundscape Strategy

- **Audio Engine**: Pure-Code DSP Chiptune Audio (`procedural-composer`). Zero external WAV/MP3 files; 16-bit 44.1kHz stereo PCM generated in memory.
- **Dynamic Background Music (BGM)**:
  - Base Track: Upbeat, groovy kitchen arcade chiptune (120 BPM).
  - Panic Modulation: Dynamically accelerates to 150+ BPM with a higher pitch and added siren arpeggios when 2+ stations enter Panic state.
- **Sound Effects (SFX)**:
  - `sfx_warning_beep`: High-pitched intermittent beeping for stations entering Warning state.
  - `sfx_alarm_klaxon`: Urgent two-tone siren for stations in critical Panic state.
  - `sfx_steam_hiss`: White noise burst with low-pass filter for venting valves.
  - `sfx_extinguisher`: Continuous gentle hiss for fire smothering.
  - `sfx_defuse_clutch`: Bright ascending chime for successful recovery.
  - `sfx_explosion`: Deep synthesized noise burst with low pitch drop and screen shake.
  - `sfx_cat_meow`: Cute synthesized chirp/meow when cat is shooed or leaps onto counter.

---

## 6. Scene Sequence & HUD Layout

### 6.1 State Flow
```text
Boot -> Title Screen (with 10s Attract Demo Mode) -> Gameplay (2-min Shift) -> Pause -> Game Over Debrief
```

- **Attract Demo Mode**: If idle on the Title Screen for 10 seconds, an autonomous AI chef runs around the kitchen managing stations until user presses any key.
- **Game Over Debrief**: Shows Shift Rating (1–3 Stars), Total Score, Clutches Pulled, Detonations, and Local High Score.

### 6.2 HUD Layout (Virtual $320 \times 180$)
- **Top-Left**: Score (`SCORE: 004500`) and Clutch Badge (`CLUTCH x3`).
- **Top-Center**: Shift Timer (`01:45`).
- **Top-Right**: Chaos Meter (`[██████░░░░] 60%`).
- **Bottom-Center**: Currently Equipped Tool Icon (`[EXTINGUISHER]`).

---

## 7. Technical Scope & Architecture

Following [`ebitengineer`](file:///Users/dev/source/opensource/joubertredrat/gclatam26/parte3-mini-game-jam/.agents/skills/ebitengineer/SKILL.md) standards:
- **Project Location**: `parte3-mini-game-jam/exploding-kitchens/`
- **Go Version**: Go 1.26+
- **Virtual Resolution**: Constant $320 \times 180$ internal canvas scaled smoothly to window/fullscreen.
- **Draw Loop**: Zero memory allocations inside `Draw(screen *ebiten.Image)`.
- **Target Deployment**: Single cross-platform Go binary + WebAssembly build (`GOOS=js GOARCH=wasm`).
