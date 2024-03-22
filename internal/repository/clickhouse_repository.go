package repository

import (
	model "Project/internal/model/clickhouse"
	"Project/internal/util"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/mailru/go-clickhouse/v2"
	"log"
	"sync"
)

type ClickHouseRepository struct {
	Db *sql.DB
	mu sync.Mutex
}

func newClickHouseRepository() *ClickHouseRepository {
	db, err := connect()
	if err != nil {
		log.Fatalf("Нет подключения к базе ClickHouse")
	}
	rep := &ClickHouseRepository{
		Db: db,
	}
	return rep
}
func connect() (*sql.DB, error) {
	ClickHouseDbUrl := util.GetEnv("CLICKHOUSE_DB_URL", "http://localhost:8123/default")
	ClickHouseMigrationUrl := util.GetEnv("CLICKHOUSE_MIGRATION_URL", "file://database/migration/clickhouse")
	db, err := sql.Open("chhttp", ClickHouseDbUrl)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	runClickHouseMigration(ClickHouseMigrationUrl, db)
	return db, nil
}

func runClickHouseMigration(migrationURL string, db *sql.DB) {
	driver, err := clickhouse.WithInstance(db, &clickhouse.Config{})
	migration, err := migrate.NewWithDatabaseInstance(migrationURL, "clickhouse", driver)
	if err != nil {
		log.Fatal(err)
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}

	log.Println("db clickhouse migrated successfully")
}

func (repository *ClickHouseRepository) InsertBatchLogs(logs []model.ClickHouseLog) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	tx, err := repository.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("INSERT INTO logs (*) VALUES (?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, logItem := range logs {
		if _, err := stmt.Exec(
			logItem.Id,
			logItem.ProjectId,
			logItem.Name,
			logItem.Description,
			logItem.Priority,
			logItem.Removed,
			logItem.EventTime,
		); err != nil {
			fmt.Println(err)
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
	}
	return nil
}
