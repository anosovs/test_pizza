package config

import (
	"os"

	"github.com/joho/godotenv"
)



const (
	EnvLocal = "local"
	EnvProd = "prod"
)

type Config struct {
	Env string
	DBDSN string
	ServerHost string
	ServerPort int
}

func Init() (*Config, error) {
	err := godotenv.Load("./config/.env")
	if err != nil {
		return &Config{}, err
	}

	env := os.Getenv("ENV")
	if env == "" || (env != EnvLocal && env!= EnvProd){
		env = EnvLocal
	}

	dsn := ""
	if env != EnvLocal {
		dsn = os.Getenv("DBDSN")
	}

	server := os.Getenv("SERVER_HOST")

	return &Config{
		Env: env,
		DBDSN: dsn,
		ServerHost: server,		
	}, nil
}