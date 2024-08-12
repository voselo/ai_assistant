package handler

import (
	"bytes"
	"io"
	customersRepo "messages_handler/internal/customers/repository"
	"messages_handler/internal/wazzup/model"
	wazzupRepo "messages_handler/internal/wazzup/repository"
	"messages_handler/pkg/logging"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WazzupHandler struct {
	wazzupRepo    *wazzupRepo.WazzupRepository
	customersRepo *customersRepo.CustomersRepository
}

func New(wazzupRepo *wazzupRepo.WazzupRepository, customersRepo *customersRepo.CustomersRepository) *WazzupHandler {
	return &WazzupHandler{
		wazzupRepo:    wazzupRepo,
		customersRepo: customersRepo,
	}
}

func (handler *WazzupHandler) HandleMessage(ctx *gin.Context) {
	logger := logging.GetLogger("Info")

	// Сохраняем тело запроса
	data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.Error("Failed to read request body:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err 1"})
		return
	} // Восстанавливаем тело запроса для первого парсинга
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(data))

	// Попытка преобразовать в структуру для тестового запроса
	var testRequest map[string]interface{}
	if err := ctx.ShouldBindJSON(&testRequest); err == nil && testRequest["test"] == true {
		logger.Info("Test webhook received and acknowledged")
		ctx.JSON(http.StatusOK, gin.H{"status": "webhook connected"})
		return
	}

	// Восстанавливаем тело запроса для второго парсинга
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(data))

	// // Валидация лицензии
	license := ctx.GetHeader("x-license")

	isValid, _ := handler.customersRepo.ValidateLicense(license)

	if !isValid {
		logger.Error("Invalid license was provided: ", license)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err 2"})
		return
	}

	// Попытка преобразовать в структуру сообщения
	var messageRequest model.MessageRequest
	if err := ctx.ShouldBindJSON(&messageRequest); err == nil {
		// Обработка сообщения...
		ctx.JSON(http.StatusOK, gin.H{"status": "message processed"})
		return
	}

	// Если ни один из парсингов не удался
	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err 3"})

	// // Декодируем входящее тело запроса в JSON
	// var payload map[string]interface{}
	// if err := ctx.ShouldBindJSON(&payload); err != nil {
	// 	logger.Error("Failed to decode JSON:", err)
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err 1"})
	// 	return
	// }

	// // Обрабатываем тестовый запрос
	// if test, ok := payload["test"].(bool); ok && test {
	// 	logger.Info("Test webhook received and acknowledged")
	// 	ctx.JSON(http.StatusOK, gin.H{"status": "webhook connected"})
	// 	return
	// }

	// // Преобразуем JSON в структуру MessageRequest
	// var messageRequest model.MessageRequest
	// if err := ctx.BindJSON(&messageRequest); err != nil {
	// 	logger.Error("Failed to decode message request JSON:", err)
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message request"})
	// 	return
	// }

	// // Проверяем наличие сообщений и обрабатываем первое сообщение
	// if len(messageRequest.Messages) > 0 {
	// 	message := messageRequest.Messages[0]
	// 	handler.wazzupRepo.ProcessMessage(message)
	// 	// if err := handler.repo.ProcessMessage(message); err != nil {
	// 	// 	logger.Error("Failed to process message:", err)
	// 	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process message"})
	// 	// 	return
	// 	// }

	// 	logger.WithExtraFields(map[string]interface{}{
	// 		"channel_id": message.ChannelId,
	// 		"message_id": message.MessageId,
	// 	}).Info("Message handled successfully")

	// 	ctx.JSON(http.StatusOK, gin.H{"status": "message has been processed"})
	// } else {
	// 	logger.Warn("Message array was empty")
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"status": "message array was empty"})
	// }

	// // Buffer the request body for multiple reads
	// var buf bytes.Buffer
	// tee := io.TeeReader(ctx.Request.Body, &buf)

	// // Attempt to decode as test request
	// var testReq map[string]interface{}
	// if err := json.NewDecoder(tee).Decode(&testReq); err == nil {
	// 	if test, ok := testReq["test"].(bool); ok && test {
	// 		ctx.JSON(http.StatusOK, gin.H{"status": "webhook connected"})
	// 		return
	// 	}
	// }

	// // Перезаписываем тело запроса
	// ctx.Request.Body = io.NopCloser(&buf)

	// var messageRequest model.MessageRequest

	// if err := ctx.BindJSON(&messageRequest); err != nil {
	// 	logger.Error("Failed to decode JSON:", err)
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error 1"})
	// 	return
	// }

	// // Обработка сообщений
	// if len(messageRequest.Messages) > 0 {
	// 	message := messageRequest.Messages[0]
	// 	handler.MessageRepository.ProcessMessage(message)

	// 	ctx.JSON(http.StatusOK, gin.H{"status": "message handled successfully"})
	// } else {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"status": "unknown error, message array was empty"})
	// }

}
