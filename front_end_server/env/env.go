package env

import "os"

var (
	Debug        bool
	Secure       bool
	GrpcHostName string
	GrpcPort     string
	Port         string
	CertFile     string
	KeyFile      string
)

func Init() {
	_, Debug = os.LookupEnv("DEBUG")
	_, Secure = os.LookupEnv("SECURE")
	GrpcHostName = os.Getenv("GRPC_HOSTNAME")
	GrpcPort = os.Getenv("GRPC_PORT")
	Port = os.Getenv("PORT")
	CertFile = os.Getenv("CERTFILE")
	KeyFile = os.Getenv("KEYFILE")
}
