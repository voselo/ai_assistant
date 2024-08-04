package bootstrap

import (
	"net/http"

	"github.com/gorilla/mux"

	"messages_handler/internal/config"
	"messages_handler/internal/messages_handler/handler"
	"messages_handler/pkg/logging"
)

func InitRouter(
	config config.Config,
	logger logging.Logger,
	messageHandler handler.MessageHandler,
) {
	router := mux.NewRouter()

	// Wazzup messages handling
	router.HandleFunc("/ai/api/v1/messages_handler/handle_message", messageHandler.HandleMessage).Methods("POST")

	address := config.Server.Host + ":" + config.Server.Port
	logger.Infof("Server is working on %s", address)

	err := http.ListenAndServe(address, router)
	if err != nil {
		logger.Fatal(err)
	}
}
