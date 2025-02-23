package systems

import (
	"github.com/neokofg/mygame/pkg/components"
	"github.com/neokofg/mygame/pkg/ecs"
	"math"
)

type VelocitySystem struct {
}

func (vs *VelocitySystem) Update(entities []*ecs.Entity) {
	fixedDeltaTime := 1.0 / 60

	for _, entity := range entities {
		if entity.HasComponent("Input") && entity.HasComponent("Velocity") {
			input := entity.GetComponent("Input").(*components.Input)
			velocity := entity.GetComponent("Velocity").(*components.Velocity)

			if velocity.DashTimer > 0 {
				velocity.DashTimer -= fixedDeltaTime
				if velocity.DashTimer < 0 {
					velocity.DashTimer = 0
				}
			}

			velocity.X = 0
			velocity.Y = 0

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

			if velocity.X != 0 && velocity.Y != 0 {
				factor := velocity.MaxSpeed / math.Sqrt(velocity.X*velocity.X+velocity.Y*velocity.Y)
				velocity.X *= factor
				velocity.Y *= factor
			}

			if input.Actions["sprint"] {
				velocity.X *= 1.5
				velocity.Y *= 1.5
			}

			if input.Actions["dash"] && !input.PreviousActions["dash"] && velocity.DashTimer <= 0 {
				velocity.X *= 10
				velocity.Y *= 10

				velocity.DashTimer = velocity.DashCooldown
			}

			for key, value := range input.Actions {
				input.PreviousActions[key] = value
			}
		}
	}
}
