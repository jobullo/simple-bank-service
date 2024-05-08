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

	printYellow(">> Loading console settings ...")
	LoadDotEnv()

	printYellow(">> Connecting to database ...")
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

	printYellow(">> Setting up database ...")
	database.BuildDatabase()

	fmt.Println("-----------------------------")
	fmt.Println("Example Commands for usage...")
	fmt.Println("-----------------------------")
	printBlue("$ list")
	printGray("     Will list all records in the SampleEntity table ")
	printBlue("$ read —id 7")
	printGray("     Will read the record with the id of 7 in the SampleEntity table")
	printBlue("$ delete —id 8 ")
	printGray("     Will delete the record with the id of 8 in the SampleEntity table.")
	printBlue("$ insert —name Chris —description \"Check this out\"")
	printGray("     Will insert a record with those details.")
	printBlue("$ update —id 3 —name Chris —description  \"Check this out again\"")
	printGray("     Will update a updated record with that id.")
	printBlue("-----------------------------")

	reader := bufio.NewReader(os.Stdin)
	for {
		printGreen("Enter your command: (type 'quit' to quit): ")
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		if cmd == "quit" {
			printRed("Exiting the application...")
			break
		}
		HandleCommands(cmd, db)
	}
}
