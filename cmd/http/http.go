package http

import (
	"fmt"

	routes "zeroslope/cmd/http/routes"
	"zeroslope/config"
	"zeroslope/database"

	"github.com/gin-contrib/cors"
)

// @title ZeroSlope API
// @version 2.0
// @description Zeroslope microservice architecture written in go
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func Execute() {
	fmt.Println("Thank you for using the ZeroSlope microservice!")
	fmt.Println("-----------------------------")

	fmt.Println(">> Loading configuration ...")
	cfg := config.LoadConfigFromPath("config.yaml")

	fmt.Println(">> Loading routes ...")
	router := routes.SetupRouter(cfg)

	fmt.Println(">> Loading middleware ...")
	router.Use(AuthMiddleware(cfg))

	fmt.Println(">> Loading api settings ...")
	router.Use(cors.Default())

	fmt.Println(">> Connecting to database ...")
	database.Init(cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Database)

	fmt.Println(">> Setting up database ...")
	database.BuildDatabase()

	fmt.Println(">> Starting service ...")
	router.Run(cfg.Server.Port)

	fmt.Println("## Success! Service is available at http://localhost: %d", cfg.Server.Port)
}
