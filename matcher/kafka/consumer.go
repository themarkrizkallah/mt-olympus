package kafka

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"

	"matcher/engine"
	"matcher/env"
)

var (
	consumer       Consumer
	consumerClient sarama.ConsumerGroup
	topics         = []string{"orders"}
)

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	Ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.Ready)
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
		var order engine.Order

		order.FromProto(message.Value) // Decode the message

		log.Printf("Processing: %+v", order)

		orderbook := engine.GetOrderBook()
		trades := orderbook.Process(order) // Process the order

		if len(trades) > 0 {
			fmt.Printf("Completed Trade(s): %+v\n", trades)
		}

		// Send trades to message queue
		for _, trade := range trades {
			rawTrade := trade.ToProto()
			Producer.Input() <- &sarama.ProducerMessage{
				Topic: "trades",
				Value: sarama.ByteEncoder(rawTrade),
			}
		}

		// Mark the message as processed
		session.MarkMessage(message, "")
	}

	return nil
}

func GetConsumerClient() *sarama.ConsumerGroup{
	return &consumerClient
}

func GetConsumer() *Consumer {
	return &consumer
}

func GetConsumerTopics() []string{
	return topics
}

func CreateConsumer() {
	log.Println("Starting a new Sarama consumer")

	version, err := sarama.ParseKafkaVersion(env.KafkaVersion)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	/*
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()
	config.Version = version
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	//Setup a new Sarama consumer group
	consumer = Consumer{
		Ready: make(chan bool),
	}

	brokers := []string{fmt.Sprintf("%v:%v", env.KafkaHost, env.KafkaPort)}

	consumerClient, err = sarama.NewConsumerGroup(brokers, env.KafkaConsGroup, config)
	if err != nil {
		log.Fatalln("Unable to connect consumer to kafka server")
	}

	log.Println("Done setting up consumer group...")
}
