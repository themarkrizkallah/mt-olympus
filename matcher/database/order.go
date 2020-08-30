package database

import (
	"database/sql"
	"github.com/golang/protobuf/ptypes"
	"log"
	pb "matcher/proto"
)

const (
	insertOrderSql = `insert into orders(id, product_id, user_id, amount, price, type, side, status, created_at)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	updateOrderSql = `update orders set status = $1 where id = $2`
)

func InsertOrder(tx *sql.Tx, conf pb.OrderConf, productId string) error {
	var err error = nil

	stmt, err := tx.Prepare(insertOrderSql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	t, _ := ptypes.Timestamp(conf.GetCreatedAt())

	_, err = stmt.Exec(
		conf.GetOrderId(),
		productId,
		conf.GetUserId(),
		conf.GetAmount(),
		conf.GetPrice(),
		conf.GetType().String(),
		conf.GetSide().Number(),
		conf.GetStatus(),
		t,
	)

	return err
}

func UpdateOrderStatus(tx *sql.Tx, id, status string) error {
	var err error = nil

	stmt, err := tx.Prepare(updateOrderSql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	log.Printf("Updating order %v: %v\n", id, status)
	_, err = stmt.Exec(status, id)

	return err
}
