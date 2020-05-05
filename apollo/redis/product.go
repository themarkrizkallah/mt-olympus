package redis

import (
	"context"
	"time"

	"github.com/go-redis/cache/v7"

	"apollo/database"
	"apollo/types"
)

const productsKey = "products"

func GetProductsList(ctx context.Context) ([]types.Product, error){
	var (
		products []types.Product = nil
		err error
	)

	if Codec.ExistsContext(ctx, productsKey) {
		err = Codec.GetContext(ctx, productsKey, &products)
	} else {
		products, err = database.GetProducts()
		productsItem := cache.Item{
			Key:        productsKey,
			Object:     products,
			Func: func() (i interface{}, err error) {
				return products, nil
			},
			Expiration: 72*time.Hour,
		}

		err = Codec.Once(&productsItem)
	}

	return products, err
}

func GetProductsMap(ctx context.Context) (map[string]types.Product, error){
	var productsMap map[string]types.Product

	products, err := GetProductsList(ctx)
	if err != nil {
		return productsMap, err
	}

	productsMap = make(map[string]types.Product, len(products))
	for _, product := range products {
		productsMap[product.Id] = product
	}

	return productsMap, err
}