package env

import (
	"fmt"
	"os"
)

var (
	RetryTimes   uint
	RetrySeconds uint

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
	KafkaProdReturnSuccesses bool
	KafkaProdReturnErrors    bool
)

func Init() {
	RetryTimes = 10
	RetrySeconds = 1

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
	_, KafkaProdReturnSuccesses = os.LookupEnv("KAFKA_PROD_RETURN_SUCCESSES")
	_, KafkaProdReturnErrors = os.LookupEnv("KAFKA_PROD_RETURN_ERRORS")
}

func GetKafkaBroker() string {
	return fmt.Sprintf("%s:%s", KafkaHost, KafkaPort)
}