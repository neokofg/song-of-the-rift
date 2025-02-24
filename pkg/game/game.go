package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/neokofg/mygame/pkg/components"
	"github.com/neokofg/mygame/pkg/ecs"
	"github.com/neokofg/mygame/pkg/entities"
	"github.com/neokofg/mygame/pkg/leveling"
	"github.com/neokofg/mygame/pkg/mapping"
	"github.com/neokofg/mygame/pkg/systems"
	"log"
	"math"
	"time"
)

type Game struct {
	EntityManager  *ecs.EntityManager
	Systems        []interface{}
	LevelManager   *leveling.LevelManager
	Accumulator    float64
	FixedDeltaTime float64
	LastUpdateTime time.Time
}

func NewGame() *Game {
	em := ecs.NewEntityManager()

	lm := &leveling.LevelManager{
		CurrentLevel:    0,
		Levels:          make(map[int]*leveling.Level),
		EntityManager:   em,
		CurrentEntities: []*ecs.Entity{},
	}

	level1 := &leveling.Level{
		ID: 1,
		CreateEntities: func(em *ecs.EntityManager) []*ecs.Entity {
			tileMap := mapping.CreateTileMapEntity(em, 100, 100)

			player := entities.NewPlayer(em)
			spawnPoint := tileMap.GetComponent("SpawnPoint").(struct{ X, Y int })
			playerPos := player.GetComponent("Position").(*components.Position)
			tileSize := 64
			playerPos.X = float64(spawnPoint.X * tileSize)
			playerPos.Y = float64(spawnPoint.Y * tileSize)
			
			return []*ecs.Entity{player, tileMap}
		},
	}
	lm.Levels[1] = level1
	lm.LoadLevel(1)

	gameSystems := []interface{}{
		&systems.InputSystem{},
		&systems.CameraSystem{},
		&systems.RenderSystem{},
		&systems.VelocitySystem{},
		&systems.CollisionSystem{},
		&systems.MovementSystem{},
	}
	g := &Game{
		EntityManager:  em,
		Systems:        gameSystems,
		LevelManager:   lm,
		FixedDeltaTime: 1.0 / 60,
		Accumulator:    0,
		LastUpdateTime: time.Now(),
	}
	log.Printf("Игра инициализирована, сущностей: %d", len(g.EntityManager.GetAllEntities()))
	return g
}
func (g *Game) Update() error {
	currentTime := time.Now()
	frameTime := currentTime.Sub(g.LastUpdateTime).Seconds()
	g.LastUpdateTime = currentTime

	if frameTime < 0 {
		log.Println("Отрицательный frameTime, устанавливаем в 0:", frameTime)
		frameTime = 0
	}

	g.Accumulator += frameTime
	gameEntities := g.EntityManager.GetAllEntities()

	for _, sys := range g.Systems {
		if variableSys, ok := sys.(ecs.VariableUpdateSystem); ok {
			variableSys.Update(gameEntities)
		}
	}

	maxUpdatesPerFrame := 5
	updateCount := 0
	for g.Accumulator >= g.FixedDeltaTime && updateCount < maxUpdatesPerFrame {
		for _, sys := range g.Systems {
			if fixedSys, ok := sys.(ecs.FixedUpdateSystem); ok {
				fixedSys.Update(gameEntities, g.FixedDeltaTime)
			}
		}
		g.Accumulator -= g.FixedDeltaTime
		updateCount++
	}

	if updateCount == maxUpdatesPerFrame {
		g.Accumulator = 0
	}

	if math.IsNaN(g.Accumulator) {
		log.Println("Accumulator стал NaN, сбрасываем в 0")
		g.Accumulator = 0
	}

	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	gameEntities := g.EntityManager.GetAllEntities()
	tileSize := 64

	var interpolationFactor float64
	var camera *components.Camera
	var cameraGeoM ebiten.GeoM

	for _, entity := range gameEntities {
		if entity.HasComponent("Camera") {
			camera = entity.GetComponent("Camera").(*components.Camera)
			camera.ApplyTransform(&cameraGeoM)
			break
		}
	}

	if camera == nil {
		log.Println("Камера не найдена")
		return
	}

	for _, entity := range gameEntities {
		if entity.HasComponent("TileMap") {
			tileMap := entity.GetComponent("TileMap").(*mapping.TileMap)

			cameraX, cameraY := camera.X, camera.Y
			cameraWidth, cameraHeight := camera.Width, camera.Height

			startX := int(math.Max(0, math.Floor(cameraX/float64(tileSize))))
			endX := int(math.Min(float64(tileMap.Width), math.Ceil((cameraX+cameraWidth)/float64(tileSize))))
			startY := int(math.Max(0, math.Floor(cameraY/float64(tileSize))))
			endY := int(math.Min(float64(tileMap.Height), math.Ceil((cameraY+cameraHeight)/float64(tileSize))))

			for y := startY; y < endY; y++ {
				for x := startX; x < endX; x++ {
					tile := tileMap.Tiles[y][x]
					if tile.Sprite != nil {
						op := &ebiten.DrawImageOptions{}
						op.GeoM = cameraGeoM
						op.GeoM.Translate(float64(x*tileSize), float64(y*tileSize))
						screen.DrawImage(tile.Sprite, op)
					}
				}
			}
		}
	}

	if g.FixedDeltaTime <= 0 {
		log.Println("FixedDeltaTime некорректен:", g.FixedDeltaTime)
		interpolationFactor = 0
	} else {
		interpolationFactor = g.Accumulator / g.FixedDeltaTime
		if math.IsNaN(interpolationFactor) {
			log.Println("interpolationFactor стал NaN, Accumulator:", g.Accumulator)
			interpolationFactor = 0
		}
	}

	for _, sys := range g.Systems {
		if drawSys, ok := sys.(interface {
			Draw(*ebiten.Image, []*ecs.Entity, ebiten.GeoM, float64)
		}); ok {
			drawSys.Draw(screen, gameEntities, cameraGeoM, interpolationFactor)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}
