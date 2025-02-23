package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/neokofg/mygame/pkg/game"
)

func main() {
	g := game.NewGame()
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("My Roguelike")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
