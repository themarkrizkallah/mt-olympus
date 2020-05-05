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
	//config.Net.SASL.User = env.KafkaUser
	//config.Net.SASL.Password = env.KafkaPassword
	config.Producer.RequiredAcks = sarama.WaitForLocal

	for i := uint(0); i < env.RetryTimes; i++ {
		log.Println("async producer retry #", i+1)
		producer, err = sarama.NewAsyncProducer(brokers, config)

		if err != nil {
			time.Sleep(time.Duration(env.RetrySeconds) * time.Second)
			continue
		}

		return &producer, nil
	}

	return nil, err
}

func ProduceMessage(topic string, value []byte) {
	producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}
}
