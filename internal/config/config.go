package config

import (
	"os"
	"sync"

	"messages_handler/pkg/logging"

	"github.com/joho/godotenv"
)

type Config struct {
	ENVIRONMENT string

	Server struct {
		Host string `yaml:"host" env-default:"localhost"`
		Port string `yaml:"port" env-default:"8080"`
	}

	Webhook string
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

		instance.ENVIRONMENT = os.Getenv("APP_ENV")
		instance.Server.Host = os.Getenv("SERVER_HOST")
		instance.Server.Port = os.Getenv("SERVER_PORT")
		instance.Webhook = os.Getenv("WEBHOOK_URI")

		logger.Info("aplication is configured ", instance)
	})

	return instance

}
