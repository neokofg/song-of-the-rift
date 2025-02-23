package components

import "github.com/hajimehoshi/ebiten/v2"

type Camera struct {
	X, Y, Width, Height, Zoom     float64
	DeadZoneWidth, DeadZoneHeight float64
	Smoothness                    float64
}

func (c *Camera) ApplyTransform(g *ebiten.GeoM) {
	g.Translate(-c.X, -c.Y)
	g.Scale(c.Zoom, c.Zoom)
}
