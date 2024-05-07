package routes

import (
	"zeroslope/config"

	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter creates a router using middleware and controllers
func SetupRouter(cfg config.Configuration) *gin.Engine {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

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

	return router
}
