package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"github.com/SnkSynthesis/sandbox/sand"
)

const zoom = 10
const WindowHeight, WindowWidth = 800, 800
const BoxWidth, BoxHeight = WindowWidth/zoom, WindowHeight/zoom

type Game struct {
	img *ebiten.Image
	op *ebiten.DrawImageOptions
	particles []*sand.Particle
}

type SandParticle struct {
	img *ebiten.Image
	x, y float64
}

func (g *Game) Init() {
	g.op = &ebiten.DrawImageOptions{}
	g.img = ebiten.NewImage(sand.Size, sand.Size)
	g.img.Fill(color.RGBA{255, 255, 0, 255})

	g.particles = make([]*sand.Particle, BoxWidth*BoxHeight)

	for i := 0; i <= 5; i += sand.Size {
		p := &sand.Particle{g.img, float64(i), float64(i), 0, 0}
		g.particles[i * BoxWidth + i] = p
	}
}

func (g *Game) Update() error {
	
	x, y := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if x > 0 && y > 0 && x < BoxWidth && y < BoxHeight {
			p := &sand.Particle{g.img, float64(x), float64(y), 0, 0}
			g.particles[y * BoxWidth + x] = p
		}
	}

	for i, _ := range g.particles {
		if g.particles[i] != nil {
			p := g.particles[i]
			if p.Y < float64(BoxHeight-sand.Size) {
				// p.Y += 1
				g.particles[int(p.Y * BoxWidth + p.X)] = p
				g.particles[i] = nil
			} else {
				p.Y = float64(BoxHeight-sand.Size)
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
	return outsideWidth/zoom, outsideHeight/zoom
}

func main() {
	ebiten.SetMaxTPS(10)
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle("Sandbox")

	game := &Game{}
	game.Init()

	fmt.Println("Starting game...")

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
