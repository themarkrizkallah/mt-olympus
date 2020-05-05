package database

import (
	"context"
	"log"
	"time"

	"apollo/types"
)

func GetProduct(base, quote string) (types.Product, error) {
	const getProductSql = `select base_id, quote_id, created_at from products where id = $1`

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	// Get product info
	product := types.Product{Id: base + "-" + quote}
	err := db.QueryRowContext(ctx, getProductSql, product.Id).Scan(
		&product.BaseId,
		&product.QuoteId,
		&product.CreatedAt,
	)
	if err != nil {
		log.Println("An error occurred getting assets:", err)
		return types.Product{}, err
	}

	return product, nil
}

func GetProducts() ([]types.Product, error) {
	var products []types.Product
	const getProductsSql = `select id, base_id, quote_id, base_tick, quote_tick, created_at from products`

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	// Get products
	rows, err := db.QueryContext(ctx, getProductsSql)
	if err != nil {
		log.Println("An error occurred getting assets:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product types.Product

		err := rows.Scan(
			&product.Id,
			&product.BaseId,
			&product.QuoteId,
			&product.BaseTick,
			&product.QuoteTick,
			&product.CreatedAt,
		)
		if err != nil {
			log.Println("An error occurred reading products:", err)
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}
