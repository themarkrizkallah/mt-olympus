package database

import (
	"database/sql"
	"log"

	"matcher/types"
)

const (
	insertOrderSql = `insert into orders(id, product_id, user_id, amount, price, type, side, status, created_at)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	updateOrderSql = `update orders set status = $1 where id = $2`
)

func InsertOrder(tx *sql.Tx, order types.Order, status, productId string) error {
	var err error = nil

	stmt, err := tx.Prepare(insertOrderSql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		order.OrderId,
		productId,
		order.UserId,
		order.Amount,
		order.Price,
		order.Type.String(),
		order.Side.Number(),
		status,
		order.CreatedAt,
	)

	return err
}

func UpdateOrderStatus(tx *sql.Tx, orderUpdate types.OrderUpdate) error {
	var err error = nil

	stmt, err := tx.Prepare(updateOrderSql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	log.Printf("Updating order %v: %v\n", orderUpdate.OrderId, orderUpdate.Status)
	_, err = stmt.Exec(orderUpdate.Status, orderUpdate.OrderId)

	return err
}
