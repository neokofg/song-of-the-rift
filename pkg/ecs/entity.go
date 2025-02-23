package ecs

type Entity struct {
	ID         int
	components map[string]interface{}
}

func NewEntity(id int) *Entity {
	return &Entity{
		ID:         id,
		components: make(map[string]interface{}),
	}
}

// AddComponent добавляет компонент к сущности
func (e *Entity) AddComponent(name string, component interface{}) {
	e.components[name] = component
}

// RemoveComponent удаляет компонент по имени
func (e *Entity) RemoveComponent(name string) {
	delete(e.components, name)
}

// GetComponent возвращает компонент по имени
func (e *Entity) GetComponent(name string) interface{} {
	return e.components[name]
}

// HasComponent проверяет, есть ли компонент с данным именем
func (e *Entity) HasComponent(name string) bool {
	_, exists := e.components[name]
	return exists
}
