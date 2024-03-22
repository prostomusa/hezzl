package manager

import (
	"Project/internal/model/clickhouse"
	"Project/internal/repository"
	"Project/internal/util"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"
	"time"
)

const TopicName = "log"

type NatsManager struct {
	repository *repository.ClickHouseRepository
	NatsUrl    string
}

func newNatsManager(repository *repository.ClickHouseRepository) *NatsManager {
	natsUrl := util.GetEnv("NATS_URL", nats.DefaultURL)
	manager := &NatsManager{
		repository: repository,
		NatsUrl:    natsUrl,
	}
	ticker := time.NewTicker(time.Minute)
	go manager.updateLogsEveryMinute(ticker)
	return manager
}

func (manager *NatsManager) PublishLog(logMessage clickhouse.ClickHouseLog) {
	nc, err := nats.Connect(manager.NatsUrl)
	defer nc.Close()
	marshal, err := msgpack.Marshal(logMessage)
	if err != nil {
		return
	}
	err = nc.Publish(TopicName, marshal)
	if err != nil {
		fmt.Println(err)
	}
}

func (manager *NatsManager) updateLogsEveryMinute(ticker *time.Ticker) {
	for range ticker.C {
		nc, err := nats.Connect(manager.NatsUrl)
		defer nc.Close()
		sub, err := nc.SubscribeSync(TopicName)
		messages, err := sub.Fetch(1000)
		if err == nil {
			messageList := make([]clickhouse.ClickHouseLog, len(messages))
			for i, v := range messages {
				var logMessage clickhouse.ClickHouseLog
				err = msgpack.Unmarshal(v.Data, &logMessage)
				if err == nil {
					messageList[i] = logMessage
				}
			}
			manager.repository.InsertBatchLogs(messageList)
		}
	}
}
