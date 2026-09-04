# Game Design Document (GDD): Panic!!! (& recover?)

<p align="center">
  <img src="assets/promo_banner.jpg" alt="Panic!!! (& recover?) Promo Banner" width="100%">
</p>

> **Game Title**: Panic!!! (& recover?)  
> **Target Genre**: Arcade Code-Puzzle / Typing Reflex  
> **Target Platform**: Desktop (Windows/macOS/Linux) & WebAssembly (Google Cloud Run)  
> **Target Aspect Ratio & Resolution**: 16:9 Widescreen (`640x360` Virtual Pixel Canvas)  
> **Author / Lead Designer**: GopherCon LATAM 2026 Mini Game Jam Team  

---

## 1. Executive Summary & Elevator Pitch
- **Elevator Pitch**: A high-stakes arcade puzzle game where buggy Go code descends toward a terminal crash threshold like Tetris blocks—race against the panic countdown timer, select the fatal line, rewrite the bug, and trigger `recover()` before the runtime panics!
- **Core Inspiration**: *Tetris* meets *Go Runtime Debugger* meets *Mavis Beacon / Typing Arcade*.
- **Target Audience / Mood**: Urgent, humorous, educational, and intensely satisfying for gophers and developers who dread production panics.

---

## 2. Core Gameplay Loop & Mechanics
- **Primary Gameplay Loop**:
  1. **Spawn & Fall**: A Go code snippet descends from the top of the terminal at a constant rate, accompanied by an urgent countdown timer.
  2. **Navigate & Identify**: The player moves a selection cursor with `Up`/`Down` arrow keys over the falling lines to spot the bug causing a crash/panic.
  3. **Edit & Rewrite**: Pressing `Enter` or `Space` engages inline editing mode with a blinking cursor. The player rewrites the buggy statement or provides the patch.
  4. **Evaluate & Recover**: Pressing `Enter` validates the submitted line against valid bug patches:
     - **Success (`recover()`)**: Panic is averted! Score bonus awarded based on remaining time. Sound fanfare plays, screen sparks green, and game advances to the next level.
     - **Failure**: Syntax/runtime warning flashes, -2s time penalty applied, alarm beeps, player can retry immediately.
  5. **Panic Loss**: If the timer hits `0.0s` or the code reaches the bottom crash threshold, runtime panic occurs (`panic: runtime error`), dumping a fatal stack trace.

- **Difficulty Progression**:
  - 10 distinct levels covering real Go runtime crashes (nil pointer, divide by zero, index out of range, nil map write, closed channel send, recursive stack overflow, etc.).
  - Fall speed and timer duration get progressively tighter as levels increase (Level 1: 40s fall time, Level 10: 15s fall time).

---

## 3. Controls & Input Mapping Scheme
| Logical Action | Keyboard | Mobile / Touch | Description |
| :--- | :--- | :--- | :--- |
| `ActionSelectUp` | `Up` / `W` | Swipe Up / Tap Prev | Move cursor to previous code line |
| `ActionSelectDown` | `Down` / `S` | Swipe Down / Tap Next | Move cursor to next code line |
| `ActionEditOrSubmit` | `Enter` / `Space` | Tap Edit/Submit Button | Enter edit mode OR submit rewritten line |
| `ActionCancelEdit` | `Escape` | Tap Cancel Button | Cancel edit mode and return to navigation |
| `ActionHint` | `H` / `F1` | Tap Hint Button | Reveal a diagnostic hint about the bug |
| `ActionToggleFullscreen`| `F11` | N/A | Toggle Fullscreen mode |

---

## 4. Visual Style & Asset Strategy
- **Aesthetic Direction**: Retro Cyberpunk IDE / CRT Monospace Terminal.
- **Procedural Pure-Code (`procedural-art`)**:
  - Monospace font rendering with custom syntax highlighting (keywords `func`, `var`, `defer`, `go`, `if`, `return` in cyan; types in yellow; strings in green; comments in gray; bug highlight in neon red).
  - CRT scanlines and subtle curvature vignette.
  - Laser "PANIC THRESHOLD" hazard line glowing at the bottom of the screen.
  - Terminal glitch and red warning flash as timer approaches zero.
  - Victory green matrix spark particle emitter upon successful `recover()`.
  - Zero external bitmap assets: 100% generated in Go memory.

---

## 5. Audio & Soundscape Strategy
- **Procedural Sound Engine (`procedural-composer`)**:
  - Zero external `.wav` or `.mp3` files: 16-bit 44.1 kHz PCM synthesized dynamically.
  - **Mechanical Keyboard Click**: Short crisp noise transient on every keypress.
  - **Cursor Blip**: 880 Hz sine blip when navigating lines.
  - **Panic Siren**: Pulsing 440 Hz / 880 Hz square wave alarm when time < 5 seconds.
  - **Panic Crash / Dump**: Deep descending sawtooth noise burst on game over.
  - **Recover Fanfare**: Ascending major arpeggio (`C5 -> E5 -> G5 -> C6`) on level clear.
  - **Cyberpunk BGM**: Ambient FM synthesized bassline and arpeggio pulse.

---

## 6. Level Progression (10 Core Scenarios)
1. **Level 1**: Nil Pointer Dereference (`var u *User; u.Name = "Gopher"` $\rightarrow$ `u := &User{Name: "Gopher"}`)
2. **Level 2**: Divide by Zero (`ratio := total / count` $\rightarrow$ `ratio := total / (count + 1)`)
3. **Level 3**: Slice Index Out of Range (`val := items[len(items)]` $\rightarrow$ `val := items[len(items)-1]`)
4. **Level 4**: Nil Map Assignment (`var m map[string]int; m["hp"] = 100` $\rightarrow$ `m := make(map[string]int); m["hp"] = 100`)
5. **Level 5**: Send on Closed Channel (`close(ch); ch <- 42` $\rightarrow$ `ch <- 42; close(ch)`)
6. **Level 6**: Close of Closed Channel (`close(ch); close(ch)` $\rightarrow$ `close(ch) // closed once`)
7. **Level 7**: Infinite Recursion Stack Overflow (`func loop(n int) { loop(n+1) }` $\rightarrow$ `if n >= 10 { return }`)
8. **Level 8**: Unchecked Type Assertion Panic (`s := val.(string)` $\rightarrow$ `s, ok := val.(string)`)
9. **Level 9**: Concurrent Map Writes (`go write(m); go write(m)` $\rightarrow$ `mu.Lock(); m[k]=v; mu.Unlock()`)
10. **Level 10**: Panic Without Recover (`panic("fatal!")` $\rightarrow$ `defer func() { recover() }()`)

---

## 7. Technical Scope & Architecture Notes
- **Go Version**: Go 1.26+
- **Engine**: Ebitengine v2 (`github.com/hajimehoshi/ebitengine/v2`)
- **Package Layout**:
  - `panic-recover/main.go`
  - `panic-recover/internal/game/` (Game loop, Scene FSM, Level Manager)
  - `panic-recover/internal/editor/` (Text input controller, line buffer, syntax parsing)
  - `panic-recover/internal/art/` (Procedural CRT renderer, syntax coloring, particles)
  - `panic-recover/internal/audio/` (Procedural DSP sound synthesis)
- **Deployment**: Local desktop app + WebAssembly (`GOOS=js GOARCH=wasm`)
