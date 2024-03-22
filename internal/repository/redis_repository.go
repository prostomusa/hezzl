package repository

import (
	"Project/internal/dto"
	goods "Project/internal/model/good"
	"Project/internal/util"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
	_ "github.com/vmihailenco/msgpack/v5"
	"strconv"
	"time"
)

var exparation = time.Minute

type RedisRepository struct {
	Client *redis.Client
}

func newRedisRepository() *RedisRepository {
	redisUrl := util.GetEnv("REDIS_URL", "localhost:6379")
	return &RedisRepository{
		Client: redis.NewClient(&redis.Options{
			Addr:     redisUrl,
			Password: "",
			DB:       0,
		}),
	}
}

func (repository *RedisRepository) GetGoodByProjectId(projectId int) (*goods.GoodEntity, error) {
	ctx := context.Background()
	res := repository.Client.Get(ctx, getKeyByProjectId(projectId))
	var root []byte
	err := res.Scan(&root)
	if err != nil {
		return nil, fmt.Errorf("Good с projectId = %v в кэше не найден", projectId)
	}
	var good goods.GoodEntity
	err = msgpack.Unmarshal(root, &good)
	if err != nil {
		return nil, fmt.Errorf("Good с projectId = %v в кэше не найден", projectId)
	}
	return &good, nil
}

func (repository *RedisRepository) GetGood(id int, projectId int) (*goods.GoodEntity, error) {
	ctx := context.Background()
	res := repository.Client.Get(ctx, getKey(id, projectId))
	var root []byte
	err := res.Scan(&root)
	if err != nil {
		return nil, fmt.Errorf("Good с id = %v и projectId = %v в кэше не найден", id, projectId)
	}
	var good goods.GoodEntity
	err = msgpack.Unmarshal(root, &good)
	if err != nil {
		return nil, fmt.Errorf("Good с id = %v и projectId = %v в кэше не найден", id, projectId)
	}
	return &good, nil
}

func (repository *RedisRepository) SetGood(good *goods.GoodEntity) {
	ctx := context.Background()
	marshal, err := msgpack.Marshal(good)
	if err != nil {
		return
	}
	repository.Client.SetEx(ctx, getKey(good.Id, good.ProjectId), marshal, exparation)
}

func (repository *RedisRepository) SetGoodProjectId(good *goods.GoodEntity) {
	ctx := context.Background()
	marshal, err := msgpack.Marshal(good)
	if err != nil {
		return
	}
	repository.Client.SetEx(ctx, getKeyByProjectId(good.ProjectId), marshal, exparation)
}

func (repository *RedisRepository) SetGoods(response dto.GetGoodsResponse) {
	ctx := context.Background()
	marshal, err := msgpack.Marshal(response)
	if err != nil {
		return
	}
	repository.Client.SetEx(ctx, getKeyForList(response.Meta.Limit, response.Meta.Offset), marshal, exparation)
}

func (repository *RedisRepository) GetGoods(limit int, offset int) (*dto.GetGoodsResponse, error) {
	ctx := context.Background()
	res := repository.Client.Get(ctx, getKeyForList(limit, offset))
	var root []byte
	err := res.Scan(&root)
	if err != nil {
		return nil, fmt.Errorf("Good с id = %v и projectId = %v в кэше не найден", limit, offset)
	}
	var response dto.GetGoodsResponse
	err = msgpack.Unmarshal(root, &response)
	if err != nil {
		return nil, fmt.Errorf("Good с id = %v и projectId = %v в кэше не найден", limit, offset)
	}
	return &response, nil
}

func getKey(id int, projectId int) string {
	return strconv.Itoa(id) + "good" + strconv.Itoa(projectId)
}

func getKeyByProjectId(projectId int) string {
	return "good" + strconv.Itoa(projectId)
}

func getKeyForList(limit int, offset int) string {
	return strconv.Itoa(limit) + "goodList" + strconv.Itoa(offset)
}
