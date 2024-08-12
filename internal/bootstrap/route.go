package bootstrap

import (
	"messages_handler/config"
	"messages_handler/pkg/logging"

	"github.com/gin-gonic/gin"
	// ginSwagger "github.com/swaggo/gin-swagger"

	customersHandler "messages_handler/internal/customers/handler"
	"messages_handler/internal/middleware"
	wazzupHandler "messages_handler/internal/wazzup/handler"
)

func InitRouter(
	config *config.Config,
	logger *logging.Logger,
	factory *Factory,
) {
	gin.SetMode(config.Mode)
	r := gin.Default()
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Customers
	customersHandler := customersHandler.New(config, *factory.CustomersRepository)
	customerRoutes := r.Group("/ai/api/v1/customers")
	customerRoutes.Use(middleware.AdminAuthMiddleware(config))

	{

		customerRoutes.GET("/get", customersHandler.Get)
		customerRoutes.POST("/create", customersHandler.Create)
		customerRoutes.PUT("/update", customersHandler.Update)
		customerRoutes.DELETE("/delete", customersHandler.Delete)
	}

	// Wazzup messages processing
	messagesHandler := wazzupHandler.New(factory.WazzupRepository, factory.CustomersRepository)
	wazzupRoutes := r.Group("/ai/api/v1/wazzup")
	{
		wazzupRoutes.POST("/handle", messagesHandler.HandleMessage)
	}

	host := config.Server.Host + ":" + config.Server.Port
	logger.Infof("Server is working on %s", host)

	if err := r.Run(host); err != nil {
		logger.Fatal(err)
	}
}
