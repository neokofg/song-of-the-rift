package components

import "github.com/hajimehoshi/ebiten/v2"

type Input struct {
	KeyMap          map[ebiten.Key]string
	Actions         map[string]bool
	PreviousActions map[string]bool
}

func NewInput(keyMap map[ebiten.Key]string) *Input {
	actionSet := make(map[string]struct{})
	for _, action := range keyMap {
		actionSet[action] = struct{}{}
	}

	actions := make(map[string]bool)
	for action := range actionSet {
		actions[action] = false
	}

	return &Input{
		KeyMap:          keyMap,
		Actions:         actions,
		PreviousActions: make(map[string]bool),
	}
}
