package redis

import (
	"context"
	"time"

	"github.com/go-redis/cache/v7"

	"apollo/database"
	"apollo/types"
)

const assetsKey = "assets"

func GetAssetList(ctx context.Context) ([]types.Asset, error){
	var (
		assets []types.Asset = nil
		err error
	)

	if Codec.ExistsContext(ctx, assetsKey) {
		err = Codec.GetContext(ctx, assetsKey, &assets)
	} else {
		assets, err = database.GetAssets()
		assetsItem := cache.Item{
			Key:        assetsKey,
			Object:     assets,
			Func: func() (i interface{}, err error) {
				return assets, nil
			},
			Expiration: 72*time.Hour,
		}

		err = Codec.Once(&assetsItem)
	}

	return assets, err
}