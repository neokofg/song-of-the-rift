package systems

import (
	"github.com/neokofg/mygame/pkg/components"
	"github.com/neokofg/mygame/pkg/ecs"
)

type MovementSystem struct{}

func (ms *MovementSystem) Update(entities []*ecs.Entity) {
	for _, entity := range entities {
		if entity.HasComponent("Position") && entity.HasComponent("Velocity") {
			position := entity.GetComponent("Position").(*components.Position)
			velocity := entity.GetComponent("Velocity").(*components.Velocity)
			fixedDeltaTime := 1.0 / 60

			position.X += velocity.X * fixedDeltaTime
			position.Y += velocity.Y * fixedDeltaTime
		}
	}
}
