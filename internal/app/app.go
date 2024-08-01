package app

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"messages_handler/internal/config"
	"messages_handler/internal/handler"
	"messages_handler/pkg/logging"
)

func StartApp(
	config config.Config,
	logger logging.Logger,
	messageHandler handler.MessageHandler,
) {
	router := httprouter.New()
	router.HandlerFunc("POST", "/handle_message", messageHandler.HandleMessage)
	

	address := config.Server.Host + ":" + config.Server.Port
	logger.Infof("Server is working on %s", address)
	err := http.ListenAndServe(address, router)
	if err != nil {
		logger.Fatal(err)
	}
}
