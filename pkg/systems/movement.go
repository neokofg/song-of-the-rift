package systems

import (
	"github.com/neokofg/mygame/pkg/components"
	"github.com/neokofg/mygame/pkg/ecs"
)

type MovementSystem struct{}

func (ms *MovementSystem) Update(entities []*ecs.Entity) {
	for _, entity := range entities {
		if entity.HasComponent("Position") && entity.HasComponent("Velocity") {
			pos := entity.GetComponent("Position").(*components.Position)
			vel := entity.GetComponent("Velocity").(*components.Velocity)

			// Сохраняем предыдущую позицию
			pos.PrevX = pos.X
			pos.PrevY = pos.Y

			// Обновляем позицию
			pos.X += vel.X
			pos.Y += vel.Y
		}
	}
}
