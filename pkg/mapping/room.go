package mapping

import (
	"math/rand"
	"time"
)

type Room struct {
	X, Y, Width, Height int
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func roomsIntersect(r1, r2 Room) bool {
	return !(r1.X >= r2.X+r2.Width || r1.X+r1.Width <= r2.X ||
		r1.Y >= r2.Y+r2.Height || r1.Y+r1.Height <= r2.Y)
}

func intersectsWithAny(room Room, rooms []Room) bool {
	for _, r := range rooms {
		if roomsIntersect(room, r) {
			return true
		}
	}
	return false
}
