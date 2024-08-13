package main

import (
	"ai_assistant/config"
	"ai_assistant/internal/bootstrap"
	"ai_assistant/internal/repository"
	"ai_assistant/pkg/logging"

	"github.com/gin-gonic/gin"

	_ "ai_assistant/docs"

	_ "github.com/jackc/pgx/stdlib"
)

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						x-token
func main() {
	logger := logging.GetLogger("Info")
	logger.Info("App is started")

	// bootstrap
	config := config.Init()
	database := bootstrap.InitDB(config)

	factory := repository.NewFactory(database)

	gin.SetMode(config.Mode)
	router := gin.Default()

	bootstrap.InitRouter(router, config, factory)


}
