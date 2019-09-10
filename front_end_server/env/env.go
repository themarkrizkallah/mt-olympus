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
	MongoUri          string
	MongoDb           string
	MongoRetrySeconds uint64
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
	MongoUri = os.Getenv("MONGO_URI")
	MongoDb = os.Getenv("MONGO_DB")
	MongoRetrySeconds, _ = strconv.ParseUint(os.Getenv("MONGO_RETRY_SECONDS"), 10, 64)
}
