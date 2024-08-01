package config

import (
	"os"
	"sync"

	"messages_handler/pkg/logging"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Host string `yaml:"host" env-default:"localhost"`
		Port string `yaml:"port" env-default:"6969"`
	}
}

var (
	instance *Config
	once     sync.Once
)

func LoadConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger("warn")
		logger.Info("Read aplication config")

		instance = &Config{}

		if err := godotenv.Load(); err != nil {
			logger.Fatal("Error loading .env file")
		}

		instance.Server.Host = os.Getenv("SERVER_HOST")
		instance.Server.Port = os.Getenv("SERVER_PORT")
	})
	
	return instance

}
