package mapping

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/neokofg/mygame/pkg/ecs"
	"image/color"
	"log"
)

type TileMap struct {
	Width  int
	Height int
	Tiles  [][]TileData
}

func CreateTileMapEntity(em *ecs.EntityManager, width, height int) *ecs.Entity {
	tileMap := &TileMap{
		Width:  width,
		Height: height,
		Tiles:  make([][]TileData, height),
	}

	floorTile, _, err := ebitenutil.NewImageFromFile("pkg/assets/sprites/floor.png")
	if err != nil {
		log.Fatal(err)
	}

	wallTile := createGrayTile(64, 64)

	for y := 0; y < height; y++ {
		tileMap.Tiles[y] = make([]TileData, width)
		for x := 0; x < width; x++ {
			if x == 0 || x == width-1 || y == 0 || y == height-1 {
				tileMap.Tiles[y][x] = TileData{
					Type:     "wall",
					Passable: false,
					Sprite:   wallTile,
				}
			} else {
				tileMap.Tiles[y][x] = TileData{
					Type:     "floor",
					Passable: true,
					Sprite:   floorTile,
				}
			}
		}
	}

	entity := em.CreateEntity()
	entity.AddComponent("TileMap", tileMap)
	return entity
}

func createGrayTile(width, height int) *ebiten.Image {
	img := ebiten.NewImage(width, height)
	img.Fill(color.Gray{128})
	return img
}
