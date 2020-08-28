package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/avast/retry-go"
	_ "github.com/lib/pq"

	"apollo/env"
)

const dbType = "postgres"

var db *sql.DB

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

	err = retry.Do(
		func() error {
			if err := db.Ping(); err != nil {
				return err
			}
			return nil
		},
		retry.Attempts(env.RetryTimes),
		retry.Delay(time.Second),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("Database - retry %v, error: %s", n, err)
		}),
	)
	if err != nil {
		log.Println("Error creating consumer group client:\n", err)
	}

	return db, err
}

func GetDB() *sql.DB {
	return db
}
