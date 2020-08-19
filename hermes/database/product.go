package database

import (
	"context"
	"log"
	"time"
)

func GetProductIDs() ([]string, error) {
	var productIDs []string
	const getProductsSql = `select id from products`

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	// Get products
	rows, err := db.QueryContext(ctx, getProductsSql)
	if err != nil {
		log.Println("An error occurred getting assets:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string

		if err := rows.Scan(&id); err != nil {
			log.Println("An error occurred reading products:", err)
			return nil, err
		}

		productIDs = append(productIDs, id)
	}

	return productIDs, nil
}