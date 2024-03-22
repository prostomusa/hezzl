package manager

import "Project/internal/repository"

type Manager struct {
	GoodManager *GoodManager
	NatsManager *NatsManager
}

func NewManager(store *repository.Store) *Manager {
	natsManager := newNatsManager(store.ClickHouseRepository)
	return &Manager{
		NatsManager: natsManager,
		GoodManager: newGoodManager(store.PostgresRepository.GoodRepository, store.RedisRepository, natsManager),
	}
}
