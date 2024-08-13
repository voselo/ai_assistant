package middleware

import (
	"ai_assistant/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	tokenHeaderName = "x-token"
)

func AdminAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.GetHeader(tokenHeaderName)

		if key != cfg.ApiKey {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Err auth 1"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
