package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/neokofg/mygame/pkg/components"
	"github.com/neokofg/mygame/pkg/ecs"
)

type RenderSystem struct{}

func (rs *RenderSystem) Draw(screen *ebiten.Image, entities []*ecs.Entity, cameraGeoM ebiten.GeoM) {
	for _, entity := range entities {
		if entity.HasComponent("Position") && entity.HasComponent("Render") {
			pos := entity.GetComponent("Position").(*components.Position)
			render := entity.GetComponent("Render").(*components.Render)
			op := &ebiten.DrawImageOptions{}
			op.GeoM = cameraGeoM
			op.GeoM.Translate(pos.X, pos.Y)
			screen.DrawImage(render.Image, op)
		}
	}
}
