package systems

import (
	"github.com/neokofg/mygame/pkg/components"
	"github.com/neokofg/mygame/pkg/ecs"
	"math"
)

type VelocitySystem struct {
}

func (vs *VelocitySystem) Update(entities []*ecs.Entity) {
	for _, entity := range entities {
		if entity.HasComponent("Input") && entity.HasComponent("Velocity") {
			input := entity.GetComponent("Input").(*components.Input)
			velocity := entity.GetComponent("Velocity").(*components.Velocity)

			velocity.X = 0
			velocity.Y = 0

			if velocity.X != 0 && velocity.Y != 0 {
				factor := velocity.MaxSpeed / math.Sqrt(velocity.X*velocity.X+velocity.Y*velocity.Y)
				velocity.X *= factor
				velocity.Y *= factor
			}

			if input.Actions["moveUp"] {
				velocity.Y = -velocity.MaxSpeed
			} else if input.Actions["moveDown"] {
				velocity.Y = velocity.MaxSpeed
			}

			if input.Actions["moveLeft"] {
				velocity.X = -velocity.MaxSpeed
			} else if input.Actions["moveRight"] {
				velocity.X = velocity.MaxSpeed
			}
		}
	}
}
