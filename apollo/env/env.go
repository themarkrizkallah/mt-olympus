package env

import "os"

var (
	RetryTimes   uint
	RetrySeconds uint

	// Postgres
	PostgresUser string
	PostgresPass string
	PostgresDB   string
	PostgresHost string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string

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
	RetryTimes = 100
	RetrySeconds = 1

	// Postgres
	PostgresUser = os.Getenv("POSTGRES_USER")
	PostgresPass = os.Getenv("POSTGRES_PASSWORD")
	PostgresDB = os.Getenv("POSTGRES_DB")
	PostgresHost = os.Getenv("POSTGRES_HOST")

	// Redis
	RedisHost = os.Getenv("REDIS_HOST")
	RedisPort = os.Getenv("REDIS_PORT")
	RedisPassword = os.Getenv("REDIS_PASSWORD")

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
