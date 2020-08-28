package kafka

import (
	"apollo/redis"
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"sync"

	"github.com/golang/protobuf/proto"

	"apollo/env"
	pb "apollo/proto"
)

const (
	chanSize = 10
	orderRequestTopic = "order.request"
)

var pipeliner *RequestPipeliner

type SendOp struct {
	id       string
	topic    string
	value    []byte
	receiver chan pb.OrderConf
}

type ReceiveOp struct {
	OrderConf pb.OrderConf
}

type RequestPipeliner struct {
	stateMap map[string]chan pb.OrderConf

	send    chan SendOp
	receive chan ReceiveOp

	consumer Consumer
	client   sarama.ConsumerGroup
	consumerWg *sync.WaitGroup

	producer Producer
}

func NewRequestPipeliner() *RequestPipeliner {
	productIds, err := redis.GetProductIDs(context.TODO())
	if err != nil {
		log.Println("RequestPipeliner - error retrieving product IDs,", err)
	}

	log.Println("RequestPipeliner - setting up consumer group")
	consumer, client := newConsumerGroup(env.GetKafkaBroker(), getConsumptionTopics(productIds))

	log.Println("RequestPipeliner - setting up producer")
	producer := Producer{newAsyncProducer(env.GetKafkaBroker())}

	pipeliner = &RequestPipeliner{
		stateMap: make(map[string]chan pb.OrderConf, chanSize),
		send:     make(chan SendOp, chanSize),
		receive:  make(chan ReceiveOp, chanSize),
		consumer: consumer,
		client:   client,
		consumerWg: &sync.WaitGroup{},
		producer: producer,
	}

	return pipeliner
}

func (rp *RequestPipeliner) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer rp.cleanup()

	rp.consumerWg.Add(1)
	go rp.consumer.run(ctx, rp.client, rp.consumerWg)

	// Await till the consumer has been set up
	<-rp.consumer.ready
	log.Println("RequestPipeliner - consumer up and running!")

	for {
		select {
		case sendOp := <-rp.send:
			rp.stateMap[sendOp.id] = sendOp.receiver
			rp.producer.sendMessage(sendOp.topic, sendOp.value)

		case receiveOp := <-rp.receive:
			rp.stateMap[receiveOp.OrderConf.OrderId] <- receiveOp.OrderConf
			delete(rp.stateMap, receiveOp.OrderConf.UserId)

		case <-ctx.Done():
			break
		}
	}

}

func (rp *RequestPipeliner) cleanup() {
	close(rp.send)
	close(rp.receive)
}

// Safe for concurrent calls
func SendOrderRequest(req pb.OrderRequest) (<-chan pb.OrderConf, error) {
	return pipeliner.sendOrderRequest(req)
}

func (rp *RequestPipeliner) sendOrderRequest(req pb.OrderRequest) (<-chan pb.OrderConf, error) {
	data, err := proto.Marshal(&req)
	if err != nil {
		return nil, err
	}

	sendOp := SendOp{
		id:       req.OrderId,
		topic:    fmt.Sprintf("%s.%s", orderRequestTopic, req.GetProductId()),
		value:    data,
		receiver: make(chan pb.OrderConf, 1),
	}
	rp.send <- sendOp

	return sendOp.receiver, nil
}

func receiveOrderConf(conf pb.OrderConf) {
	receiveOp := ReceiveOp{conf}
	pipeliner.receive <- receiveOp
}
