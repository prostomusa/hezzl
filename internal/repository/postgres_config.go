package repository

import (
	"Project/internal/util"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
	"log"
	"os"
	"sync"
)

const (
	OsVarMigrationUrl = "MIGRATION_URL"
	OsVarDatabaseUrl  = "DB_SOURCE"
)

type AbstractPostgresRepository struct {
	db *DbPostgres
	mu sync.RWMutex
}

type PostgresContainerRepository struct {
	GoodRepository *GoodRepository
}

type DbPostgres struct {
	*sql.DB
}

func createPostgresClient() (*DbPostgres, error) {
	host := util.GetEnv("DATABASE_HOST", "localhost")
	port, err := util.GetIntEnv("DATABASE_PORT", 5432)
	if err != nil {
		log.Fatal(err)
	}
	user := util.GetEnv("DATABASE_USERNAME", "postgres")
	password := util.GetEnv("DATABASE_PASSWORD", "postgres")
	dbname := util.GetEnv("DATABASE_DB", "postgres")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	dbPostgres := DbPostgres{DB: db}
	dbPostgres.SetMaxOpenConns(16)
	dbPostgres.SetMaxIdleConns(16)

	return &dbPostgres, nil
}

func (dbPostgres *DbPostgres) CloseDb() {
	err := dbPostgres.Close()
	if err != nil {
		return
	}
}

func newPostgresContainerRepository() *PostgresContainerRepository {
	container := PostgresContainerRepository{}
	dbPostgres, err := createPostgresClient()
	if err != nil {
		log.Fatal(err)
	}
	migrationUrl := os.Getenv(OsVarMigrationUrl)
	dbUrl := os.Getenv(OsVarDatabaseUrl)
	checkFunc := func(variable string, key string) {
		if variable == "" {
			log.Fatalf("Переменная %V не может быть пустой", key)
		}
	}
	checkFunc(migrationUrl, OsVarMigrationUrl)
	checkFunc(dbUrl, OsVarDatabaseUrl)
	runDBMigration(migrationUrl, dbPostgres.DB)
	container.createGoodRepository(dbPostgres)
	return &container
}

func runDBMigration(migrationURL string, db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	migration, err := migrate.NewWithDatabaseInstance(migrationURL, "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}

	log.Println("db migrated successfully")
}
