package config

import (
	"os"
	"sync"

	"messages_handler/pkg/logging"

	"github.com/joho/godotenv"
)

type Config struct {
	ApiKey string

	Mode string

	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}

	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		SSLMode  string `yaml:"sslmode"`
	}
}

var (
	instance *Config
	once     sync.Once
)

func Init() *Config {
	once.Do(func() {
		logger := logging.GetLogger("warn")
		logger.Info("Read aplication config")

		instance = &Config{}

		// Environment setup
		if _, err := os.Stat(".env.dev"); err == nil {
			if err := godotenv.Load(".env.dev"); err != nil {
				logger.Error("Error loading .env file")
			}
		}

		instance.Mode = os.Getenv("MODE")

		instance.ApiKey = os.Getenv("API_KEY")

		// Server setup
		instance.Server.Host = os.Getenv("SERVER_HOST")
		instance.Server.Port = os.Getenv("SERVER_PORT")

		// Database setup
		instance.Database.Host = os.Getenv("DB_HOST")
		instance.Database.Port = os.Getenv("DB_PORT")
		instance.Database.User = os.Getenv("DB_USER")
		instance.Database.Password = os.Getenv("DB_PASSWORD")
		instance.Database.Name = os.Getenv("DB_NAME")
		instance.Database.SSLMode = os.Getenv("DB_SSL_MODE")

		logger.Info("aplication is configured ", instance)

	})

	return instance

}
