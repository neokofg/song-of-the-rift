package systems

import (
	"github.com/neokofg/mygame/pkg/components"
	"github.com/neokofg/mygame/pkg/ecs"
	"github.com/neokofg/mygame/pkg/mapping"
	"math"
	"math/rand"
)

type EnemyAISystem struct {
	TileMap *mapping.TileMap
	Player  *ecs.Entity
}

func (eas *EnemyAISystem) Update(entities []*ecs.Entity, deltaTime float64) {
	for _, entity := range entities {
		if entity.HasComponent("AI") && entity.HasComponent("Velocity") && entity.HasComponent("Position") {
			ai := entity.GetComponent("AI").(*components.AI)
			pos := entity.GetComponent("Position").(*components.Position)
			vel := entity.GetComponent("Velocity").(*components.Velocity)

			switch ai.Type {
			case "random":
				eas.handleRandomBehavior(ai, vel)
			case "patrol":
				eas.handlePatrolBehavior(ai, pos, vel, deltaTime)
			case "patrol_chase":
				eas.handlePatrolChaseBehavior(ai, pos, vel, deltaTime)
			case "idle_chase":
				eas.handleIdleChaseBehavior(ai, pos, vel, deltaTime)
			}
		}
	}
}

func (eas *EnemyAISystem) handleRandomBehavior(ai *components.AI, vel *components.Velocity) {
	// Простая реализация: случайное направление
	vel.X = (rand.Float64()*2 - 1) * vel.MaxSpeed
	vel.Y = (rand.Float64()*2 - 1) * vel.MaxSpeed
}

func (eas *EnemyAISystem) handlePatrolBehavior(ai *components.AI, pos *components.Position, vel *components.Velocity, deltaTime float64) {
	if len(ai.PatrolPoints) == 0 {
		return
	}

	// Находим текущую цель
	target := ai.PatrolPoints[ai.CurrentPatrolIdx]

	// Проверяем, достиг ли враг цели
	if math.Abs(float64(target.X*64)-pos.X) < 1 && math.Abs(float64(target.Y*64)-pos.Y) < 1 {
		ai.CurrentPatrolIdx = (ai.CurrentPatrolIdx + 1) % len(ai.PatrolPoints)
		target = ai.PatrolPoints[ai.CurrentPatrolIdx]
	}

	// Обновляем путь, если таймер истек
	ai.PathUpdateTimer -= deltaTime
	if ai.PathUpdateTimer <= 0 {
		start := components.Node{X: int(pos.X / 64), Y: int(pos.Y / 64)}
		goal := target
		ai.Path = AStar(start, goal, eas.TileMap)
		ai.PathUpdateTimer = 1.0 // Обновляем путь раз в секунду
	}

	// Следуем по пути
	if len(ai.Path) > 1 {
		next := ai.Path[1]
		dirX := float64(next.X*64) - pos.X
		dirY := float64(next.Y*64) - pos.Y
		magnitude := math.Sqrt(dirX*dirX + dirY*dirY)
		if magnitude > 0 {
			vel.X = (dirX / magnitude) * vel.MaxSpeed
			vel.Y = (dirY / magnitude) * vel.MaxSpeed
		}
	}
}

func (eas *EnemyAISystem) handlePatrolChaseBehavior(ai *components.AI, pos *components.Position, vel *components.Velocity, deltaTime float64) {
	playerPos := eas.Player.GetComponent("Position").(*components.Position)
	distX := playerPos.X - pos.X
	distY := playerPos.Y - pos.Y
	distance := math.Sqrt(distX*distX + distY*distY)

	const chaseDistance = 500.0 // Расстояние, на котором враг начинает преследование
	if distance < chaseDistance {
		// Преследование игрока
		ai.PathUpdateTimer -= deltaTime
		if ai.PathUpdateTimer <= 0 {
			start := components.Node{X: int(pos.X / 64), Y: int(pos.Y / 64)}
			goal := components.Node{X: int(playerPos.X / 64), Y: int(playerPos.Y / 64)}
			ai.Path = AStar(start, goal, eas.TileMap)
			ai.PathUpdateTimer = 0.5 // Обновляем путь чаще при преследовании
		}

		if len(ai.Path) > 1 {
			next := ai.Path[1]
			dirX := float64(next.X*64) - pos.X
			dirY := float64(next.Y*64) - pos.Y
			magnitude := math.Sqrt(dirX*dirX + dirY*dirY)
			if magnitude > 0 {
				vel.X = (dirX / magnitude) * vel.MaxSpeed
				vel.Y = (dirY / magnitude) * vel.MaxSpeed
			}
		}
	} else {
		// Патрулирование
		eas.handlePatrolBehavior(ai, pos, vel, deltaTime)
	}
}

func (eas *EnemyAISystem) handleIdleChaseBehavior(ai *components.AI, pos *components.Position, vel *components.Velocity, deltaTime float64) {
	playerPos := eas.Player.GetComponent("Position").(*components.Position)
	distX := playerPos.X - pos.X
	distY := playerPos.Y - pos.Y
	distance := math.Sqrt(distX*distX + distY*distY)

	const chaseDistance = 500.0
	const cellSize = 64

	if distance < chaseDistance {
		ai.PathUpdateTimer -= deltaTime
		if ai.PathUpdateTimer <= 0 {
			start := components.Node{X: int(pos.X / cellSize), Y: int(pos.Y / cellSize)}
			goal := components.Node{X: int(playerPos.X / cellSize), Y: int(playerPos.Y / cellSize)}
			ai.Path = AStar(start, goal, eas.TileMap)
			ai.PathUpdateTimer = 0.1
		}

		if len(ai.Path) > 1 {
			next := ai.Path[1]
			targetX := float64(next.X * cellSize)
			targetY := float64(next.Y * cellSize)

			// Проверяем, достиг ли враг текущей точки пути
			if math.Abs(targetX-pos.X) < 1 && math.Abs(targetY-pos.Y) < 1 {
				// Переходим к следующей точке пути
				ai.Path = ai.Path[1:]
				if len(ai.Path) > 1 {
					next = ai.Path[1]
					targetX = float64(next.X * cellSize)
					targetY = float64(next.Y * cellSize)
				} else {
					vel.X = 0
					vel.Y = 0
					return
				}
			}

			dirX := targetX - pos.X
			dirY := targetY - pos.Y
			magnitude := math.Sqrt(dirX*dirX + dirY*dirY)
			if magnitude > 0 {
				vel.X = (dirX / magnitude) * vel.MaxSpeed
				vel.Y = (dirY / magnitude) * vel.MaxSpeed
			}
		} else {
			vel.X = 0
			vel.Y = 0
		}
	} else {
		vel.X = 0
		vel.Y = 0
	}
}
