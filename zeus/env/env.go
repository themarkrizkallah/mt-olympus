package env

import (
	"os"
	"strconv"
)

var (
	Debug    bool
	Secure   bool
	Port     string
	CertFile string
	Secret   string

	// GRPC variables
	GrpcHostName string
	GrpcPort     string

	// Redis variables
	RedisHostName string
	RedisPort     string
	RedisPassword string

	// MongoDB variables
	MongoHostName     string
	MongoPort         string
	MongoUser         string
	MongoPassword     string
	MongoDb           string
	MongoRetrySeconds uint64
	MongoRetryTimes   uint64
)

func Init() {
	_, Debug = os.LookupEnv("DEBUG")
	_, Secure = os.LookupEnv("SECURE")
	Port = os.Getenv("PORT")
	CertFile = os.Getenv("CERTFILE")
	Secret = os.Getenv("SECRET")

	// GRPC
	GrpcHostName = os.Getenv("GRPC_HOSTNAME")
	GrpcPort = os.Getenv("GRPC_PORT")

	// Redis
	RedisHostName = os.Getenv("REDIS_HOSTNAME")
	RedisPort = os.Getenv("REDIS_PORT")
	RedisPassword = os.Getenv("REDIS_PASSWORD")

	// MongoDB
	MongoHostName = os.Getenv("MONGO_HOSTNAME")
	MongoPort = os.Getenv("MONGO_PORT")
	MongoUser = os.Getenv("MONGO_USER")
	MongoPassword = os.Getenv("MONGO_PASSWORD")
	MongoDb = os.Getenv("MONGO_DB")
	MongoRetrySeconds, _ = strconv.ParseUint(os.Getenv("MONGO_RETRY_SECONDS"), 10, 64)
	MongoRetryTimes, _ = strconv.ParseUint(os.Getenv("MONGO_RETRY_TIMES"), 10, 64)
}
