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

	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err 3"})

}
