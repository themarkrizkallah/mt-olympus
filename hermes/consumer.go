package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/avast/retry-go"
	"github.com/golang/protobuf/proto"

	"hermes/env"
	pb "hermes/proto"
)

// Sarama configuration options
const (
	assignor = "roundrobin"
	oldest   = true
	verbose  = false
)

var topicPrefixes = []string{
	"order.request",
	"order.conf",
	"trades",
}

// Parse topic and return the topic prefix and product-id
func parseTopic(topic string) (string, string) {
	elements := strings.Split(topic, ".")

	topicPrefix := strings.Join(elements[:len(elements)-1], ".")
	prodId := elements[len(elements)-1]

	return topicPrefix, prodId
}

func constructTopics(productIDs []string) []string {
	topics := make([]string, len(productIDs)*len(topicPrefixes))
	i := 0

	for _, id := range productIDs {
		for _, prefix := range topicPrefixes {
			topics[i] = fmt.Sprintf("%s.%s", prefix, id)
			i += 1
		}
	}

	return topics
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool

	// OrderRequest channel
	requestChan chan<- pb.OrderRequest

	// OrderConf channel
	confChan chan<- pb.OrderConf

	// TradeMessage channel
	tradeMsgChan chan<- pb.TradeMessage
}

func (consumer *Consumer) run(
	topics []string,
	client sarama.ConsumerGroup,
	wg *sync.WaitGroup,
	ctx context.Context,
	) {
	defer wg.Done()

	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		if err := client.Consume(ctx, topics, consumer); err != nil {
			log.Panicf("Error from consumer: %v", err)
		}
		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			return
		}
		consumer.ready = make(chan bool)
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29

	for message := range claim.Messages() {
		topicPrefix, _ := parseTopic(message.Topic)

		switch topicPrefix {
		case "order.request":
			var request pb.OrderRequest
			if err := proto.Unmarshal(message.Value, &request); err != nil {
				log.Panicf("Consumer - error unmarshalling message: %s", err)
			}
			consumer.requestChan <- request

		case "order.conf":
			var conf pb.OrderConf
			if err := proto.Unmarshal(message.Value, &conf); err != nil {
				log.Panicf("Consumer - error unmarshalling message: %s", err)
			}
			consumer.confChan <- conf

		case "trades":
			var tradeMsg pb.TradeMessage
			if err := proto.Unmarshal(message.Value, &tradeMsg); err != nil {
				log.Panicf("Consumer - error unmarshalling message: %s", err)
			}
			consumer.tradeMsgChan <- tradeMsg

		default:
			log.Printf("Consumer - new topic %s encountered", topicPrefix)
		}

		// Mark the message as processed
		session.MarkMessage(message, "")
	}

	return nil
}

func newConsumerGroup(brokers string) (Consumer, sarama.ConsumerGroup) {
	var (
		consumer Consumer
		client sarama.ConsumerGroup
		err error
	)

	if verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	version, err := sarama.ParseKafkaVersion(env.KafkaVersion)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()
	config.Version = version

	switch assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", assignor)
	}

	if oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	/**
	 * Setup a new Sarama consumer group
	 */
	consumer = Consumer{
		ready: make(chan bool),
	}

	err = retry.Do(
		func() error {
			client, err = sarama.NewConsumerGroup(strings.Split(brokers, ","), env.KafkaConsGroup, config)
			if err != nil {
				return err
			}
			return nil
		},
		retry.Attempts(env.RetryTimes),
		retry.Delay(time.Second),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("Consumer - retry %v, error: %s", n, err)
		}),
	)
	if err != nil {
		log.Printf("Error creating consumer group client: %v\n", err)
	}

	return consumer, client
}

