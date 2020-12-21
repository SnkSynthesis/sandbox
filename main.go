package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Game struct{}

var img *ebiten.Image
var imgX, imgY float64

func (g *Game) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		imgY -= 10
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		imgY += 10
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		imgX -= 10
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		imgX += 10
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(imgX, imgY)
	screen.DrawImage(img, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(700, 700)
	ebiten.SetWindowTitle("Sandbox")

	game := &Game{}

	img = ebiten.NewImage(30, 30)
	img.Fill(color.RGBA{255, 0, 0, 255})

	fmt.Println("Starting game...")

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
