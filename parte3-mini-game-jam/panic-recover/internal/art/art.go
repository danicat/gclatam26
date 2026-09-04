package art

import (
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	once sync.Once

	// Cached procedural images
	PlayerShip        *ebiten.Image
	DroneShip         *ebiten.Image
	EnemyNilPointer   *ebiten.Image
	EnemyConcurrent   *ebiten.Image
	EnemyDeadlock     *ebiten.Image
	EnemyMemoryLeak   *ebiten.Image
	EnemyGoroutine    *ebiten.Image
	BulletPlayer      *ebiten.Image
	BulletPanic       *ebiten.Image
	BulletEnemy       *ebiten.Image
	PickupRecover     *ebiten.Image
	PickupMutex       *ebiten.Image
	PickupWorker      *ebiten.Image
	ParticleGlow      *ebiten.Image
	ShieldAura        *ebiten.Image
	ScanlineOverlay   *ebiten.Image
	BossSigsegv       *ebiten.Image
)

// Init generates all procedural raster textures into GPU memory.
func Init() {
	once.Do(func() {
		PlayerShip = generatePlayerShip()
		DroneShip = generateDroneShip()
		BossSigsegv = generateBossSigsegv()
		EnemyNilPointer = generateEnemyNilPointer()
		EnemyConcurrent = generateEnemyConcurrent()
		EnemyDeadlock = generateEnemyDeadlock()
		EnemyMemoryLeak = generateEnemyMemoryLeak()
		EnemyGoroutine = generateEnemyGoroutine()
		BulletPlayer = generateBulletPlayer()
		BulletPanic = generateBulletPanic()
		BulletEnemy = generateBulletEnemy()
		PickupRecover = generatePickupRecover()
		PickupMutex = generatePickupMutex()
		PickupWorker = generatePickupWorker()
		ParticleGlow = generateParticleGlow()
		ShieldAura = generateShieldAura()
		ScanlineOverlay = generateScanlines()
	})
}

// ----------------------------------------------------------------------------
// Texture Generators
// ----------------------------------------------------------------------------

func generatePlayerShip() *ebiten.Image {
	w, h := 36, 36
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	// Go Cyan & Gopher Palette
	gopherBody := color.RGBA{0, 173, 216, 255}       // #00ADD8
	gopherOutline := color.RGBA{0, 115, 150, 255}    // #007396
	gopherInnerEar := color.RGBA{248, 187, 208, 255} // soft pink
	gopherMuzzle := color.RGBA{245, 230, 215, 255}   // cream snout
	gopherNose := color.RGBA{50, 35, 25, 255}        // dark brown nose
	white := color.RGBA{255, 255, 255, 255}
	pupil := color.RGBA{15, 20, 25, 255}
	wingMetal := color.RGBA{30, 50, 75, 255}
	wingEdge := color.RGBA{70, 180, 230, 255}
	thrusterCyan := color.RGBA{60, 220, 255, 255}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			fx, fy := float64(x), float64(y)

			// 1. Cyber Wings (behind the gopher body)
			if fy >= 16 && fy <= 31 {
				span := (fy - 16.0) * 1.25 // expands down
				if math.Abs(fx-18.0) <= span+4.0 && math.Abs(fx-18.0) >= 6.0 {
					img.Set(x, y, wingMetal)
					// Wing edge highlights
					if math.Abs(math.Abs(fx-18.0)-(span+4.0)) < 1.2 {
						img.Set(x, y, wingEdge)
					}
				}
			}

			// Wingtip blasters
			if (fx >= 3 && fx <= 5 || fx >= 30 && fx <= 32) && (fy >= 14 && fy <= 26) {
				img.Set(x, y, wingEdge)
			}

			// Dual engine thrusters
			if (fy >= 28 && fy <= 34) && ((fx >= 8 && fx <= 11) || (fx >= 24 && fx <= 27)) {
				img.Set(x, y, thrusterCyan)
			}

			// 2. Gopher Ears (Round at top sides)
			dLe := math.Hypot(fx-10.0, fy-10.0)
			if dLe <= 4.2 {
				img.Set(x, y, gopherBody)
				if dLe <= 2.2 {
					img.Set(x, y, gopherInnerEar)
				} else if dLe >= 3.4 {
					img.Set(x, y, gopherOutline)
				}
			}
			dRe := math.Hypot(fx-26.0, fy-10.0)
			if dRe <= 4.2 {
				img.Set(x, y, gopherBody)
				if dRe <= 2.2 {
					img.Set(x, y, gopherInnerEar)
				} else if dRe >= 3.4 {
					img.Set(x, y, gopherOutline)
				}
			}

			// 3. Gopher Main Body (Chubby round ellipse)
			dxBody := (fx - 18.0) / 10.0
			dyBody := (fy - 19.0) / 11.0
			distBody := dxBody*dxBody + dyBody*dyBody
			if distBody <= 1.0 {
				img.Set(x, y, gopherBody)
				if distBody >= 0.82 {
					img.Set(x, y, gopherOutline)
				}
			}

			// 4. Snout / Muzzle
			dxMuzzle := (fx - 18.0) / 4.5
			dyMuzzle := (fy - 20.5) / 3.0
			distMuzzle := dxMuzzle*dxMuzzle + dyMuzzle*dyMuzzle
			if distMuzzle <= 1.0 {
				img.Set(x, y, gopherMuzzle)
			}

			// 5. Two Buck Teeth!
			if (fy >= 22 && fy <= 25) && (fx == 16 || fx == 17 || fx == 19 || fx == 20) {
				img.Set(x, y, white)
			}
			if (fy >= 22 && fy <= 25) && fx == 18 {
				img.Set(x, y, gopherOutline) // Separator
			}

			// 6. Cute Nose
			if math.Hypot(fx-18.0, fy-19.0) <= 1.5 {
				img.Set(x, y, gopherNose)
			}

			// 7. Iconic Big Eyes
			dEyeL := math.Hypot(fx-13.0, fy-14.5)
			if dEyeL <= 3.8 {
				img.Set(x, y, white)
				if dEyeL >= 3.2 {
					img.Set(x, y, gopherOutline)
				}
				if math.Hypot(fx-13.2, fy-14.0) <= 1.7 {
					img.Set(x, y, pupil)
				}
				if x == 12 && y == 13 {
					img.Set(x, y, white)
				}
			}

			dEyeR := math.Hypot(fx-23.0, fy-14.5)
			if dEyeR <= 3.8 {
				img.Set(x, y, white)
				if dEyeR >= 3.2 {
					img.Set(x, y, gopherOutline)
				}
				if math.Hypot(fx-22.8, fy-14.0) <= 1.7 {
					img.Set(x, y, pupil)
				}
				if x == 22 && y == 13 {
					img.Set(x, y, white)
				}
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

func generateBossSigsegv() *ebiten.Image {
	w, h := 64, 64
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	cx, cy := 32.0, 32.0

	hullDark := color.RGBA{22, 24, 30, 255}
	hullPlate := color.RGBA{45, 48, 60, 255}
	crimson := color.RGBA{235, 30, 50, 255}
	brightCrimson := color.RGBA{255, 80, 80, 255}
	coreWhite := color.RGBA{255, 235, 240, 255}
	goldWarning := color.RGBA{255, 195, 30, 255}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			fx, fy := float64(x), float64(y)
			dx := math.Abs(fx - cx)
			dy := math.Abs(fy - cy)

			// Armored dreadnought octagon hull
			if dx+dy <= 28.0 && dx <= 26.0 && dy <= 24.0 {
				img.Set(x, y, hullDark)

				// Armor plates
				if dx+dy <= 24.0 && (int(fx+fy)%8 < 6) {
					img.Set(x, y, hullPlate)
				}

				// Gold warning stripes
				if (fy >= 18 && fy <= 21) && (int(fx+fy)%6 < 3) {
					img.Set(x, y, goldWarning)
				}

				// Outer crimson trim
				if dx+dy >= 26.0 || dx >= 24.0 {
					img.Set(x, y, crimson)
				}
			}

			// Twin heavy plasma cannons on sides
			if (dx >= 24.0 && dx <= 30.0) && (fy >= 22.0 && fy <= 44.0) {
				img.Set(x, y, hullDark)
				if fy >= 40.0 {
					img.Set(x, y, crimson) // Cannon nozzles
				}
			}

			// Forward prow horns
			if (fy >= 4.0 && fy <= 18.0) && (dx >= 12.0 && dx <= 18.0) {
				if fy >= 18.0-(dx-12.0)*2.0 {
					img.Set(x, y, crimson)
				}
			}

			// Central pulsing SIGSEGV Core (Reactor)
			distCenter := math.Hypot(fx-cx, fy-cy)
			if distCenter <= 10.0 {
				img.Set(x, y, crimson)
				if distCenter <= 7.0 {
					img.Set(x, y, brightCrimson)
				}
				if distCenter <= 3.5 {
					img.Set(x, y, coreWhite)
				}
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

func generateDroneShip() *ebiten.Image {
	w, h := 16, 16
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx := float64(x - 8)
			dy := float64(y - 8)
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist <= 6.0 {
				r, g, b := uint8(0), uint8(210), uint8(255)
				if dist <= 2.5 {
					r, g, b = 255, 255, 255 // Core highlight
				}
				img.Set(x, y, color.RGBA{r, g, b, 255})
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

func generateEnemyNilPointer() *ebiten.Image {
	// Sharp purple razor stealth dart pointing downward
	w, h := 26, 26
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	cx := 13.0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			fx, fy := float64(x), float64(y)
			dx := math.Abs(fx - cx)
			// Inverted delta (nose at bottom y: 22)
			if fy >= 4 && fy <= 22 {
				wingWidth := 1.5 + (22.0-fy)*0.6
				if dx <= wingWidth {
					r, g, b := uint8(190), uint8(40), uint8(240) // Vivid Purple
					if dx < wingWidth*0.3 {
						r, g, b = 245, 160, 255 // Violet highlight
					}
					// Red glitch core
					if fy >= 10 && fy <= 14 && dx <= 2.0 {
						r, g, b = 255, 50, 50
					}
					img.Set(x, y, color.RGBA{r, g, b, 255})
				}
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

func generateEnemyConcurrent() *ebiten.Image {
	// Dual red catamaran interceptor linked by energy arc
	w, h := 30, 26
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			fx, fy := float64(x), float64(y)
			// Left hull (center x=7), Right hull (center x=23)
			dxLeft := math.Abs(fx - 7.0)
			dxRight := math.Abs(fx - 23.0)

			isHull := false
			if fy >= 4 && fy <= 22 {
				if dxLeft <= 3.5 || dxRight <= 3.5 {
					isHull = true
				}
			}
			// Cross bridge conduit (y: 11..14, x: 7..23)
			isBridge := (fy >= 11 && fy <= 14 && fx >= 7 && fx <= 23)

			if isHull {
				img.Set(x, y, color.RGBA{230, 40, 50, 255})
			} else if isBridge {
				img.Set(x, y, color.RGBA{255, 180, 50, 255}) // Golden energy connection
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

func generateEnemyDeadlock() *ebiten.Image {
	// Heavy fortified bronze bunker with mutex symbol
	w, h := 34, 34
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	cx, cy := 17.0, 17.0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist <= 14.0 {
				r, g, b := uint8(190), uint8(130), uint8(40) // Heavy Bronze
				if dist > 11.0 {
					r, g, b = 120, 80, 25 // Dark border
				} else if dist <= 6.0 {
					// Mutex Lock Core
					r, g, b = 255, 215, 80
				}
				img.Set(x, y, color.RGBA{r, g, b, 255})
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

func generateEnemyMemoryLeak() *ebiten.Image {
	// Pulsating green organic slime blob
	w, h := 28, 28
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	cx, cy := 14.0, 14.0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			angle := math.Atan2(dy, dx)
			// Blobby deformed radius
			rMax := 10.0 + 2.0*math.Sin(4.0*angle)
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist <= rMax {
				r, g, b := uint8(40), uint8(220), uint8(80)
				if dist <= 4.0 {
					r, g, b = 200, 255, 120 // Glowing radioactive center
				}
				img.Set(x, y, color.RGBA{r, g, b, 255})
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

func generateEnemyGoroutine() *ebiten.Image {
	// Fast orange hornet / diamond drone
	w, h := 18, 18
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	cx, cy := 9.0, 9.0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx := math.Abs(float64(x) - cx)
			dy := math.Abs(float64(y) - cy)
			// Diamond shape (dx + dy <= 7)
			if dx+dy <= 7.0 {
				r, g, b := uint8(255), uint8(120), uint8(20)
				if dx+dy <= 3.0 {
					r, g, b = 255, 240, 160
				}
				img.Set(x, y, color.RGBA{r, g, b, 255})
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

func generateBulletPlayer() *ebiten.Image {
	w, h := 6, 14
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if x == 0 || x == 5 {
				img.Set(x, y, color.RGBA{80, 200, 255, 180})
			} else {
				img.Set(x, y, color.RGBA{220, 250, 255, 255}) // Bright laser core
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

func generateBulletPanic() *ebiten.Image {
	w, h := 8, 18
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if x == 0 || x == 7 {
				img.Set(x, y, color.RGBA{255, 60, 20, 200})
			} else {
				img.Set(x, y, color.RGBA{255, 230, 100, 255})
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

func generateBulletEnemy() *ebiten.Image {
	w, h := 8, 8
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx := float64(x - 4)
			dy := float64(y - 4)
			if dx*dx+dy*dy <= 12.0 {
				img.Set(x, y, color.RGBA{255, 40, 120, 255})
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

func generatePickupRecover() *ebiten.Image {
	// Green health cross / recovery badge
	w, h := 18, 18
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	cx, cy := 9.0, 9.0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx := math.Abs(float64(x) - cx)
			dy := math.Abs(float64(y) - cy)
			// Cross shape
			isCross := (dx <= 2.0 && dy <= 6.0) || (dy <= 2.0 && dx <= 6.0)
			if isCross {
				img.Set(x, y, color.RGBA{50, 255, 120, 255})
			} else if dx <= 7.0 && dy <= 7.0 {
				// Outer subtle ring
				img.Set(x, y, color.RGBA{20, 80, 40, 180})
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

func generatePickupMutex() *ebiten.Image {
	// Golden Shield
	w, h := 18, 18
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			fx, fy := float64(x), float64(y)
			dx := math.Abs(fx - 9.0)
			if fy >= 3 && fy <= 15 && dx <= 6.0 {
				img.Set(x, y, color.RGBA{255, 200, 30, 255})
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

func generatePickupWorker() *ebiten.Image {
	// Cyan Mini Drone icon
	w, h := 18, 18
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx := math.Abs(float64(x) - 9.0)
			dy := math.Abs(float64(y) - 9.0)
			if dx+dy <= 6.0 {
				img.Set(x, y, color.RGBA{40, 220, 255, 255})
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

func generateParticleGlow() *ebiten.Image {
	w, h := 16, 16
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx := float64(x - 8)
			dy := float64(y - 8)
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist < 8.0 {
				alpha := uint8((1.0 - dist/8.0) * 255.0)
				img.Set(x, y, color.RGBA{255, 255, 255, alpha})
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

func generateShieldAura() *ebiten.Image {
	w, h := 48, 48
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	cx, cy := 24.0, 24.0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist >= 18.0 && dist <= 23.0 {
				alpha := uint8((1.0 - math.Abs(dist-20.5)/2.5) * 200.0)
				img.Set(x, y, color.RGBA{255, 215, 60, alpha})
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}

func generateScanlines() *ebiten.Image {
	w, h := 640, 360
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		if y%2 == 0 {
			for x := 0; x < w; x++ {
				img.Set(x, y, color.RGBA{0, 0, 0, 45})
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}
