package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/neokofg/mygame/pkg/components"
	"github.com/neokofg/mygame/pkg/ecs"
	"github.com/neokofg/mygame/pkg/entities"
	"github.com/neokofg/mygame/pkg/systems"
	"log"
)

type Game struct {
	EntityManager *ecs.EntityManager
	Systems       []interface{}
}

func NewGame() *Game {
	em := ecs.NewEntityManager()

	entities.NewPlayer(em)

	gameSystems := []interface{}{
		&systems.InputSystem{},
		&systems.MovementSystem{},
		&systems.RenderSystem{},
		&systems.VelocitySystem{},
		&systems.CameraSystem{},
	}

	return &Game{
		EntityManager: em,
		Systems:       gameSystems,
	}
}
func (g *Game) Update() error {
	gameEntities := g.EntityManager.GetAllEntities()
	for _, sys := range g.Systems {
		if updateSys, ok := sys.(interface {
			Update([]*ecs.Entity)
		}); ok {
			updateSys.Update(gameEntities)
		}
	}
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	gameEntities := g.EntityManager.GetAllEntities()

	var cameraGeoM ebiten.GeoM
	for _, entity := range gameEntities {
		if entity.HasComponent("Camera") {
			camera := entity.GetComponent("Camera").(*components.Camera)
			camera.ApplyTransform(&cameraGeoM)
			break
		}
	}

	floorTile, _, err := ebitenutil.NewImageFromFile("pkg/assets/sprites/floor.png")
	if err != nil {
		log.Fatal(err)
	}
	tileWidth := 64
	tileHeight := 64
	for x := 0; x < 800; x += tileWidth {
		for y := 0; y < 600; y += tileHeight {
			op := &ebiten.DrawImageOptions{}
			op.GeoM = cameraGeoM
			op.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(floorTile, op)
		}
	}

	for _, sys := range g.Systems {
		if drawSys, ok := sys.(interface {
			Draw(*ebiten.Image, []*ecs.Entity, ebiten.GeoM)
		}); ok {
			drawSys.Draw(screen, gameEntities, cameraGeoM)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}
