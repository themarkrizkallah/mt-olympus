package env

import (
	"fmt"
	"os"
)

var (
	Debug        bool
	RetryTimes   uint
	RetrySeconds uint

	OrderBookCap uint64
	Base         string
	Quote        string

	// Postgres
	PostgresUser string
	PostgresPass string
	PostgresDB   string
	PostgresHost string

	// Kafka
	KafkaHost                string
	KafkaPort                string
	KafkaVersion             string
	KafkaUser                string
	KafkaPassword            string
	KafkaConsGroup           string
	KafkaConsReturnErrors    bool
	KafkaConsReturnNotifs    bool
	KafkaProdReturnSuccesses bool
	KafkaProdReturnErrors    bool
)

func Init() {
	_, Debug = os.LookupEnv("DEBUG")
	RetryTimes = 100
	RetrySeconds = 1

	Base = os.Getenv("BASE")
	Quote = os.Getenv("QUOTE")

	// Postgres
	PostgresUser = os.Getenv("POSTGRES_USER")
	PostgresPass = os.Getenv("POSTGRES_PASSWORD")
	PostgresDB = os.Getenv("POSTGRES_DB")
	PostgresHost = os.Getenv("POSTGRES_HOST")

	// Kafka
	KafkaHost = os.Getenv("KAFKA_HOST")
	KafkaPort = os.Getenv("KAFKA_PORT")
	KafkaVersion = os.Getenv("KAFKA_VERSION")
	KafkaUser = os.Getenv("KAFKA_BROKER_USER")
	KafkaPassword = os.Getenv("KAFKA_BROKER_PASSWORD")
	KafkaConsGroup = os.Getenv("KAFKA_CONS_GROUP")
	_, KafkaConsReturnErrors = os.LookupEnv("KAFKA_CONS_RETURN_ERRORS")
	_, KafkaConsReturnNotifs = os.LookupEnv("KAFKA_CONS_RETURN_NOTIFS")
	_, KafkaProdReturnSuccesses = os.LookupEnv("KAFKA_PROD_RETURN_SUCCESSES")
	_, KafkaProdReturnErrors = os.LookupEnv("KAFKA_PROD_RETURN_ERRORS")
}

func GetKafkaBroker() string {
	return fmt.Sprintf("%s:%s", KafkaHost, KafkaPort)
}