# Game Design Document (GDD): Panic Recover - Runtime Defender

> **Game Title**: Panic Recover: Runtime Defender  
> **Theme**: Panic (or recover?)  
> **Genre**: Top-Down Endless Scrolling Space Shooter / Shmup  
> **Target Platform**: Desktop (Linux / macOS / Windows) & WebAssembly  
> **Target Aspect Ratio & Resolution**: 16:9 Widescreen (`640x360` Virtual Pixel Canvas)  
> **Lead Designer & Developers**: Cesar & Antigravity  

---

## 1. Executive Summary & Elevator Pitch
- **Elevator Pitch**: Pilot the Gopher Vanguard through the turbulent Go runtime, blasting incoming runtime errors and concurrency hazards in an endless scrolling vertical battlefield. When fatal damage strikes, enter high-stakes **PANIC MODE** (a 5-second adrenaline frenzy with invulnerability and overclocked fire); eliminate elite bugs or grab a `recover()` token before the call stack unwinds to resurrect your ship!
- **Core Inspiration**: Classic arcade shmups (*Galaga*, *Ikaruga*) infused with authentic Go runtime mechanics and debugging humor.
- **Mood**: High-octane cyberpunk arcade with retro neon vector visuals, pulse-pounding chiptune beats, and chaotic runtime panics.

---

## 2. Core Gameplay Loop & Mechanics
```
[Fly & Dodge] ──> [Shoot Bugs / Avoid Fire] ──> [Near-Death: Fatal Damage]
      ▲                                                    │
      │                                             [PANIC MODE! (5s)]
      │                                            (Overclocked Fire)
      │                                                    │
[Survive & Escalate] <── [Catch recover() Drop] <───────────┤
                                                           ▼ (Time Expires)
                                                  [Game Over: Stack Trace]
```

### Primary Mechanics
1. **Vertical Endless Scrolling**:
   - Parallax cyber-grid starfield streaming downwards with floating hex memory addresses and code fragments.
2. **Player Ship (Gopher Vanguard)**:
   - Agile 8-way movement with retro engine trail particles.
   - Primary weapon: Twin Laser Beams (or Spread Cannon with upgrades).
3. **The Panic / Recover Hook (Near-Death Surge)**:
   - Taking a fatal blow does not immediately end the game: it triggers **PANIC MODE**.
   - During Panic Mode:
     - Warning sirens blare, screen pulses red with scanline distortion.
     - Ship is invulnerable for 5.0 seconds.
     - Fire rate increases by 300% (erratic hyper-burst).
     - A ticking Call Stack countdown appears on screen (`5.0s ... 0.0s`).
     - Defeating an enemy during Panic drops a **`recover()` token**.
     - Catching `recover()` clears the panic state, stabilizes health, and unleashes a shockwave destroying enemy projectiles!
     - If the timer hits `0.0s` without recovery, the call stack unwinds into a full crash (**Game Over** with formatted stack trace).

---

## 3. Enemy Roster (Classic Go Problems)

| Enemy Bug | Visual Silhouette | Behavior Pattern | Attack / Hazard |
| :--- | :--- | :--- | :--- |
| **`nil pointer`** | Ultra-sharp purple/violet dart | Zippy zigzag dashes; evades straight fire | Swift straight plasma bolts |
| **`concurrent map writes`** | Dual red linked ships | Fly in pairs, mirroring movements | Cross-angled laser streams |
| **`deadlock`** | Armored bronze square fortress | Moves slowly down the screen, high HP | Mutex shield barrier that deflects shots |
| **`memory leak`** | Pulsing green amoeba/blob | Slowly expands in size until shot | Splits into 2 smaller leaks when destroyed |
| **`goroutine leak`** | Swarm of tiny orange nano-bots | Fast dive-bomb formations | High volume, low individual HP |

---

## 4. Power-ups & Pickups

- **`recover()`**: Restores ship from Panic mode to normal state, or acts as an EMP bomb clearing bullets if collected while healthy.
- **`sync.Mutex`**: Grants a 6-second frontal energy shield absorbing incoming damage.
- **`go worker`**: Spawns an autonomous companion drone that shoots alongside the player.

---

## 5. Controls & Input Mapping Scheme

| Logical Action | Keyboard / Mouse | Gamepad Button |
| :--- | :--- | :--- |
| **Movement** | `W`/`A`/`S`/`D` or Arrow Keys | Left Stick / D-Pad |
| **Fire** | `Space` / `J` / Left Click (Hold for autofire) | Button South (`A` / `X`) |
| **Deploy Recover (if stocked)**| `K` / `E` / Right Click | Button East (`B` / `Circle`) |
| **Fullscreen** | `F11` / `Alt+Enter` | - |
| **Restart (Game Over)** | `R` or `Space` | Start Button |

---

## 6. Visual Art Strategy (Zero-Asset Procedural)
- Built 100% in pure Go code using `procedural-art` principles:
  - Vector-rendered Gopher Vanguard ship with dual engine exhausts.
  - Distinct geometric silhouettes for all 5 Go enemy types.
  - Dynamic particle emitters for thrusters, sparks, and enemy explosions.
  - Retro CRT scanline and red warning flash overlay during Panic mode.
  - Infinite parallax background with scrolling cyber grid and code bits.

---

## 7. Audio Strategy (Zero-Asset Procedural Synthesis)
- Built 100% in pure Go code using DSP audio synthesis:
  - Laser blaster sound (pitch drop square wave).
  - Explosion crunch (shaped white noise with low-pass decay).
  - Panic alert siren (rapid two-tone frequency oscillation).
  - Recover chime (uplifting arpeggiated major chord).
  - High-energy chiptune bassline and lead loop.
