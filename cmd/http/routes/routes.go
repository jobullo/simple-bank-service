package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jobullo/go-api-example/config"
	"github.com/jobullo/go-api-example/database"
	_ "github.com/jobullo/go-api-example/docs"
	"github.com/jobullo/go-api-example/service"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// SetupRouter creates a router using middleware and controllers
func SetupRouter(cfg config.Configuration) *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.SetTrustedProxies([]string{"127.0.0.1"}) //only trust local proxy

	//set up swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Redirect the root to swagger
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
		c.Abort()
	})

	// Health endpoint
	health := new(HealthController)
	healthRoutes := router.Group("/health")
	{
		healthRoutes.GET("/", health.Status)
	}

	// Authentication endpoint
	auth := new(AuthController)
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", auth.Login)
	}

	//get database instance pointer
	db := database.GetDatabase()

	//initialize account service and controller
	accountService := service.NewAccountService(db)
	accountController := NewAccountController(accountService)

	// Account endpoints
	accountRoutes := router.Group("/accounts")
	{
		accountRoutes.GET("/", accountController.List)
		accountRoutes.GET("/:id", accountController.FetchById)
		accountRoutes.POST("/", accountController.Create)
		accountRoutes.PUT("/:id", accountController.Update)
		accountRoutes.DELETE("/:id", accountController.Delete)
	}

	//initialize transaction service and controller
	transactionService := service.NewTransactionService(db, *accountService)
	transactionController := NewTransactionController(transactionService)

	// Transaction endpoints
	transactionRoutes := router.Group("/transactions")
	{
		transactionRoutes.GET("/", transactionController.List)
		transactionRoutes.GET("/:id", transactionController.FetchById)
		transactionRoutes.POST("/", transactionController.Create)
		transactionRoutes.PUT("/:id", transactionController.Update)
		transactionRoutes.DELETE("/:id", transactionController.Delete)
	}

	return router
}
