package ecs

type FixedUpdateSystem interface {
	Update(entities []*Entity, deltaTime float64)
}

type VariableUpdateSystem interface {
	Update(entities []*Entity)
}
