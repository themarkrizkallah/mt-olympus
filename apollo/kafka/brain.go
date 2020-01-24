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

func PipelineRequests(wg *sync.WaitGroup) {
	defer wg.Done()

	log.Println("PipelineRequests thread running...")
	Receiver = make(chan ReceiveOp)
	Sender = make(chan SendOp)

	stateMap := make(map[string]chan pb.OrderConf)

	for {
		select {
		case sendOp := <-Sender:
			//log.Println("Received sendOp, producing message")
			stateMap[sendOp.ID] = sendOp.Receiver
			ProduceMessage(sendOp.Topic, sendOp.Value)
		case receiveOp := <-Receiver:
			//log.Println("Received receiveOp, passing message")
			stateMap[receiveOp.OrderConf.UserId] <- receiveOp.OrderConf
			delete(stateMap, receiveOp.OrderConf.UserId)
		}
	}
}
