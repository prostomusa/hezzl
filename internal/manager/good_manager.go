package manager

import (
	"Project/internal/dto"
	"Project/internal/model/clickhouse"
	goods "Project/internal/model/good"
	"Project/internal/repository"
	"database/sql"
	"errors"
	"time"
)

var NotFound = errors.New("NotFound")
var AlreadyExist = errors.New("AlreadyExist")

type GoodManager struct {
	GoodRepository  *repository.GoodRepository
	RedisRepository *repository.RedisRepository
	NatsManager     *NatsManager
}

func newGoodManager(GoodRepository *repository.GoodRepository, RedisRepository *repository.RedisRepository, natsManager *NatsManager) *GoodManager {
	return &GoodManager{
		GoodRepository:  GoodRepository,
		RedisRepository: RedisRepository,
		NatsManager:     natsManager,
	}
}

func (manager *GoodManager) GetGoodByProjectId(projectId int) (*goods.GoodEntity, error) {
	good, err := manager.RedisRepository.GetGoodByProjectId(projectId)
	if err == nil {
		return good, nil
	}
	good, err = manager.GoodRepository.GetGoodByProjectId(projectId)
	if err != nil {
		return nil, NotFound
	}
	manager.RedisRepository.SetGoodProjectId(good)
	return good, nil
}

func (manager *GoodManager) GetGood(id int, projectId int) (*goods.GoodEntity, error) {
	good, err := manager.RedisRepository.GetGood(id, projectId)
	if err == nil {
		return good, nil
	}
	good, err = manager.GoodRepository.GetGood(id, projectId)
	if err != nil {
		return nil, NotFound
	}
	manager.RedisRepository.SetGood(good)
	return good, nil
}

func (manager *GoodManager) UpdateGood(id int, projectId int, requestBody dto.UpdateGoodRequest) (*dto.Good, error) {
	good, err := manager.GetGood(id, projectId)
	if err != nil {
		return nil, NotFound
	}
	description := sql.NullString{}
	if requestBody.Description == "" {
		description = good.Description
	} else {
		description = sql.NullString{String: requestBody.Description, Valid: true}
	}
	good.Name = requestBody.Name
	good.Description = description
	good, err = manager.GoodRepository.UpdateGood(good)
	if err != nil {
		return nil, err
	}
	manager.RedisRepository.SetGood(good)
	result := dto.Good{
		Id:          good.Id,
		ProjectId:   good.ProjectId,
		Name:        good.Name,
		Description: good.Description.String,
		Priority:    good.Priority,
		Removed:     good.Removed,
		CreatedAt:   good.CreatedAt,
	}
	logMessage := clickhouse.ClickHouseLog{
		Id:          good.Id,
		ProjectId:   good.ProjectId,
		Name:        good.Name,
		Description: good.Description,
		Priority:    good.Priority,
		Removed:     good.Removed,
		EventTime:   time.Now(),
	}
	go manager.NatsManager.PublishLog(logMessage)
	return &result, nil
}

func (manager *GoodManager) CreateGood(projectId int, name string) (*dto.Good, error) {
	good, err := manager.GetGoodByProjectId(projectId)
	if err == nil {
		return nil, AlreadyExist
	}
	good, err = manager.GoodRepository.CreateGood(projectId, name)
	if err != nil {
		return nil, err
	}
	manager.RedisRepository.SetGood(good)
	result := dto.Good{
		Id:          good.Id,
		ProjectId:   good.ProjectId,
		Name:        good.Name,
		Description: good.Description.String,
		Priority:    good.Priority,
		Removed:     good.Removed,
		CreatedAt:   good.CreatedAt,
	}
	return &result, nil
}

func (manager *GoodManager) DeleteGood(id int, projectId int) (*dto.DeleteGoodResponse, error) {
	good, err := manager.GetGood(id, projectId)
	if err != nil {
		return nil, NotFound
	}
	good, err = manager.GoodRepository.DeleteGood(good)
	if err != nil {
		return nil, err
	}
	manager.RedisRepository.SetGood(good)
	result := dto.DeleteGoodResponse{
		Id:        good.Id,
		ProjectId: good.ProjectId,
		Removed:   good.Removed,
	}
	logMessage := clickhouse.ClickHouseLog{
		Id:          good.Id,
		ProjectId:   good.ProjectId,
		Name:        good.Name,
		Description: good.Description,
		Priority:    good.Priority,
		Removed:     good.Removed,
		EventTime:   time.Now(),
	}
	go manager.NatsManager.PublishLog(logMessage)
	return &result, nil
}

func (manager *GoodManager) ReprioritiizeGood(id int, projectId int, newPriority int) (*dto.ReprioritizeGoodResponse, error) {
	goodEntity, err := manager.GetGood(id, projectId)
	if err != nil {
		return nil, NotFound
	}
	goodList, err := manager.GoodRepository.Reprioritize(goodEntity, newPriority)
	if err != nil {
		return nil, err
	}
	result := make([]dto.PriorityResponse, len(goodList))
	for i, val := range goodList {
		result[i] = dto.PriorityResponse{Id: val.Id, Priority: val.Priority}
		manager.RedisRepository.SetGood(&val)
		logMessage := clickhouse.ClickHouseLog{
			Id:          val.Id,
			ProjectId:   val.ProjectId,
			Name:        val.Name,
			Description: val.Description,
			Priority:    val.Priority,
			Removed:     val.Removed,
			EventTime:   time.Now(),
		}
		go manager.NatsManager.PublishLog(logMessage)
	}
	return &dto.ReprioritizeGoodResponse{Priorities: result}, nil
}

func (manager *GoodManager) ListGood(limit int, offset int) (*dto.GetGoodsResponse, error) {
	cacheGoods, err := manager.RedisRepository.GetGoods(limit, offset)
	if err == nil {
		return cacheGoods, nil
	}
	returnedGoods, err := manager.GoodRepository.GetGoods(limit, offset)
	if err != nil {
		return nil, err
	}
	goodList := make([]dto.Good, len(returnedGoods))
	removed := 0
	for i, v := range returnedGoods {
		if v.Removed {
			removed++
		}
		goodList[i] = dto.Good{
			Id:          v.Id,
			ProjectId:   v.ProjectId,
			Name:        v.Name,
			Description: v.Description.String,
			Priority:    v.Priority,
			Removed:     v.Removed,
			CreatedAt:   v.CreatedAt,
		}
	}
	result := dto.GetGoodsResponse{
		Meta: dto.MetaGoodsResponse{
			Total:   len(goodList),
			Removed: removed,
			Limit:   limit,
			Offset:  offset,
		},
		Goods: goodList,
	}
	manager.RedisRepository.SetGoods(result)
	return &result, err
}
