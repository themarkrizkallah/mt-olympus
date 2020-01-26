package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"apollo/env"
)

const dbType = "postgres"

var (
	db       *sql.DB
	AssetIds []string
)

func Init(sslMode string) (*sql.DB, error) {
	var err error

	db = nil

	connStr := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v sslmode=%v",
		env.PostgresHost,
		env.PostgresUser,
		env.PostgresPass,
		env.PostgresDB,
		sslMode,
	)

	db, err = sql.Open(dbType, connStr)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(0)

	for i := uint(0); i < env.RetryTimes; i++ {
		log.Println("database retry #", i+1)

		if err := db.Ping(); err != nil {
			time.Sleep(time.Duration(env.RetrySeconds) * time.Second)
			continue
		}

		return db, nil
	}

	return db, nil
}

func GetDB() *sql.DB {
	return db
}

func GetAssetIds() ([]string, error) {
	const getAssetsSql = `select id from assets`

	if len(AssetIds) > 0 {
		return AssetIds, nil
	}

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	// Get asset ids
	rows, err := db.QueryContext(ctx, getAssetsSql)
	if err != nil {
		log.Println("An error occurred getting assets:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string

		if err := rows.Scan(&id); err != nil {
			log.Println("An error occurred reading assets:", err)
			return nil, err
		}

		AssetIds = append(AssetIds, id)
	}

	return AssetIds, nil
}
