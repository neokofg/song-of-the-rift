package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/neokofg/mygame/pkg/components"
	"github.com/neokofg/mygame/pkg/ecs"
	"image/color"
)

func NewPlayer(em *ecs.EntityManager) *ecs.Entity {
	player := em.CreateEntity()
	player.AddComponent("Position", &components.Position{X: 400, Y: 300})
	player.AddComponent("Velocity", &components.Velocity{X: 0, Y: 0, MaxSpeed: 150, DashCooldown: 1, DashTimer: 0})
	playerColor := color.RGBA{135, 206, 235, 255}
	playerImage := ebiten.NewImage(32, 32)
	playerImage.Fill(playerColor)
	player.AddComponent("Render", &components.Render{Image: playerImage})
	input := components.NewInput(map[ebiten.Key]string{
		ebiten.KeyW:     "moveUp",
		ebiten.KeyS:     "moveDown",
		ebiten.KeyA:     "moveLeft",
		ebiten.KeyD:     "moveRight",
		ebiten.KeyShift: "sprint",
		ebiten.KeySpace: "dash",
	})
	player.AddComponent("Input", input)
	player.AddComponent("Camera", &components.Camera{
		X: 400, Y: 300, Width: 800, Height: 600, Zoom: 1,
		DeadZoneWidth: 75, DeadZoneHeight: 75, Smoothness: 0.1,
	})

	return player
}
