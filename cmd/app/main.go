package main

import (
	"ai_assistant/config"
	"ai_assistant/internal/bootstrap"
	"ai_assistant/pkg/logging"

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

	factory := bootstrap.NewFactory(database)

	bootstrap.InitRouter(config, logger, factory)
}
