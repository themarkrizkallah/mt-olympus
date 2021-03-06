package database

import (
	"context"
	"log"
	"time"
)

type Product struct {
	Id        string
	BaseId    string
	QuoteId   string
	CreatedAt time.Time
}

func GetProduct(base, quote string) (Product, error) {
	const getProductSql = `select base_id, quote_id, created_at from products where id = $1`

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	// Get product info
	product := Product{Id: base + "-" + quote}
	err := db.QueryRowContext(ctx, getProductSql, product.Id).Scan(
		&product.BaseId,
		&product.QuoteId,
		&product.CreatedAt,
	)
	if err != nil {
		log.Println("An error occurred getting assets:", err)
		return Product{}, err
	}

	return product, nil
}
