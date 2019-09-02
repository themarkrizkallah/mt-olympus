package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"

	"matcher/env"
)

var Producer sarama.AsyncProducer

// Create the producer
func CreateProducer(){
	config := sarama.NewConfig()
	config.Producer.Return.Successes = env.KafkaProdReturnSuccesses
	config.Producer.Return.Errors = env.KafkaProdReturnErrors
	config.Producer.RequiredAcks = sarama.WaitForAll

	var (
		prod      sarama.AsyncProducer
		connected bool
		err       error
	)

	brokers := []string{fmt.Sprintf("%v:%v", env.KafkaHost, env.KafkaPort)}

	for i := uint64(0); i < env.KafkaProdRetryTimes; i++ {
		log.Printf("Producer Retry #%v\n", i+1)
		prod, err = sarama.NewAsyncProducer(brokers, config)

		if err != nil {
			time.Sleep(time.Duration(env.KafkaProdRetrySeconds) * time.Second)
		} else {
			connected = true
			break
		}
	}

	if !connected {
		log.Fatal("Unable to connect producer to kafka server")
	}

	log.Println("Sarama producer up and running!...")

	Producer = prod
}
