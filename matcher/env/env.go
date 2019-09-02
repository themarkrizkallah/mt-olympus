package env

import (
	"log"
	"os"
	"strconv"
)

var (
	Debug        bool
	OrderBookCap uint64

	KafkaHost    string
	KafkaPort    string
	KafkaVersion string

	// Kafka Consumer Vars
	KafkaConsGroup        string
	KafkaConsReturnErrors bool
	KafkaConsReturnNotifs bool
	KafkaConsRetrySeconds uint64
	KafkaConsRetryTimes   uint64

	// Kafka Producer Vars
	KafkaProdReturnSuccesses bool
	KafkaProdReturnErrors    bool
	KafkaProdRetrySeconds    uint64
	KafkaProdRetryTimes      uint64
)

func Init() {
	var err error

	_, Debug = os.LookupEnv("DEBUG")
	OrderBookCap, err = strconv.ParseUint(os.Getenv("ORDERBOOK_CAP"), 10, 64)
	if err != nil {
		log.Println("Error reading ORDERBOOK_CAP, defaulting to 100")
		OrderBookCap = 100
	}

	KafkaHost = os.Getenv("KAFKA_HOST")
	KafkaPort = os.Getenv("KAFKA_PORT")
	KafkaVersion = os.Getenv("KAFKA_VERSION")

	// Kafka Consumer Vars
	KafkaConsGroup = os.Getenv("KAFKA_CONS_GROUP")
	_, KafkaConsReturnErrors = os.LookupEnv("KAFKA_CONS_RETURN_ERRORS")
	_, KafkaConsReturnNotifs = os.LookupEnv("KAFKA_CONS_RETURN_NOTIFS")

	KafkaConsRetrySeconds, err = strconv.ParseUint(os.Getenv("KAFKA_CONS_RETRY_SECONDS"), 10, 64)
	if err != nil {
		log.Println("Error reading KAFKA_CONS_RETRY_SECONDS, defaulting to 10")
		KafkaConsRetrySeconds = 10
	}

	KafkaConsRetryTimes, err = strconv.ParseUint(os.Getenv("KAFKA_CONS_RETRY_TIMES"), 10, 64)
	if err != nil {
		log.Println("Error reading KAFKA_CONS_RETRY_TIMES, defaulting to 100")
		KafkaConsRetryTimes = 100
	}

	// Kafka Producer Vars
	_, KafkaProdReturnSuccesses = os.LookupEnv("KAFKA_PROD_RETURN_SUCCESSES")
	_, KafkaProdReturnErrors = os.LookupEnv("KAFKA_PROD_RETURN_ERRORS")

	KafkaProdRetrySeconds, err = strconv.ParseUint(os.Getenv("KAFKA_PROD_RETRY_SECONDS"), 10, 64)
	if err != nil {
		log.Println("Error reading KAFKA_PROD_RETRY_SECONDS, defaulting to 10")
		KafkaProdRetrySeconds = 10
	}

	KafkaProdRetryTimes, err = strconv.ParseUint(os.Getenv("KAFKA_PROD_RETRY_TIMES"), 10, 64)
	if err != nil {
		log.Println("Error reading KAFKA_PROD_RETRY_TIMES, defaulting to 100")
		KafkaProdRetryTimes = 100
	}
}