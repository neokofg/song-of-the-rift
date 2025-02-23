package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/neokofg/mygame/pkg/entities"
	"log"
)

type Game struct {
	Player  *entities.Player
	Floor   *ebiten.Image
	CameraX float64
	CameraY float64
}

func NewGame() *Game {
	floorImg, _, err := ebitenutil.NewImageFromFile("pkg/assets/sprites/floor.png")
	if err != nil {
		log.Fatal(err)
	}
	return &Game{
		Player: entities.NewPlayer(),
		Floor:  floorImg,
	}
}
func (g *Game) Update() error {
	g.Player.Update()
	g.CameraX = g.Player.X - 400
	g.CameraY = g.Player.Y - 300
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	tileWidth := 64.0
	tileHeight := 64.0

	for x := -tileWidth; x < 800; x += tileWidth {
		for y := -tileHeight; y < 600; y += tileHeight {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(x-g.CameraX, y-g.CameraY)
			screen.DrawImage(g.Floor, op)
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(400, 300)
	screen.DrawImage(g.Player.Image, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}
