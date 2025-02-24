package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/neokofg/mygame/pkg/game"
	"log"
)

func main() {
	g := game.NewGame()
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("mygame")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
