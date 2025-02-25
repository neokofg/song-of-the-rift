package systems

import (
	"github.com/neokofg/mygame/pkg/components"
	"github.com/neokofg/mygame/pkg/mapping"
	"math"
)

func manhattanDistance(a, b components.Node) float64 {
	return math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y))
}

func AStar(start, goal components.Node, tileMap *mapping.TileMap) []components.Node {
	// Тип для позиции узла
	type Pos [2]int

	// Инициализация списков и карт
	openSet := []components.Node{start} // Список узлов для проверки
	closedSet := make(map[Pos]bool)     // Список уже проверенных позиций
	cameFrom := make(map[Pos]Pos)       // Карта для восстановления пути

	startPos := Pos{start.X, start.Y}
	goalPos := Pos{goal.X, goal.Y}

	// Устанавливаем начальную позицию (родитель сам себе)
	cameFrom[startPos] = startPos

	for len(openSet) > 0 {
		// Находим узел с наименьшим значением f = g + h
		currentIdx := 0
		for i, node := range openSet {
			if node.G+node.H < openSet[currentIdx].G+openSet[currentIdx].H {
				currentIdx = i
			}
		}
		current := openSet[currentIdx]
		openSet = append(openSet[:currentIdx], openSet[currentIdx+1:]...)

		currentPos := Pos{current.X, current.Y}

		// Если достигли цели, восстанавливаем путь
		if currentPos == goalPos {
			path := []components.Node{current}
			for currentPos != startPos {
				parentPos, exists := cameFrom[currentPos]
				if !exists {
					break
				}
				parent := components.Node{X: parentPos[0], Y: parentPos[1]}
				path = append([]components.Node{parent}, path...)
				currentPos = parentPos
			}
			return path
		}

		// Добавляем текущую позицию в закрытый список
		closedSet[currentPos] = true

		// Список соседних позиций (вправо, влево, вверх, вниз)
		neighbors := []Pos{
			{current.X + 1, current.Y},
			{current.X - 1, current.Y},
			{current.X, current.Y + 1},
			{current.X, current.Y - 1},
		}

		for _, neighborPos := range neighbors {
			// Проверяем границы карты
			if neighborPos[0] < 0 || neighborPos[0] >= tileMap.Width ||
				neighborPos[1] < 0 || neighborPos[1] >= tileMap.Height {
				continue
			}

			// Проверяем проходимость клетки
			if !tileMap.Tiles[neighborPos[1]][neighborPos[0]].Passable {
				continue
			}

			// Пропускаем, если позиция уже в закрытом списке
			if closedSet[neighborPos] {
				continue
			}

			tentativeG := current.G + 1

			// Ищем соседний узел в открытом списке
			neighborInOpen := false
			for i, node := range openSet {
				if node.X == neighborPos[0] && node.Y == neighborPos[1] {
					neighborInOpen = true
					// Если нашли лучший путь, обновляем существующий узел
					if tentativeG < node.G {
						openSet[i].G = tentativeG
						openSet[i].H = manhattanDistance(node, goal)
						cameFrom[Pos{node.X, node.Y}] = currentPos
					}
					break
				}
			}

			// Если узел не в открытом списке, добавляем его
			if !neighborInOpen {
				neighbor := components.Node{
					X: neighborPos[0],
					Y: neighborPos[1],
					G: tentativeG,
					H: manhattanDistance(components.Node{X: neighborPos[0], Y: neighborPos[1]}, goal),
				}
				openSet = append(openSet, neighbor)
				cameFrom[neighborPos] = currentPos
			}
		}
	}

	// Путь не найден
	return nil
}

func contains(set []components.Node, node components.Node) bool {
	for _, n := range set {
		if n.X == node.X && n.Y == node.Y {
			return true
		}
	}
	return false
}
