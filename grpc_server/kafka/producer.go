package kafka

import (
	"log"
	"time"

	"github.com/Shopify/sarama"

	"grpc_server/env"
)

var Producer sarama.SyncProducer

func CreateSyncProducer(brokers []string) (*sarama.SyncProducer, error){
	var err error

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForLocal

	for i := uint64(0); i < env.KafkaProdRetryTimes; i++ {
		log.Println("sync producer retry #", i+1)
		Producer, err = sarama.NewSyncProducer(brokers, config)

		if err != nil {
			time.Sleep(time.Duration(env.KafkaProdRetrySeconds) * time.Second)
			continue
		}

		return &Producer, nil
	}

	return nil, err
}
