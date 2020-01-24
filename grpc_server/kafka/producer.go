package kafka

import (
	"github.com/Shopify/sarama"
	"grpc_server/env"
	"log"
	"time"
)

var (
	AsyncProducer sarama.AsyncProducer
	SyncProducer  sarama.SyncProducer
)

func CreateSyncProducer(brokers []string) (*sarama.SyncProducer, error) {
	var err error

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForLocal

	for i := uint64(0); i < env.KafkaProdRetryTimes; i++ {
		log.Println("sync producer retry #", i+1)
		SyncProducer, err = sarama.NewSyncProducer(brokers, config)

		if err != nil {
			time.Sleep(time.Duration(env.KafkaProdRetrySeconds) * time.Second)
			continue
		}

		return &SyncProducer, nil
	}

	return nil, err
}

func CreateAsyncProducer(brokers []string) (*sarama.AsyncProducer, error) {
	var err error

	config := sarama.NewConfig()
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForLocal

	for i := uint64(0); i < env.KafkaProdRetryTimes; i++ {
		log.Println("async producer retry #", i+1)
		AsyncProducer, err = sarama.NewAsyncProducer(brokers, config)

		if err != nil {
			time.Sleep(time.Duration(env.KafkaProdRetrySeconds) * time.Second)
			continue
		}

		return &AsyncProducer, nil
	}

	return nil, err
}
