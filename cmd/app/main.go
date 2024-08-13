package main

import (
	"ai_assistant/config"
	"ai_assistant/internal/bootstrap"
	"ai_assistant/pkg/logging"

	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/voselo/ai_assistant/docs"
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
