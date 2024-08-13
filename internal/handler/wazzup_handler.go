package handler

import (
	model "ai_assistant/internal/model/messages"
	"ai_assistant/internal/repository"
	"ai_assistant/pkg/logging"
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WazzupHandler struct {
	factory *repository.Factory
}

func NewWazzupHandler(factory *repository.Factory) *WazzupHandler {
	return &WazzupHandler{
		factory: factory,
	}
}

// Wazzup handle message
//
// @Tags         wazzup
// @Description  process message
// @Param		 hash  path string true "License hash"
// @Param        message  body      model.MessageRequest  true  "message Request"
// @Accept       json
// @Produce      json
// @Router       /ai/api/v1/wazzup/handle/{hash} [post]
func (handler *WazzupHandler) HandleMessage(ctx *gin.Context) {
	logger := logging.GetLogger("Info")

	data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.Error("Failed to read request body:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err 1"})
		return
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(data))

	var testRequest map[string]interface{}
	if err := ctx.ShouldBindJSON(&testRequest); err == nil && testRequest["test"] == true {
		logger.Info("Test webhook received and acknowledged")
		ctx.JSON(http.StatusOK, gin.H{"status": "webhook connected"})
		return
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(data))

	// License validation
	license := ctx.Param("hash")

	uid, err := handler.factory.CustomersRepository.ValidateLicense(license)

	if err != nil {
		logger.Error("Invalid license was provided: ", license)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err Auth 1"})
		return
	}

	var messageRequest model.MessageRequest
	if err := ctx.ShouldBindJSON(&messageRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error messagee format"})
	}

	// Message processing
	if len(messageRequest.Messages) > 0 {
		handler.factory.WazzupRepository.ProcessMessage(uid, messageRequest.Messages[0], handler.factory.CustomersRepository)
		ctx.JSON(http.StatusOK, gin.H{"response": "message was processed"})

	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "empty messages list"})
		return
	}

}
