package console

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jobullo/go-api-example/database"
)

func main() {
	Execute()
}

func Execute() {
	fmt.Println("Thank you for using the bank example console application!")
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
	printBlue("$ these commands work for the following entities: Account and Transaction")
	printBlue("$ list -entity Transaction -accountID 1")
	printGray("     Will list all records from the transactions table associated with account 1.")
	printBlue("$ list -entity Account")
	printGray("     Will list all records from the Account table.")
	printBlue("$ read -entity <entity name> —id 1")
	printGray("     Will read the record with the id of 1 in the <entity name> table")
	printBlue("$ delete -entity <entity name> —id 1 ")
	printGray("     Will delete the record with the id of 1 in the <entity name> table.")
	printBlue("$ insert -entity Account -owner \"John Doe\" -type savings -balance 1000")
	printGray("     Will create a record in Account table with owner John Does with a $1000 balance in a savings account.")
	printBlue("$ insert -entity Transaction -account 1 -amount 1000 -type <deposit/withdrawal>")
	printGray("     Will create a transaction record that corresponds to account 1 for a deposit or withdrawal in the amount of $1000.")
	printBlue("$ update -entity Account -id 1 -owner \"John Doe\"")
	printGray("     Will update the owner of the account with id 1 to John Doe.")
	printBlue("$ update -entity Transaction -id 1 -account 1 -amount 1000")
	printGray("     Will update the transaction with id 1 in account 1 to have an amount of $1000.")
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
