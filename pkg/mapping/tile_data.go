package mapping

import "github.com/hajimehoshi/ebiten/v2"

type TileData struct {
	Type     string
	Passable bool
	Sprite   *ebiten.Image
}
