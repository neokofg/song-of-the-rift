package mapping

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/neokofg/mygame/pkg/ecs"
	"image/color"
	"log"
	"math/rand"
)

type TileMap struct {
	Width  int
	Height int
	Tiles  [][]TileData
}

func CreateTileMapEntity(em *ecs.EntityManager, width, height int) *ecs.Entity {
	tileMap, spawnPoint := GenerateTileMap(width, height)
	entity := em.CreateEntity()
	entity.AddComponent("TileMap", tileMap)
	entity.AddComponent("SpawnPoint", spawnPoint)
	return entity
}

func CreateEmptyTileMap(width, height int) *TileMap {
	tileMap := &TileMap{
		Width:  width,
		Height: height,
		Tiles:  make([][]TileData, height),
	}
	wallTile := createGrayTile(64, 64)
	for y := 0; y < height; y++ {
		tileMap.Tiles[y] = make([]TileData, width)
		for x := 0; x < width; x++ {
			tileMap.Tiles[y][x] = TileData{
				Type:     "wall",
				Passable: false,
				Sprite:   wallTile,
			}
		}
	}
	return tileMap
}

func GenerateTileMap(width, height int) (*TileMap, struct{ X, Y int }) {
	tileMap := CreateEmptyTileMap(width, height)
	floorTile, _, err := ebitenutil.NewImageFromFile("pkg/assets/sprites/floor.png")
	if err != nil {
		log.Fatal(err)
	}

	// Генерация главного пути и комнат
	rooms, spawnPoint := generateMainPathWithRooms(tileMap, floorTile)

	// Генерация развилок
	generateBranches(tileMap, rooms, floorTile)

	// Добавление естественности
	addNaturalNoise(tileMap, floorTile)

	return tileMap, spawnPoint
}

func generateMainPathWithRooms(tileMap *TileMap, floorTile *ebiten.Image) ([]Room, struct{ X, Y int }) {
	startX, startY := 1, tileMap.Height/2 // Начало слева по центру
	endX := tileMap.Width - 2             // Конец справа
	x, y := startX, startY
	rooms := []Room{}
	mainPathWidth := 2 // Уменьшаем ширину пути, чтобы комнаты выделялись
	spawnPoint := struct{ X, Y int }{x, y}

	// Параметры комнат
	numRooms := rand.Intn(3) + 4 // 4-6 комнат
	stepsBetweenRooms := tileMap.Width / (numRooms + 1)

	for steps := 0; x < endX; steps++ {
		// Плавное движение
		if rand.Float32() < 0.85 {
			x++
		}
		if rand.Float32() < 0.3 && y > mainPathWidth+2 && y < tileMap.Height-mainPathWidth-2 {
			y += rand.Intn(3) - 1
		}

		// Заполняем путь (меньшая ширина)
		for dy := -mainPathWidth; dy <= mainPathWidth; dy++ {
			for dx := -mainPathWidth; dx <= mainPathWidth; dx++ {
				nx, ny := x+dx, y+dy
				if nx >= 0 && nx < tileMap.Width && ny >= 0 && ny < tileMap.Height {
					tileMap.Tiles[ny][nx] = TileData{"floor", true, floorTile}
				}
			}
		}

		// Добавляем комнату
		if len(rooms) < numRooms && steps >= stepsBetweenRooms*(len(rooms)+1) {
			roomWidth := rand.Intn(10) + 8 // 8-17 тайлов, больше для заметности
			roomHeight := rand.Intn(8) + 6 // 6-13 тайлов
			roomX := x - roomWidth/2
			roomY := y - roomHeight/2
			if roomX < 1 || roomX+roomWidth >= tileMap.Width-1 || roomY < 1 || roomY+roomHeight >= tileMap.Height-1 {
				continue
			}

			newRoom := Room{roomX, roomY, roomWidth, roomHeight}
			rooms = append(rooms, newRoom)
			fillRoom(tileMap, newRoom, floorTile)
			x += roomWidth / 2 // Перепрыгиваем через комнату
		}
	}

	// Дотягиваем путь до конца
	for x < endX {
		x++
		for dy := -mainPathWidth; dy <= mainPathWidth; dy++ {
			nx, ny := x, y+dy
			if nx >= 0 && nx < tileMap.Width && ny >= 0 && ny < tileMap.Height {
				tileMap.Tiles[ny][nx] = TileData{"floor", true, floorTile}
			}
		}
	}

	return rooms, spawnPoint
}

func generateBranches(tileMap *TileMap, rooms []Room, floorTile *ebiten.Image) {
	directions := [8][2]int{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
		{1, 1}, {-1, -1}, {1, -1}, {-1, 1},
	}

	for _, room := range rooms {
		numBranches := rand.Intn(2) + 2 // 2-3 развилки для каждой комнаты
		for i := 0; i < numBranches; i++ {
			var x, y int
			switch rand.Intn(4) {
			case 0: // Верх
				x = room.X + rand.Intn(room.Width)
				y = room.Y
			case 1: // Низ
				x = room.X + rand.Intn(room.Width)
				y = room.Y + room.Height - 1
			case 2: // Лево
				x = room.X
				y = room.Y + rand.Intn(room.Height)
			case 3: // Право
				x = room.X + room.Width - 1
				y = room.Y + rand.Intn(room.Height)
			}

			// Генерация длинной развилки
			branchLength := rand.Intn(15) + 10 // 10-25 шагов, чтобы было заметно
			branchWidth := 1                   // Узкая развилка для контраста
			for step := 0; step < branchLength; step++ {
				dir := directions[rand.Intn(len(directions))]
				x += dir[0] * (1 + rand.Intn(2))
				y += dir[1] * (1 + rand.Intn(2))

				if x < 1 || x >= tileMap.Width-1 || y < 1 || y >= tileMap.Height-1 {
					break
				}

				for dy := -branchWidth; dy <= branchWidth; dy++ {
					for dx := -branchWidth; dx <= branchWidth; dx++ {
						nx, ny := x+dx, y+dy
						if nx >= 0 && nx < tileMap.Width && ny >= 0 && ny < tileMap.Height {
							tileMap.Tiles[ny][nx] = TileData{"floor", true, floorTile}
						}
					}
				}
			}

			// Зона в конце развилки
			zoneWidth := rand.Intn(6) + 4  // 4-9 тайлов
			zoneHeight := rand.Intn(5) + 3 // 3-7 тайлов
			for dy := -zoneHeight / 2; dy <= zoneHeight/2; dy++ {
				for dx := -zoneWidth / 2; dx <= zoneWidth/2; dx++ {
					nx, ny := x+dx, y+dy
					if nx >= 0 && nx < tileMap.Width && ny >= 0 && ny < tileMap.Height {
						tileMap.Tiles[ny][nx] = TileData{"floor", true, floorTile}
					}
				}
			}
		}
	}
}

func addNaturalNoise(tileMap *TileMap, floorTile *ebiten.Image) {
	for y := 1; y < tileMap.Height-1; y++ {
		for x := 1; x < tileMap.Width-1; x++ {
			if tileMap.Tiles[y][x].Type == "wall" && rand.Float32() < 0.15 {
				hasFloorNearby := false
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						if tileMap.Tiles[y+dy][x+dx].Type == "floor" {
							hasFloorNearby = true
							break
						}
					}
				}
				if hasFloorNearby {
					tileMap.Tiles[y][x] = TileData{"floor", true, floorTile}
				}
			}
		}
	}
}

func fillRoom(tileMap *TileMap, room Room, floorTile *ebiten.Image) {
	for y := room.Y; y < room.Y+room.Height; y++ {
		for x := room.X; x < room.X+room.Width; x++ {
			if x >= 0 && x < tileMap.Width && y >= 0 && y < tileMap.Height {
				tileMap.Tiles[y][x] = TileData{"floor", true, floorTile}
			}
		}
	}
}

func createGrayTile(width, height int) *ebiten.Image {
	img := ebiten.NewImage(width, height)
	img.Fill(color.Gray{128})
	return img
}
