package systems

import (
	"github.com/neokofg/mygame/pkg/components"
	"github.com/neokofg/mygame/pkg/ecs"
	"github.com/neokofg/mygame/pkg/mapping"
	"math"
)

type CollisionSystem struct{}

func (cs *CollisionSystem) Update(entities []*ecs.Entity) {
	const (
		tileSize     = 64
		playerWidth  = 32
		playerHeight = 32
	)

	for _, entity := range entities {
		if entity.HasComponent("Position") && entity.HasComponent("Velocity") {
			pos := entity.GetComponent("Position").(*components.Position)
			vel := entity.GetComponent("Velocity").(*components.Velocity)

			// Целевая позиция
			targetX := pos.X + vel.X
			targetY := pos.Y + vel.Y

			// Проверяем коллизии по X
			if vel.X != 0 {
				newX := targetX
				tileXLeft := int(math.Floor(newX / float64(tileSize)))
				tileXRight := int(math.Floor((newX + float64(playerWidth-1)) / float64(tileSize)))
				tileYTop := int(math.Floor(pos.Y / float64(tileSize)))
				tileYBottom := int(math.Floor((pos.Y + float64(playerHeight-1)) / float64(tileSize)))

				if vel.X > 0 {
					if cs.checkCollision(entities, tileXRight, tileYTop) || cs.checkCollision(entities, tileXRight, tileYBottom) {
						pos.X = float64(tileXRight*tileSize) - float64(playerWidth)
						vel.X = 0
					} else {
						pos.X = newX
					}
				} else if vel.X < 0 {
					if cs.checkCollision(entities, tileXLeft, tileYTop) || cs.checkCollision(entities, tileXLeft, tileYBottom) {
						pos.X = float64((tileXLeft + 1) * tileSize)
						vel.X = 0
					} else {
						pos.X = newX
					}
				}
			}

			// Проверяем коллизии по Y
			if vel.Y != 0 {
				newY := targetY
				tileYTop := int(math.Floor(newY / float64(tileSize)))
				tileYBottom := int(math.Floor((newY + float64(playerHeight-1)) / float64(tileSize)))
				tileXLeft := int(math.Floor(pos.X / float64(tileSize)))
				tileXRight := int(math.Floor((pos.X + float64(playerWidth-1)) / float64(tileSize)))

				if vel.Y > 0 {
					if cs.checkCollision(entities, tileXLeft, tileYBottom) || cs.checkCollision(entities, tileXRight, tileYBottom) {
						pos.Y = float64(tileYBottom*tileSize) - float64(playerHeight)
						vel.Y = 0
					} else {
						pos.Y = newY
					}
				} else if vel.Y < 0 {
					if cs.checkCollision(entities, tileXLeft, tileYTop) || cs.checkCollision(entities, tileXRight, tileYTop) {
						pos.Y = float64((tileYTop + 1) * tileSize)
						vel.Y = 0
					} else {
						pos.Y = newY
					}
				}
			}

			// Выталкиваем игрока, если он оказался внутри стены
			cs.resolveOverlap(entities, pos, playerWidth, playerHeight, tileSize)
		}
	}
}

func (cs *CollisionSystem) checkCollision(entities []*ecs.Entity, tileX, tileY int) bool {
	for _, e := range entities {
		if e.HasComponent("TileMap") {
			tileMap := e.GetComponent("TileMap").(*mapping.TileMap)
			if tileX >= 0 && tileX < tileMap.Width && tileY >= 0 && tileY < tileMap.Height {
				tile := tileMap.Tiles[tileY][tileX]
				return !tile.Passable
			}
		}
	}
	return true // По умолчанию считаем тайл непроходимым, если он вне карты
}

func (cs *CollisionSystem) resolveOverlap(entities []*ecs.Entity, pos *components.Position, playerWidth, playerHeight, tileSize int) {
	tileXLeft := int(math.Floor(pos.X / float64(tileSize)))
	tileXRight := int(math.Floor((pos.X + float64(playerWidth-1)) / float64(tileSize)))
	tileYTop := int(math.Floor(pos.Y / float64(tileSize)))
	tileYBottom := int(math.Floor((pos.Y + float64(playerHeight-1)) / float64(tileSize)))

	// Проверяем пересечение с тайлами
	overlapping := false
	for x := tileXLeft; x <= tileXRight; x++ {
		for y := tileYTop; y <= tileYBottom; y++ {
			if cs.checkCollision(entities, x, y) {
				overlapping = true
				break
			}
		}
	}

	if !overlapping {
		return
	}

	// Направления выталкивания
	directions := [][2]float64{
		{0, -1}, // вверх
		{0, 1},  // вниз
		{-1, 0}, // влево
		{1, 0},  // вправо
	}

	step := float64(tileSize) / 4
	for _, dir := range directions {
		newX, newY := pos.X, pos.Y
		for i := 0; i < 4; i++ {
			newX += dir[0] * step
			newY += dir[1] * step

			tileXLeft = int(math.Floor(newX / float64(tileSize)))
			tileXRight = int(math.Floor((newX + float64(playerWidth-1)) / float64(tileSize)))
			tileYTop = int(math.Floor(newY / float64(tileSize)))
			tileYBottom = int(math.Floor((newY + float64(playerHeight-1)) / float64(tileSize)))

			overlap := false
			for x := tileXLeft; x <= tileXRight; x++ {
				for y := tileYTop; y <= tileYBottom; y++ {
					if cs.checkCollision(entities, x, y) {
						overlap = true
						break
					}
				}
			}

			if !overlap {
				pos.X = newX
				pos.Y = newY
				return
			}
		}
	}
}
