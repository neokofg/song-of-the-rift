package systems

import (
	"github.com/neokofg/mygame/pkg/components"
	"github.com/neokofg/mygame/pkg/ecs"
)

type MovementSystem struct{}

func (ms *MovementSystem) Update(entities []*ecs.Entity) {
	for _, entity := range entities {
		if entity.HasComponent("Input") && entity.HasComponent("Velocity") {
			input := entity.GetComponent("Input").(*components.Input)
			position := entity.GetComponent("Position").(*components.Position)
			velocity := entity.GetComponent("Velocity").(*components.Velocity)

			if input.Actions["moveUp"] {
				position.Y = position.Y - velocity.Y
			} else if input.Actions["moveDown"] {
				position.Y = position.Y + velocity.Y
			}

			if input.Actions["moveLeft"] {
				position.X = position.X - velocity.X
			} else if input.Actions["moveRight"] {
				position.X = position.X + velocity.X
			}
		}
	}
}
