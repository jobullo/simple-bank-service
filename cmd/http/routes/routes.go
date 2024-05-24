package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jobullo/go-api-example/config"
)

// SetupRouter creates a router using middleware and controllers
func SetupRouter(cfg config.Configuration) *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.SetTrustedProxies([]string{"127.0.0.1"}) //only trust local proxy

	router.Static("/swagger", "./path-to-swagger-ui")
	router.StaticFile("/swagger.yaml", "./swagger.yaml")

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

	// Sample endpoints
	sample := new(SampleController)
	sampleRoutes := router.Group("/samples")
	{
		sampleRoutes.GET("/", sample.List)
		sampleRoutes.GET("/:id", sample.Read)
		sampleRoutes.POST("/", sample.Create)
		sampleRoutes.PUT("/", sample.Update)
		sampleRoutes.DELETE("/:id", sample.List)
	}

	// Account endpoints
	account := new(AccountController)
	accountRoutes := router.Group("/accounts")
	{
		accountRoutes.GET("/", account.List)
		accountRoutes.GET("/:id", account.FetchById)
		accountRoutes.POST("/", account.Create)
		accountRoutes.PUT("/:id", account.Update)
		accountRoutes.DELETE("/:id", account.Delete)
	}

	transaction := new(TransactionController)
	transactionRoutes := router.Group("/transactions")
	{
		transactionRoutes.GET("/", transaction.List)
		transactionRoutes.GET("/:id", transaction.FetchById)
		transactionRoutes.POST("/", transaction.Create)
		transactionRoutes.PUT("/:id", transaction.Update)
		transactionRoutes.DELETE("/:id", transaction.Delete)
	}

	return router
}
