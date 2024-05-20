package lambda

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	lambda "github.com/aws/aws-lambda-go/lambda"
	gorm "github.com/jinzhu/gorm"
	database "github.com/jobullo/go-api-example/database"
)

func connectToDatabase() *gorm.DB {
	fmt.Println(">> Connecting to database ...")
	portString := os.Getenv("DB_PORT")
	port, err := strconv.Atoi(portString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		port = 5432
	}
	database.Init(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		port,
		os.Getenv("DB_NAME"),
	)
	return database.GetDatabase()
}

func handler(ctx context.Context) {
	db := connectToDatabase()
	readEntity, _ := Read(ctx, db, 1) // Replace with the operation you want to execute
	fmt.Println(readEntity)
}

func main() {
	lambda.Start(handler)
}
