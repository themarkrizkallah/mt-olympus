package kafka

import (
	"log"
	"sync"

	pb "apollo/proto"
)

type SendOp struct {
	ID       string
	Topic    string
	Value    []byte
	Receiver chan pb.OrderConf
}

type ReceiveOp struct {
	OrderConf pb.OrderConf
}

var (
	Receiver chan ReceiveOp
	Sender   chan SendOp
)

const chanSize = 10

func PipelineRequests(wg *sync.WaitGroup) {
	defer wg.Done()

	log.Println("Request Pipeliner running...")
	Receiver = make(chan ReceiveOp, chanSize)
	Sender = make(chan SendOp, chanSize)

	stateMap := make(map[string]chan pb.OrderConf)

	for {
		select {
		case sendOp := <-Sender:
			stateMap[sendOp.ID] = sendOp.Receiver
			ProduceMessage(sendOp.Topic, sendOp.Value)
		case receiveOp := <-Receiver:
			stateMap[receiveOp.OrderConf.UserId] <- receiveOp.OrderConf
			delete(stateMap, receiveOp.OrderConf.UserId)
		}
	}
}
