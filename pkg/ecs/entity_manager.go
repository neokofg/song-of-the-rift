package ecs

type EntityManager struct {
	entities map[int]*Entity
	nextID   int
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		entities: make(map[int]*Entity),
		nextID:   1,
	}
}

// CreateEntity создаёт новую сущность и возвращает её
func (em *EntityManager) CreateEntity() *Entity {
	entity := NewEntity(em.nextID)
	em.entities[em.nextID] = entity
	em.nextID++
	return entity
}

// GetEntity возвращает сущность по ID
func (em *EntityManager) GetEntity(id int) *Entity {
	return em.entities[id]
}

// RemoveEntity удаляет сущность по ID
func (em *EntityManager) RemoveEntity(id int) {
	delete(em.entities, id)
}

// GetAllEntities возвращает все сущности
func (em *EntityManager) GetAllEntities() []*Entity {
	entities := make([]*Entity, 0, len(em.entities))
	for _, entity := range em.entities {
		entities = append(entities, entity)
	}
	return entities
}
