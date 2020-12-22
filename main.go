package main

import (
	"fmt"
	"github.com/SnkSynthesis/sandbox/sand"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

const zoom = 10
const WindowHeight, WindowWidth = 800, 800
const BoxWidth, BoxHeight = WindowWidth / zoom, WindowHeight / zoom

type Game struct {
	img       *ebiten.Image
	op        *ebiten.DrawImageOptions
	particles []*sand.Particle
}

func (g *Game) Init() {
	g.op = &ebiten.DrawImageOptions{}
	g.img = ebiten.NewImage(sand.Size, sand.Size)
	g.img.Fill(color.RGBA{255, 255, 0, 255})

	g.particles = make([]*sand.Particle, BoxWidth*BoxHeight)

	for i := 0; i <= 5; i += sand.Size {
		p := &sand.Particle{g.img, float64(i), float64(i)}
		g.particles[i*BoxWidth+i] = p
	}
}

func (g *Game) Update() error {

	x, y := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if x > 0 - sand.Size && y > 0 - sand.Size && x < BoxWidth + sand.Size && y < BoxHeight + sand.Size {
			p := &sand.Particle{g.img, float64(x), float64(y)}
			g.particles[y*BoxWidth+x] = p
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

			if p.Y < float64(BoxHeight-sand.Size) {

				if g.particles[int((p.Y + 1) * BoxWidth + p.X)] == nil {
					p.Y += 1
					j := int(p.Y * BoxWidth + p.X)
					g.particles[j] = p
					doNotChange[j] = true
					g.particles[i] = nil
				} else {
					doNotChange[i] = true
				}
				
			} else {
				p.Y = float64(BoxHeight - sand.Size)
				j := int(p.Y * BoxWidth + p.X)
				g.particles[j] = p
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

	game := &Game{}
	game.Init()

	fmt.Println("Starting game...")

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
