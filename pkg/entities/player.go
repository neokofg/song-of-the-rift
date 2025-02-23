package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/neokofg/mygame/pkg/components"
	"github.com/neokofg/mygame/pkg/ecs"
	"image/color"
)

type Player struct {
	entity *ecs.Entity
}

func NewPlayer(em *ecs.EntityManager) *Player {
	player := em.CreateEntity()
	player.AddComponent("Position", &components.Position{X: 400, Y: 300})
	player.AddComponent("Velocity", &components.Velocity{X: 5, Y: 5})
	playerColor := color.RGBA{135, 206, 235, 255}
	playerImage := ebiten.NewImage(32, 32)
	playerImage.Fill(playerColor)
	player.AddComponent("Render", &components.Render{Image: playerImage})
	input := components.NewInput(map[ebiten.Key]string{
		ebiten.KeyW: "moveUp",
		ebiten.KeyS: "moveDown",
		ebiten.KeyA: "moveLeft",
		ebiten.KeyD: "moveRight",
	})
	player.AddComponent("Input", input)

	return &Player{
		entity: player,
	}
}

func (p *Player) GetEntity() *ecs.Entity {
	return p.entity
}
