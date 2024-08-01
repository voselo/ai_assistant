package main

import (
	"messages_handler/internal/app"
	"messages_handler/internal/config"
	"messages_handler/internal/domain/repository"
	"messages_handler/internal/domain/service"
	"messages_handler/internal/handler"
	"messages_handler/pkg/logging"
)

func main() {
	// Logging setup
	logger := logging.GetLogger("trace")
	logger.Info("Logger is working")

	// Read config
	config := config.LoadConfig()

	logger.Info("Config in main ", config)

	messagesRepository := repository.NewInMemoryMessageRepository()
	messagesService := service.NewMessageService(messagesRepository)
	messageHandler := handler.NewMessageHandler(messagesService)

	// Start app
	app.StartApp(*config, logger, *messageHandler)
}
