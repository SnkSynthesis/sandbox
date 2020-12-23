package sand

import "github.com/hajimehoshi/ebiten/v2"

const Size = 1

type Particle struct {
	Img  *ebiten.Image
	X, Y float64
}
