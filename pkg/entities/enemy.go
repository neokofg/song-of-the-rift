package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/neokofg/mygame/pkg/components"
	"github.com/neokofg/mygame/pkg/ecs"
	"image/color"
)

func NewEnemy(em *ecs.EntityManager) *ecs.Entity {
	enemy := em.CreateEntity()

	enemy.AddComponent("Position", &components.Position{X: 0, Y: 0})

	enemySprite := ebiten.NewImage(32, 32)
	enemySprite.Fill(color.RGBA{255, 0, 0, 255})

	enemy.AddComponent("Render", &components.Render{Image: enemySprite})
	enemy.AddComponent("Velocity", &components.Velocity{X: 0, Y: 0, MaxSpeed: 200})
	enemy.AddComponent("Collision", &components.Collision{Width: 32, Height: 32, Type: components.Solid})
	enemy.AddComponent("AI", &components.AI{
		Type:            "idle_chase",
		PathUpdateTimer: 1.0,
	})

	return enemy
}
