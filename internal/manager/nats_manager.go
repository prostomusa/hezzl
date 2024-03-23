package manager

import (
	"Project/internal/model/clickhouse"
	"Project/internal/repository"
	"Project/internal/util"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"
	"log"
	"time"
)

const TopicName = "log"
const BatchSize = 1000

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
	nc, _ := nats.Connect(manager.NatsUrl)
	subscriber, err := nc.SubscribeSync(TopicName)
	if err != nil {
		log.Fatal(err)
	}
	go manager.updateLogsEveryMinute(subscriber, ticker)
	return manager
}

func (manager *NatsManager) PublishLog(logMessage clickhouse.ClickHouseLog) {
	nc, err := nats.Connect(manager.NatsUrl)
	defer nc.Drain()
	marshal, err := msgpack.Marshal(logMessage)
	if err != nil {
		return
	}
	err = nc.Publish(TopicName, marshal)
	if err != nil {
		fmt.Println(err)
	}
}

func (manager *NatsManager) updateLogsEveryMinute(subscriber *nats.Subscription, ticker *time.Ticker) {
	for range ticker.C {
		messageList := make([]clickhouse.ClickHouseLog, 0)
		for {
			msg, err := subscriber.NextMsg(time.Second * 2)
			if err != nil {
				break
			}
			var logMessage clickhouse.ClickHouseLog
			err = msgpack.Unmarshal(msg.Data, &logMessage)
			if err == nil {
				messageList = append(messageList, logMessage)
				msg.Ack()
			}
			if len(messageList) > BatchSize {
				break
			}
		}
		if len(messageList) > 0 {
			manager.repository.InsertBatchLogs(messageList)
		}
	}
}
