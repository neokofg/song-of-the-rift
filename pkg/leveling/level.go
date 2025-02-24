package leveling

import "github.com/neokofg/mygame/pkg/ecs"

type Level struct {
	ID             int
	CreateEntities func(em *ecs.EntityManager) []*ecs.Entity
}
