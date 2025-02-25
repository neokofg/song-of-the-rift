package components

import (
	"github.com/neokofg/mygame/pkg/ecs"
)

type AI struct {
	Type             string
	Path             []Node      // Путь для патрулирования или преследования
	Target           *ecs.Entity // Цель для преследования (игрок)
	PatrolPoints     []Node      // Точки патрулирования
	PathUpdateTimer  float64     // Таймер для оптимизации
	CurrentPatrolIdx int         // Индекс текущей точки патрулирования
}
