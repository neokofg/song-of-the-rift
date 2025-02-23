package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

type Player struct {
	X, Y  float64
	Speed float64
	Image *ebiten.Image
}

func NewPlayer() *Player {
	img, _, err := ebitenutil.NewImageFromFile("pkg/assets/sprites/player.png")
	if err != nil {
		log.Fatal(err)
	}
	return &Player{
		X:     400,
		Y:     300,
		Speed: 5,
		Image: img,
	}
}

func (p *Player) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.Y -= p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.Y += p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.X -= p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.X += p.Speed
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.X, p.Y)
	screen.DrawImage(p.Image, op)
}
