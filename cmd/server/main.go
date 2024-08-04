package main

import (
	"messages_handler/internal/bootstrap"
	"messages_handler/internal/config"
	"messages_handler/internal/messages_handler/handler"
	"messages_handler/pkg/logging"
)

func main() {
	// Logging setup
	logger := logging.GetLogger("trace")
	logger.Info("Logger is working")

	// Read config
	config := config.LoadConfig()

	// bootstrap
	repositoryFactory := bootstrap.NewRepositoryFactory(config)
	messageHandler := handler.NewMessageHandler(repositoryFactory.MessageRepository)

	// Starting app
	bootstrap.InitRouter(*config, logger, *messageHandler)
}
