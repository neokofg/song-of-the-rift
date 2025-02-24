package leveling

import (
	"github.com/neokofg/mygame/pkg/ecs"
	"log"
)

type LevelManager struct {
	CurrentLevel    int
	Levels          map[int]*Level
	EntityManager   *ecs.EntityManager
	CurrentEntities []*ecs.Entity
}

func (lm *LevelManager) LoadLevel(levelID int) {
	if level, ok := lm.Levels[levelID]; ok {
		if lm.CurrentLevel != 0 {
			lm.UnloadLevel()
		}
		lm.CurrentEntities = level.CreateEntities(lm.EntityManager)
		lm.CurrentLevel = levelID
	} else {
		log.Printf("Уровень %d не найден", levelID)
	}
}

func (lm *LevelManager) UnloadLevel() {
	for _, entity := range lm.CurrentEntities {
		lm.EntityManager.RemoveEntity(entity.ID)
	}
	lm.CurrentEntities = nil
	lm.CurrentLevel = 0
}
