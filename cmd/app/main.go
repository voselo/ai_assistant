package main

import (
	"messages_handler/config"
	"messages_handler/internal/bootstrap"
	"messages_handler/pkg/logging"

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
