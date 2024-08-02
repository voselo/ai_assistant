package main

import (
	"messages_handler/config"
	"messages_handler/internal/bootstrap"
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

	// DI
	messagesRepository := repository.NewInMemoryMessageRepository()
	messagesService := service.NewMessageService(messagesRepository)
	messageHandler := handler.NewMessageHandler(messagesService)

	// Starting app
	bootstrap.InitApp(*config, logger, *messageHandler)
}
