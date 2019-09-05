package kafka

import (
	"log"
	"time"

	"github.com/Shopify/sarama"

	"matcher/env"
)

var producer sarama.AsyncProducer

func SetupProducer(brokers []string) (*sarama.AsyncProducer, error) {
	var err error

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal

	for i := uint64(0); i < env.KafkaProdRetryTimes; i++ {
		log.Println("async producer retry #", i+1)
		producer, err = sarama.NewAsyncProducer(brokers, config)

		if err != nil {
			time.Sleep(time.Duration(env.KafkaProdRetrySeconds) * time.Second)
			continue
		}

		return &producer, nil
	}

	return nil, err
}
