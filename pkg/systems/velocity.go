package systems

import (
	"github.com/neokofg/mygame/pkg/components"
	"github.com/neokofg/mygame/pkg/ecs"
	"math"
)

type VelocitySystem struct{}

func (vs *VelocitySystem) Update(entities []*ecs.Entity, deltaTime float64) {
	for _, entity := range entities {
		if entity.HasComponent("Input") && entity.HasComponent("Velocity") {
			input := entity.GetComponent("Input").(*components.Input)
			vel := entity.GetComponent("Velocity").(*components.Velocity)

			// Обновляем таймер даша
			if vel.DashTimer > 0 {
				vel.DashTimer -= deltaTime
				if vel.DashTimer < 0 {
					vel.DashTimer = 0
				}
			}

			// Вычисляем направление
			dirX, dirY := 0.0, 0.0
			if input.Actions["moveUp"] {
				dirY = -1.0
			} else if input.Actions["moveDown"] {
				dirY = 1.0
			}
			if input.Actions["moveLeft"] {
				dirX = -1.0
			} else if input.Actions["moveRight"] {
				dirX = 1.0
			}

			// Нормализуем направление
			magnitude := math.Sqrt(dirX*dirX + dirY*dirY)
			if magnitude > 0 {
				dirX /= magnitude
				dirY /= magnitude
			}

			// Базовая скорость
			vel.X = dirX * vel.MaxSpeed
			vel.Y = dirY * vel.MaxSpeed

			// Применяем модификаторы
			currentMaxSpeed := vel.MaxSpeed
			if input.Actions["sprint"] {
				currentMaxSpeed *= 1.5
				vel.X *= 1.5
				vel.Y *= 1.5
			}
			if input.Actions["dash"] && !input.PreviousActions["dash"] && vel.DashTimer <= 0 {
				currentMaxSpeed *= 10
				vel.X *= 10
				vel.Y *= 10
				vel.DashTimer = vel.DashCooldown
			}

			// Ограничиваем скорость
			speed := math.Sqrt(vel.X*vel.X + vel.Y*vel.Y)
			if speed > currentMaxSpeed {
				factor := currentMaxSpeed / speed
				vel.X *= factor
				vel.Y *= factor
			}

			// Применяем deltaTime
			vel.X *= deltaTime
			vel.Y *= deltaTime

			// Обновляем предыдущие действия
			for key, value := range input.Actions {
				input.PreviousActions[key] = value
			}
		}
	}
}
