package components

type Node struct {
	X, Y   int     // Координаты узла
	G, H   float64 // G - стоимость от старта, H - эвристика до цели
	Parent *Node   // Родительский узел для восстановления пути
}
