package types

import (
	//"time"
	//
	pb "apollo/proto"
)

//type Payload struct {
//	Amount uint64  `json:"amount"`
//	Price  uint64  `json:"price"`
//	Side   pb.Side `json:"side"`
//	Type   pb.Type `json:"type"`
//}

// ToOrderRequest converts an Order to an OrderRequest
func (o *Order) ToOrderRequest() pb.OrderRequest {
	return pb.OrderRequest{
		UserId:    o.UserId,
		OrderId:   o.OrderId,
		Amount:    o.Amount,
		Price:     o.Price,
		Side:      o.Side,
		Type:      o.Type,
	}
}

// Parse parses a Payload into an Order
//func (p *Payload) Parse() Order {
//	return Order{
//		Amount:    p.Amount,
//		Price:     p.Price,
//		Side:      p.Side,
//		Type:      p.Type,
//		CreatedAt: time.Now(),
//	}
//}

//func ProtoToConf(data []byte) Conf {
//	var protoConf pb.OrderConf
//
//	if err := proto.Unmarshal(data, &protoConf); err != nil {
//		return Conf{}
//	}
//
//	return Conf{
//		UserId:      protoConf.UserId,
//		OrderId:     protoConf.OrderId,
//		CreatedAt:   protoConf.,
//		ConfirmedAt: time.Time{},
//	}
//}
