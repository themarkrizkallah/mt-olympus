package kafka

import (
	"log"
	"time"

	"github.com/Shopify/sarama"

	"apollo/env"
)

var AsyncProducer sarama.AsyncProducer

func CreateAsyncProducer(brokers []string) (*sarama.AsyncProducer, error) {
	var err error

	config := sarama.NewConfig()
	//config.Net.SASL.User = env.KafkaUser
	//config.Net.SASL.Password = env.KafkaPassword
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForLocal

	for i := uint(0); i < env.RetryTimes; i++ {
		log.Println("async producer retry #", i+1)

		if AsyncProducer, err = sarama.NewAsyncProducer(brokers, config); err != nil {
			time.Sleep(time.Duration(env.RetrySeconds) * time.Second)
		} else {
			break
		}
	}

	return &AsyncProducer, err
}

func ProduceMessage(topic string, value []byte) {
	AsyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}
}
