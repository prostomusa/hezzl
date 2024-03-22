package repository

type Store struct {
	PostgresRepository   *PostgresContainerRepository
	RedisRepository      *RedisRepository
	ClickHouseRepository *ClickHouseRepository
}

func (st *Store) ConfigureStore() {
	st.PostgresRepository = newPostgresContainerRepository()
	st.RedisRepository = newRedisRepository()
	st.ClickHouseRepository = newClickHouseRepository()
}
