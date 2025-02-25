package systems

import (
	"github.com/neokofg/mygame/pkg/components"
	"github.com/neokofg/mygame/pkg/ecs"
	"github.com/neokofg/mygame/pkg/mapping"
	"math"
)

// CollisionSystem обрабатывает столкновения между сущностями и с тайлами.
type CollisionSystem struct{}

// Update обновляет позиции сущностей с учетом столкновений.
func (cs *CollisionSystem) Update(entities []*ecs.Entity) {
	const tileSize = 64 // Размер тайла (предполагается, что тайлы квадратные)

	// Сначала обрабатываем столкновения с тайлами
	for _, entity := range entities {
		if entity.HasComponent("Position") && entity.HasComponent("Collision") {
			pos := entity.GetComponent("Position").(*components.Position)
			coll := entity.GetComponent("Collision").(*components.Collision)
			if coll.Type == components.Solid {
				cs.resolveTileCollisions(entities, pos, coll.Width, coll.Height, tileSize)
			}
		}
	}

	// Затем обрабатываем столкновения между сущностями
	for i, e1 := range entities {
		if e1.HasComponent("Position") && e1.HasComponent("Collision") {
			pos1 := e1.GetComponent("Position").(*components.Position)
			coll1 := e1.GetComponent("Collision").(*components.Collision)
			for j := i + 1; j < len(entities); j++ {
				e2 := entities[j]
				if e2.HasComponent("Position") && e2.HasComponent("Collision") {
					pos2 := e2.GetComponent("Position").(*components.Position)
					coll2 := e2.GetComponent("Collision").(*components.Collision)
					if coll1.Type == components.Solid && coll2.Type == components.Solid {
						cs.resolveEntityCollision(pos1, coll1, pos2, coll2)
					} else if coll1.Type == components.Trigger || coll2.Type == components.Trigger {
						// Обработка триггеров
						if cs.checkOverlap(pos1, coll1, pos2, coll2) {
							// Вызвать событие триггера
							// Например, можно добавить компонент или вызвать метод
						}
					}
				}
			}
		}
	}
}

// resolveTileCollisions корректирует позицию сущности, чтобы она не пересекалась с непроходимыми тайлами.
func (cs *CollisionSystem) resolveTileCollisions(entities []*ecs.Entity, pos *components.Position, width, height, tileSize int) {
	tileMap := cs.getTileMap(entities)
	if tileMap == nil {
		return
	}

	// Находим тайлы, с которыми пересекается bounding box сущности
	minX := int(math.Floor(pos.X / float64(tileSize)))
	maxX := int(math.Floor((pos.X + float64(width-1)) / float64(tileSize)))
	minY := int(math.Floor(pos.Y / float64(tileSize)))
	maxY := int(math.Floor((pos.Y + float64(height-1)) / float64(tileSize)))

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			if x >= 0 && x < tileMap.Width && y >= 0 && y < tileMap.Height {
				tile := tileMap.Tiles[y][x]
				if !tile.Passable {
					// Находим пересечение и корректируем позицию
					cs.pushOutOfTile(pos, width, height, x, y, tileSize)
				}
			}
		}
	}
}

// pushOutOfTile выталкивает сущность из тайла в наименьшем направлении.
func (cs *CollisionSystem) pushOutOfTile(pos *components.Position, width, height, tileX, tileY, tileSize int) {
	entityLeft := pos.X
	entityRight := pos.X + float64(width)
	entityTop := pos.Y
	entityBottom := pos.Y + float64(height)

	tileLeft := float64(tileX * tileSize)
	tileRight := float64((tileX + 1) * tileSize)
	tileTop := float64(tileY * tileSize)
	tileBottom := float64((tileY + 1) * tileSize)

	overlapLeft := entityRight - tileLeft
	overlapRight := tileRight - entityLeft
	overlapTop := entityBottom - tileTop
	overlapBottom := tileBottom - entityTop

	overlaps := []float64{overlapLeft, overlapRight, overlapTop, overlapBottom}
	minOverlap := overlapLeft
	for _, overlap := range overlaps {
		if overlap < minOverlap {
			minOverlap = overlap
		}
	}

	if minOverlap == overlapLeft {
		pos.X -= minOverlap
	} else if minOverlap == overlapRight {
		pos.X += minOverlap
	} else if minOverlap == overlapTop {
		pos.Y -= minOverlap
	} else if minOverlap == overlapBottom {
		pos.Y += minOverlap
	}
}

// resolveEntityCollision корректирует позиции двух сущностей, чтобы они не пересекались.
func (cs *CollisionSystem) resolveEntityCollision(pos1 *components.Position, coll1 *components.Collision, pos2 *components.Position, coll2 *components.Collision) {
	if cs.checkOverlap(pos1, coll1, pos2, coll2) {
		// Находим пересечение
		overlapX, overlapY := cs.getOverlap(pos1, coll1, pos2, coll2)

		// Определяем направление выталкивания
		if math.Abs(overlapX) < math.Abs(overlapY) {
			if pos1.X < pos2.X {
				pos1.X -= overlapX / 2
				pos2.X += overlapX / 2
			} else {
				pos1.X += overlapX / 2
				pos2.X -= overlapX / 2
			}
		} else {
			if pos1.Y < pos2.Y {
				pos1.Y -= overlapY / 2
				pos2.Y += overlapY / 2
			} else {
				pos1.Y += overlapY / 2
				pos2.Y -= overlapY / 2
			}
		}
	}
}

func (cs *CollisionSystem) checkOverlap(pos1 *components.Position, coll1 *components.Collision, pos2 *components.Position, coll2 *components.Collision) bool {
	return pos1.X < pos2.X+float64(coll2.Width) &&
		pos1.X+float64(coll1.Width) > pos2.X &&
		pos1.Y < pos2.Y+float64(coll2.Height) &&
		pos1.Y+float64(coll1.Height) > pos2.Y
}

func (cs *CollisionSystem) getOverlap(pos1 *components.Position, coll1 *components.Collision, pos2 *components.Position, coll2 *components.Collision) (float64, float64) {
	overlapX := math.Min(pos1.X+float64(coll1.Width), pos2.X+float64(coll2.Width)) - math.Max(pos1.X, pos2.X)
	overlapY := math.Min(pos1.Y+float64(coll1.Height), pos2.Y+float64(coll2.Height)) - math.Max(pos1.Y, pos2.Y)
	if pos1.X < pos2.X {
		overlapX = -overlapX
	}
	if pos1.Y < pos2.Y {
		overlapY = -overlapY
	}
	return overlapX, overlapY
}

func (cs *CollisionSystem) getTileMap(entities []*ecs.Entity) *mapping.TileMap {
	for _, e := range entities {
		if e.HasComponent("TileMap") {
			return e.GetComponent("TileMap").(*mapping.TileMap)
		}
	}
	return nil
}
