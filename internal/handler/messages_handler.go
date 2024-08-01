package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"messages_handler/internal/domain/service"
	"messages_handler/pkg/logging"
	// "github.com/gorilla/mux"
)

type MessageHandler struct {
	service *service.MessageService
}

func NewMessageHandler(service *service.MessageService) *MessageHandler {
	return &MessageHandler{
		service: service,
	}
}

func (handler *MessageHandler) HandleMessage(responseWriter http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger("trace")

	logger.Info("Message handled")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("Failed to read request body:", err)
		http.Error(responseWriter, "Can't read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Вывод сырого тела запроса в лог (необязательно)
	logger.Info("Request body:\n", string(body))

	var request struct {
		ChannelID string `json:"channel_id"`
		Content   string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	if err := handler.service.AddMessage(request.ChannelID, request.Content); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	responseWriter.WriteHeader(http.StatusOK)
}

// func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	channelID := vars["channelID"]

// 	messages, err := h.service.GetMessages(channelID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(messages)
// }
