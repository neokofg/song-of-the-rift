package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/neokofg/mygame/pkg/components"
	"github.com/neokofg/mygame/pkg/ecs"
)

type InputSystem struct{}

func (is *InputSystem) Update(entities []*ecs.Entity) {
	for _, entity := range entities {
		if entity.HasComponent("Input") {
			input := entity.GetComponent("Input").(*components.Input)

			for action, state := range input.Actions {
				input.PreviousActions[action] = state
			}

			for action := range input.Actions {
				input.Actions[action] = false
			}

			for key, action := range input.KeyMap {
				if ebiten.IsKeyPressed(key) {
					input.Actions[action] = true
				}
			}
		}
	}
}
