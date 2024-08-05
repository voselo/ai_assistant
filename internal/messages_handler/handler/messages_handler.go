package handler

import (
	"bytes"
	"encoding/json"
	"io"

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

	// Buffer the request body for multiple reads
	var buf bytes.Buffer
	tee := io.TeeReader(request.Body, &buf)

	// Attempt to decode as test request
	var testReq map[string]interface{}
	if err := json.NewDecoder(tee).Decode(&testReq); err == nil {
		if test, ok := testReq["test"].(bool); ok && test {
			responseWriter.WriteHeader(http.StatusOK)
			response := map[string]string{"status": "webhook connected"}
			json.NewEncoder(responseWriter).Encode(response)
			return
		}
	}

	request.Body = io.NopCloser(&buf)

	var messageRequest entity.MessageRequest

	err := json.NewDecoder(request.Body).Decode(&messageRequest)
	if err != nil {
		logger.Error("Failed to decode JSON:", err)
		http.Error(responseWriter, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Processing message
	if len(messageRequest.Messages) > 0 {
		message := messageRequest.Messages[0]
		handler.MessageRepository.HandleMessage(message)

		response := map[string]string{"status": "message handled successfully"}
		if err := json.NewEncoder(responseWriter).Encode(response); err != nil {
			logger.Error("Failed to send response:", err)
			http.Error(responseWriter, "Failed to send response", http.StatusInternalServerError)

		}
	} else {
		response := map[string]string{"status": "unknown error, message array was empty "}
		if err := json.NewEncoder(responseWriter).Encode(response); err != nil {
			logger.Error("Failed to send response:", err)
			http.Error(responseWriter, "Failed to send response", http.StatusInternalServerError)
		}
	}

}
