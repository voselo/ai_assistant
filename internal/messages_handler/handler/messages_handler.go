package handler

import (
	"encoding/json"

	// "messages_handler/internal/domain/entity"
	"messages_handler/internal/messages_handler/domain/entity"
	"messages_handler/internal/messages_handler/domain/repository"
	"messages_handler/pkg/logging"
	"net/http"
	// "github.com/gorilla/mux"
)

type MessageHandler struct {
	MessageRepository repository.IMessageRepository
}

func NewMessageHandler(repository repository.IMessageRepository) *MessageHandler {
	return &MessageHandler{
		MessageRepository: repository,
	}
}

type WazzupTestRequest struct {
	Test bool `json:"test"`
}

func (handler *MessageHandler) HandleMessage(responseWriter http.ResponseWriter, request *http.Request) {
	logger := logging.GetLogger("trace")

	defer request.Body.Close()
	responseWriter.Header().Set("Content-Type", "application/json")

	var req entity.MessageRequest

	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		logger.Error("Failed to decode JSON:", err)
		http.Error(responseWriter, "Invalid JSON format", http.StatusBadRequest)
	}

	// Processing message
	message := req.Messages[0]
	handler.MessageRepository.HandleMessage(message)

	response := map[string]string{"status": "success"}
	if err := json.NewEncoder(responseWriter).Encode(response); err != nil {
		logger.Error("Failed to send response:", err)
		http.Error(responseWriter, "Failed to send response", http.StatusInternalServerError)
	}

}
