package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// Zoom controls the zoom amount
	Zoom = 8

	// WindowWidth is the width of the window
	WindowWidth = 500
	// WindowHeight is the height of the window
	WindowHeight = 500

	// BoxWidth is the width of the underlying array of particles
	BoxWidth = WindowWidth / Zoom
	// BoxHeight is the height of the underlying array of particles
	BoxHeight = WindowHeight / Zoom
)

// ParticleSize is the size of the sand particles
const ParticleSize = 1

// Particle holds the data each particle contains like its image, x position, and y position
type Particle struct {
	Img  *ebiten.Image
	X, Y float64
}

// Game is the main game struct holding the game state
type Game struct {
	particleImg *ebiten.Image
	col         uint8
	bCol        bool
	op          *ebiten.DrawImageOptions
	particles   []*Particle
}

// Init is the initialization function of Game
func (g *Game) Init() {
	g.op = &ebiten.DrawImageOptions{}
	g.particleImg = ebiten.NewImage(ParticleSize, ParticleSize)
	g.particleImg.Fill(color.RGBA{255, 255, 0, 255})

	g.particles = make([]*Particle, BoxWidth*BoxHeight)
}

// Update does the update logic of the game
func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.particles = make([]*Particle, len(g.particles))
	}

	WindowWidth, WindowHeight = ebiten.WindowSize()
	BoxWidth, BoxHeight = (WindowWidth/Zoom)+1, (WindowHeight/Zoom)+1

	diff := BoxWidth*BoxHeight - len(g.particles)
	if diff != 0 {
		g.particles = make([]*Particle, len(g.particles)+diff)
	}

	if g.bCol {
		g.col--
		if g.col == 150 {
			g.bCol = false
		}
	} else {
		g.col++
		if g.col == 255 {
			g.bCol = true
		}
	}
	g.particleImg.Fill(color.RGBA{g.col, g.col, 0, 255})

	x, y := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if x > 0-ParticleSize && y > 0-ParticleSize && x < BoxWidth && y < BoxHeight {
			g.particles[y*BoxWidth+x] = &Particle{g.particleImg, float64(x), float64(y)}
		}
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		if x > 0-ParticleSize && y > 0-ParticleSize && x < BoxWidth && y < BoxHeight {
			g.particles[y*BoxWidth+x] = nil
		}
	}

	doNotChange := map[int]bool{}

	for i := range g.particles {
		if g.particles[i] != nil {
			p := g.particles[i]

			_, ok := doNotChange[i]
			if ok {
				continue
			}

			if p.Y < float64(BoxHeight-ParticleSize) {
				if int(p.Y+1)*BoxWidth+int(p.X) < len(g.particles) && g.particles[int(p.Y+1)*BoxWidth+int(p.X)] == nil {
					p.Y++
					j := int(p.Y)*BoxWidth + int(p.X)
					g.particles[j] = p
					doNotChange[j] = true
					g.particles[i] = nil
				} else if int(p.X+1) < BoxWidth-ParticleSize && int(p.Y+1)*BoxWidth+int(p.X+1) < len(g.particles) && g.particles[int(p.Y+1)*BoxWidth+int(p.X+1)] == nil {
					p.X++
					p.Y++
					j := int(p.Y)*BoxWidth + int(p.X)
					g.particles[j] = p
					doNotChange[j] = true
					g.particles[i] = nil
				} else if p.X-1 >= 0 && int(p.Y+1)*BoxWidth+int(p.X-1) < len(g.particles) && g.particles[int(p.Y+1)*BoxWidth+int(p.X-1)] == nil {
					p.X--
					p.Y++
					j := int(p.Y)*BoxWidth + int(p.X)
					g.particles[j] = p
					doNotChange[j] = true
					g.particles[i] = nil
				} else {
					doNotChange[i] = true
				}
			} else {
				p.Y = float64(BoxHeight - ParticleSize)
				j := int(p.Y)*BoxWidth + int(p.X)
				if j > 0 && j < len(g.particles) {
					g.particles[j] = p
				}
			}
		}
	}

	return nil
}

// Draw has code that draws the particles and other items onto the window
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{97, 202, 255, 255})
	for _, p := range g.particles {
		if p != nil {
			g.op.GeoM.Reset()
			g.op.GeoM.Translate(p.X, p.Y)
			screen.DrawImage(p.Img, g.op)
		}
	}
}

// Layout takes outside side and divides by Zoom giving the zoom effect
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth / Zoom, outsideHeight / Zoom
}

func main() {
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle("Sandbox")
	ebiten.SetWindowResizable(true)

	fmt.Print("Enter zoom amount (Integers 1 and above; default: 8): ")
	fmt.Scanf("%d", &Zoom)

	if Zoom < 1 {
		fmt.Println("Invalid number! Defaulting to 8...")
		Zoom = 8
	}

	fmt.Println()
	fmt.Println(" --- Instructions --- ")
	fmt.Println("[Space] Or [Resizing Window] - Clears all sand")
	fmt.Println("[Left-Click] - Places sand")
	fmt.Println("[Right-Click] - Removes sand")
	fmt.Println(" --- Instructions --- ")
	fmt.Println()

	fmt.Println("Enjoy!")
	fmt.Println()

	game := &Game{}
	game.Init()

	fmt.Println("Starting game...")

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
