package http

import (
	"fmt"

	routes "github.com/jobullo/go-api-example/cmd/http/routes"
	config "github.com/jobullo/go-api-example/config"
	database "github.com/jobullo/go-api-example/database"

	"github.com/gin-contrib/cors"
)

// @title bank-example API
// @version 2.0
// @description bank-example microservice architecture written in go
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func Execute() {
	fmt.Println("Thank you for using the bank example microservice!")
	fmt.Println("-----------------------------")

	fmt.Println(">> Loading configuration ...")
	cfg := config.LoadConfigFromPath("config.yaml")

	fmt.Println(">> Connecting to database ...")
	database.Init(cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Database)

	fmt.Println(">> Setting up database ...")
	database.BuildDatabase()

	fmt.Println(">> Loading routes ...")
	router := routes.SetupRouter(cfg)

	fmt.Println(">> Loading middleware ...")
	router.Use(AuthMiddleware(cfg))

	fmt.Println(">> Loading api settings ...")
	router.Use(cors.Default())

	fmt.Println(">> Starting service ...")
	router.Run(cfg.Server.Port)

	fmt.Printf("## Success! Service is available at http://localhost: %v\n", cfg.Server.Port)
}
