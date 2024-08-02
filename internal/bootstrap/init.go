package bootstrap

import (
	"net/http"

	"github.com/gorilla/mux"

	"messages_handler/config"
	"messages_handler/internal/handler"
	"messages_handler/pkg/logging"
)

func InitApp(
	config config.Config,
	logger logging.Logger,
	messageHandler handler.MessageHandler,
) {
	router := mux.NewRouter()
	router.HandleFunc("/ai/api/v1/handle_message", messageHandler.HandleMessage).Methods("POST")

	address := config.Server.Host + ":" + config.Server.Port
	logger.Infof("Server is working on %s", address)
	err := http.ListenAndServe(address, router)
	if err != nil {
		logger.Fatal(err)
	}
}
