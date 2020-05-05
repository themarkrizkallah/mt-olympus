package database

import (
	"context"
	"log"
	"time"

	"apollo/types"
)

func GetAssets() ([]types.Asset, error) {
	var assets []types.Asset

	const getAssetsSql = `select id, name, tick from assets`

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	// Get asset ids
	rows, err := db.QueryContext(ctx, getAssetsSql)
	if err != nil {
		log.Println("An error occurred getting assets:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var asset types.Asset

		if err := rows.Scan(&asset.Id, &asset.Name, &asset.Tick); err != nil {
			log.Println("An error occurred reading assets:", err)
			return nil, err
		}

		assets = append(assets, asset)
	}

	return assets, nil
}

