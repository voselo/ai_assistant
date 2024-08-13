package bootstrap

import (
	"ai_assistant/config"
	// "ai_assistant/docs"
	"ai_assistant/pkg/logging"

	"ai_assistant/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	customersHandler "ai_assistant/internal/customers/handler"
	wazzupHandler "ai_assistant/internal/wazzup/handler"
)

func InitRouter(
	r *gin.Engine,
	config *config.Config,
	factory *Factory,
) {

	logger := logging.GetLogger("Info")

	// docs.SwaggerInfo.Host = config.BaseUrl
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	defaultConfig := cors.DefaultConfig()
	defaultConfig.AllowOrigins = []string{"*"}
	defaultConfig.AllowHeaders = []string{"*"}
	r.Use(cors.New(defaultConfig))

	// Customers
	customersHandler := customersHandler.New(config, *factory.CustomersRepository)
	customerRoutes := r.Group("/ai/api/v1/customers")
	customerRoutes.Use(middleware.AdminAuthMiddleware(config))

	{

		customerRoutes.POST("/create", customersHandler.Create)
		customerRoutes.PUT("/update/:id", customersHandler.Update)
		customerRoutes.GET("/get", customersHandler.GetAll)
		customerRoutes.GET("/get/:id", customersHandler.GetById)
		customerRoutes.DELETE("/delete/:id", customersHandler.Delete)
	}

	// Wazzup messages processing
	messagesHandler := wazzupHandler.New(factory.WazzupRepository, factory.CustomersRepository)
	wazzupRoutes := r.Group("/ai/api/v1/wazzup")
	{
		wazzupRoutes.POST("/handle/:hash", messagesHandler.HandleMessage)
	}

	host := config.Server.Host + ":" + config.Server.Port
	logger.Infof("Server is working on %s", host)

	if err := r.Run(host); err != nil {
		logger.Fatal(err)
	}
}
