package systems

import (
	"github.com/neokofg/mygame/pkg/components"
	"github.com/neokofg/mygame/pkg/ecs"
)

type CameraSystem struct{}

func (cs *CameraSystem) Update(entities []*ecs.Entity) {
	for _, entity := range entities {
		if entity.HasComponent("Camera") && entity.HasComponent("Position") {
			camera := entity.GetComponent("Camera").(*components.Camera)
			position := entity.GetComponent("Position").(*components.Position)

			deadZoneLeft := camera.X + (camera.Width-camera.DeadZoneWidth)/2
			deadZoneRight := deadZoneLeft + camera.DeadZoneWidth
			deadZoneTop := camera.Y + (camera.Height-camera.DeadZoneHeight)/2
			deadZoneBottom := deadZoneTop + camera.DeadZoneHeight

			targetX := camera.X
			targetY := camera.Y

			if position.X < deadZoneLeft {
				targetX = position.X - (camera.Width-camera.DeadZoneWidth)/2
			} else if position.X > deadZoneRight {
				targetX = position.X - (camera.Width+camera.DeadZoneWidth)/2
			}

			if position.Y < deadZoneTop {
				targetY = position.Y - (camera.Height-camera.DeadZoneHeight)/2
			} else if position.Y > deadZoneBottom {
				targetY = position.Y - (camera.Height+camera.DeadZoneHeight)/2
			}

			camera.X += (targetX - camera.X) * camera.Smoothness
			camera.Y += (targetY - camera.Y) * camera.Smoothness
		}
	}
}
