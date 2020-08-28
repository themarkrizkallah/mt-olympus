package engine

import (
	"log"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/avast/retry-go"

	"matcher/env"
)

// Topics
const (
	confTopic   = "order.conf"
	tradesTopic = "trades"
)

type Producer struct {
	producer sarama.AsyncProducer
}

func (p *Producer) sendMessage(topic string, value []byte) {
	p.producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}
}

func newAsyncProducer(brokers string) sarama.AsyncProducer {
	var (
		err      error
		producer sarama.AsyncProducer
	)

	config := sarama.NewConfig()
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForLocal

	brokerList := strings.Split(brokers, ",")

	err = retry.Do(
		func() error {
			if producer, err = sarama.NewAsyncProducer(brokerList, config); err != nil {
				return err
			}
			return nil
		},
		retry.Attempts(env.RetryTimes),
		retry.Delay(time.Second),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("Producer - retry %v, error: %s", n, err)
		}),
	)
	if err != nil {
		log.Panicf("Engine - error creating producer: %s\n", err)
	}

	return producer
}
