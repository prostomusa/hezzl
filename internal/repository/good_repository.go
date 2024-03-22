package repository

import (
	goods "Project/internal/model/good"
	"errors"
	"fmt"
	"log"
)

type GoodRepository struct {
	AbstractPostgresRepository
}

func (pt *PostgresContainerRepository) createGoodRepository(postgresDb *DbPostgres) {
	if pt.GoodRepository == nil {
		pt.GoodRepository = &GoodRepository{
			AbstractPostgresRepository: AbstractPostgresRepository{
				db: postgresDb,
			},
		}
	}
}

func (repository *GoodRepository) GetGood(id int, projectId int) (*goods.GoodEntity, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	good := goods.GoodEntity{}
	err := repository.db.QueryRow("SELECT * FROM goods WHERE id = $1 AND project_id = $2", int64(id), projectId).Scan(&good.Id, &good.ProjectId, &good.Name, &good.Description, &good.Priority, &good.Removed, &good.CreatedAt)
	if err != nil {
		log.Printf("Не удалось получить good c id = %v и projectId = %v", id, projectId)
		return nil, err
	}
	return &good, nil
}

func (repository *GoodRepository) GetGoods(limit int, offset int) ([]goods.GoodEntity, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	var goodList []goods.GoodEntity
	rows, err := repository.db.Query("SELECT * FROM goods LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		log.Print("Не получить список good")
		return nil, err
	}
	for rows.Next() {
		newGood := goods.GoodEntity{}
		errScan := rows.Scan(&newGood.Id, &newGood.ProjectId, &newGood.Name, &newGood.Description, &newGood.Priority, &newGood.Removed, &newGood.CreatedAt)
		if errScan == nil {
			goodList = append(goodList, newGood)
		}
	}
	return goodList, nil
}

func (repository *GoodRepository) CreateGood(projectId int, name string) (*goods.GoodEntity, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	if name == "" {
		errText := "Поле Name пустое"
		log.Print(errText)
		return nil, fmt.Errorf(errText)
	}

	tx, err := repository.db.Begin()
	if err != nil {
		log.Print("Не удалось начать транзакцию")
		return nil, err
	}
	defer tx.Rollback()
	good := goods.GoodEntity{}
	err = tx.QueryRow("INSERT INTO goods(project_id, name) values ($1, $2) RETURNING *",
		projectId, name).Scan(&good.Id, &good.ProjectId, &good.Name, &good.Description, &good.Priority, &good.Removed, &good.CreatedAt)
	if err != nil {
		log.Printf("Не удалось создать good c projectId = %v", projectId)
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &good, nil
}

func (repository *GoodRepository) UpdateGood(good *goods.GoodEntity) (*goods.GoodEntity, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	if err := checkFields(good); err != nil {
		return nil, err
	}

	tx, err := repository.db.Begin()
	if err != nil {
		log.Print("Не удалось начать транзакцию")
		return nil, err
	}
	defer tx.Rollback()

	newGood := goods.GoodEntity{}
	err = tx.QueryRow("UPDATE goods SET name = $1, description = $2 WHERE id = $3 AND project_id = $4 RETURNING *",
		good.Name, good.Description, good.Id, good.ProjectId).Scan(&newGood.Id, &newGood.ProjectId, &newGood.Name, &newGood.Description, &newGood.Priority, &newGood.Removed, &newGood.CreatedAt)
	if err != nil {
		log.Printf("Не удалось обновить good с id = %v и project_id = %v", good.Id, good.ProjectId)
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &newGood, nil
}

func (repository *GoodRepository) DeleteGood(good *goods.GoodEntity) (*goods.GoodEntity, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	if err := checkFields(good); err != nil {
		return nil, err
	}

	tx, err := repository.db.Begin()
	if err != nil {
		log.Print("Не удалось начать транзакцию")
		return nil, err
	}
	defer tx.Rollback()
	newGood := goods.GoodEntity{}
	err = tx.QueryRow("UPDATE goods SET removed = true where id = $1 and project_id = $2 RETURNING *",
		good.Id, good.ProjectId).Scan(&newGood.Id, &newGood.ProjectId, &newGood.Name, &newGood.Description, &newGood.Priority, &newGood.Removed, &newGood.CreatedAt)
	if err != nil {
		log.Printf("Не удалось удалить good c id = %v и projectId = %v", good.Id, good.ProjectId)
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &newGood, nil
}

func (repository *GoodRepository) Reprioritize(good *goods.GoodEntity, newPriority int) ([]goods.GoodEntity, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	if err := checkFields(good); err != nil {
		return nil, err
	}

	tx, err := repository.db.Begin()
	if err != nil {
		log.Print("Не удалось начать транзакцию")
		return nil, err
	}
	defer tx.Rollback()

	newGood := goods.GoodEntity{}
	err = tx.QueryRow("UPDATE goods SET priority = $1 WHERE id = $2 AND project_id = $3 returning *",
		newPriority, good.Id, good.ProjectId).Scan(&newGood.Id, &newGood.ProjectId, &newGood.Name, &newGood.Description, &newGood.Priority, &newGood.Removed, &newGood.CreatedAt)
	if err != nil {
		log.Printf("Не удалось обновить приоритет good c id = %v и projectId = %v", good.Id, good.ProjectId)
		return nil, err
	}
	var goodList []goods.GoodEntity
	goodList = append(goodList, newGood)
	rows, err := tx.Query("UPDATE goods SET priority = priority + 1 WHERE priority >= $1 and id != $2 returning *", newPriority, good.Id)
	if err != nil {
		log.Print("Не удалось обновить good")
		return nil, err
	}
	for rows.Next() {
		newGood := goods.GoodEntity{}
		errScan := rows.Scan(&newGood.Id, &newGood.ProjectId, &newGood.Name, &newGood.Description, &newGood.Priority, &newGood.Removed, &newGood.CreatedAt)
		if errScan == nil {
			goodList = append(goodList, newGood)
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return goodList, nil
}

func checkFields(good *goods.GoodEntity) error {
	if good.Name == "" || good.Id == 0 || good.ProjectId == 0 || good.CreatedAt == nil || good.Priority == 0 {
		errText := "Не все поля заполнены"
		log.Print(errText)
		return errors.New(errText)
	}
	return nil
}
