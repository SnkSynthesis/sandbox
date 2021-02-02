package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

var (
	zoom                      = 8
	WindowHeight, WindowWidth = 500, 500
	BoxWidth, BoxHeight       = (WindowWidth / zoom), (WindowHeight / zoom)
)

const ParticleSize = 1
type Particle struct {
	Img  *ebiten.Image
	X, Y float64
}

type Game struct {
	img       *ebiten.Image
	col       uint8
	bCol      bool
	op        *ebiten.DrawImageOptions
	particles []*Particle
}

func (g *Game) Init() {
	g.op = &ebiten.DrawImageOptions{}
	g.img = ebiten.NewImage(ParticleSize, ParticleSize)
	g.img.Fill(color.RGBA{255, 255, 0, 255})

	g.particles = make([]*Particle, BoxWidth*BoxHeight)
}

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.particles = make([]*Particle, len(g.particles))
	}

	WindowWidth, WindowHeight = ebiten.WindowSize()
	BoxWidth, BoxHeight = (WindowWidth / zoom), (WindowHeight / zoom)

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
	g.img.Fill(color.RGBA{g.col, g.col, 0, 255})

	x, y := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if x > 0-ParticleSize && y > 0-ParticleSize && x < BoxWidth && y < BoxHeight {
			g.particles[y*BoxWidth+x] = &Particle{g.img, float64(x), float64(y)}
		}
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		if x > 0-ParticleSize && y > 0-ParticleSize && x < BoxWidth && y < BoxHeight {
			g.particles[y*BoxWidth+x] = nil
		}
	}

	doNotChange := map[int]bool{}

	for i, _ := range g.particles {
		if g.particles[i] != nil {
			p := g.particles[i]

			_, ok := doNotChange[i]
			if ok {
				continue
			}

			if p.Y < float64(BoxHeight-ParticleSize) {
				if int(p.Y+1)*BoxWidth+int(p.X) < len(g.particles) && g.particles[int(p.Y+1)*BoxWidth+int(p.X)] == nil {
					p.Y += 1
					j := int(p.Y)*BoxWidth + int(p.X)
					g.particles[j] = p
					doNotChange[j] = true
					g.particles[i] = nil
				} else if int(p.X+1) < BoxWidth-ParticleSize && int(p.Y+1)*BoxWidth+int(p.X+1) < len(g.particles) && g.particles[int(p.Y+1)*BoxWidth+int(p.X+1)] == nil {
					p.X += 1
					p.Y += 1
					j := int(p.Y)*BoxWidth + int(p.X)
					g.particles[j] = p
					doNotChange[j] = true
					g.particles[i] = nil
				} else if p.X-1 >= 0 && int(p.Y+1)*BoxWidth+int(p.X-1) < len(g.particles) && g.particles[int(p.Y+1)*BoxWidth+int(p.X-1)] == nil {
					p.X -= 1
					p.Y += 1
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

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.particles {
		if p != nil {
			g.op.GeoM.Reset()
			g.op.GeoM.Translate(p.X, p.Y)
			screen.DrawImage(p.Img, g.op)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth / zoom, outsideHeight / zoom
}

func main() {
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle("Sandbox")
	ebiten.SetWindowResizable(true)

	fmt.Print("Enter zoom amount (Integers 1 and above; default: 8): ")
	fmt.Scanf("%d", &zoom)

	if zoom < 1 {
		fmt.Println("Invalid number! Defaulting to 8...")
		zoom = 8
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
