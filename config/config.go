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
		Port string `yaml:"port" env-default:"8080"`
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

		env := os.Getenv("APP_ENV")
		if env == "" {
			env = "dev"
		}

		envFileName := ".env." + env
		if err := godotenv.Load(envFileName); err != nil {
			logger.Fatal("Error loading " + envFileName + " file")
		}

		instance.Server.Host = os.Getenv("SERVER_HOST")
		instance.Server.Port = os.Getenv("SERVER_PORT")

		logger.Info("aplication is configured ", instance)
	})

	return instance

}
