package console

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"zeroslope/database"
)

func main() {
	Execute()
}

func Execute() {
	fmt.Println("Thank you for using the ZeroSlope console application!")
	fmt.Println("You can use this console to interact with the database. ")
	fmt.Println("-----------------------------")

	fmt.Println(">> Loading console settings ...")
	LoadDotEnv()

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
	db := database.GetDatabase()

	fmt.Println(">> Setting up database ...")
	database.BuildDatabase()

	fmt.Println("-----------------------------")
	fmt.Println("Example Commands for usage...")
	fmt.Println("-----------------------------")
	fmt.Println("$ list")
	fmt.Println("     Will list all records in the SampleEntity table ")
	fmt.Println("$ read —id 7")
	fmt.Println("     Will read the record with the id of 7 in the SampleEntity table")
	fmt.Println("$ delete —id 8 ")
	fmt.Println("     Will delete the record with the id of 8 in the SampleEntity table.")
	fmt.Println("$ insert —name Chris —description \"Check this out\"")
	fmt.Println("     Will insert a record with those details.")
	fmt.Println("$ update —id 3 —name Chris —description  \"Check this out again\"")
	fmt.Println("     Will update a updated record with that id.")
	fmt.Println("-----------------------------")

	fmt.Println(">> Please enter your first command to begin")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter command (type 'quit' to quit): ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		if cmd == "quit" {
			fmt.Println("Exiting the application...")
			break
		}
		HandleCommands(cmd, db)
	}
}
