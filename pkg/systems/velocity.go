package systems

import (
	"github.com/neokofg/mygame/pkg/components"
	"github.com/neokofg/mygame/pkg/ecs"
	"math"
)

// VelocitySystem обновляет скорости сущностей на основе ввода.
type VelocitySystem struct{}

// Update обновляет скорости сущностей.
func (vs *VelocitySystem) Update(entities []*ecs.Entity, deltaTime float64) {
	for _, entity := range entities {
		if entity.HasComponent("Input") && entity.HasComponent("Playable") && entity.HasComponent("Velocity") {
			input := entity.GetComponent("Input").(*components.Input)
			vel := entity.GetComponent("Velocity").(*components.Velocity)

			// Обновляем таймер рывка
			if vel.DashTimer > 0 {
				vel.DashTimer -= deltaTime
				if vel.DashTimer < 0 {
					vel.DashTimer = 0
				}
			}

			// Вычисляем направление движения
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

			// Устанавливаем базовую скорость
			vel.X = dirX * vel.MaxSpeed
			vel.Y = dirY * vel.MaxSpeed

			// Применяем модификаторы (спринт и рывок)
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

			// Обновляем предыдущие действия
			for key, value := range input.Actions {
				input.PreviousActions[key] = value
			}
		}
	}
}
