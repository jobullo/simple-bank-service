package console

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	gorm "github.com/jinzhu/gorm"

	database "github.com/jobullo/go-api-example/database"
	service "github.com/jobullo/go-api-example/service"
)

func handleAccountOperations(command string, argMap map[string]string, db *gorm.DB) {

	newAccountService := service.NewAccountService(db)

	switch command {
	case "list":
		accounts, err := newAccountService.List()
		if err != nil {
			fmt.Println("  Error fetching accounts:", err)
			return
		}
		for _, account := range *accounts {
			fmt.Printf("Account #: %d, Owner: %v, Balance: %v\n", account.ID, account.AccountHolder, account.Balance)
		}
	case "read":
		idString := argMap["id"]
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			fmt.Println("  Invalid ID:", idString)
			return
		}
		//var account database.Account
		account, err := newAccountService.FetchById(uint(id))
		if err != nil {
			fmt.Println("  Error fetching account:", err)
			return
		}
		fmt.Printf("Account #: %d, Owner: %v, Balance: %v\n", account.Model.ID, account.AccountHolder, account.Balance)
	case "delete":
		idString := argMap["id"]
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			fmt.Println("  Invalid ID:", idString)
			return
		}
		newAccountService.Delete(uint(id))
		fmt.Println("  Deleted account with ID:", id)
	case "update":
		idString := argMap["id"]
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			fmt.Println("  Invalid ID:", idString)
			return
		}
		ownerString := argMap["owner"]
		account, err := newAccountService.FetchById(uint(id))
		if err != nil {
			fmt.Println("  Error fetching account:", err)
			return
		}
		account.AccountHolder = ownerString
		if err := newAccountService.Update(account); err != nil {
			fmt.Println("  Error updating account:", err)
			return
		}
		fmt.Println("  Updated owner of account with ID:", id)
	case "insert":
		ownerString := argMap["owner"]
		typeString := argMap["type"]
		balanceString := argMap["balance"]
		balance, err := strconv.ParseFloat(balanceString, 64)
		if err != nil {
			fmt.Println("  Invalid balance:", balanceString)
			return
		}
		account := database.Account{AccountHolder: ownerString, AccountType: typeString, Balance: balance}
		if err := newAccountService.Create(&account); err != nil {
			fmt.Println("  Error creating account:", err)
			return
		}
		fmt.Println("  Inserted new account with ID:", account.Model.ID)
	default:
		fmt.Println("  Unknown command.")
	}
}

func handleTransactionOperations(command string, argMap map[string]string, db *gorm.DB) {
	newAccountService := service.NewAccountService(db)
	newTransactionService := service.NewTransactionService(db, *newAccountService)
	switch command {
	case "list":
		var transactions *[]database.Transaction
		idString := argMap["accountID"]
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			fmt.Println("  Invalid Account ID:", idString)
			return
		}
		transactions, err = newTransactionService.ListByAccount(uint(id))
		if err != nil {
			fmt.Println("  Error fetching transactions:", err)
		}

		for _, transaction := range *transactions {
			fmt.Printf("  ID: %d, Account: %d, Amount: %v, Type: %s\n", transaction.ID, transaction.AccountID, transaction.Amount, transaction.Type)
		}
	case "read":
		idString := argMap["id"]
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			fmt.Println("  Invalid ID:", idString)
			return
		}

		var transaction database.Transaction
		db.First(&transaction, id)
		fmt.Printf("  ID: %d, Account: %d, Amount: %v, Type: %s\n", transaction.ID, transaction.AccountID, transaction.Amount, transaction.Type)
	case "delete":
		idString := argMap["id"]
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			fmt.Println("  Invalid ID:", idString)
			return
		}

		newTransactionService.Delete(uint(id))
		db.Delete(&database.Transaction{}, id)
		fmt.Println("  Deleted transaction with ID:", id)
	case "update":
		idString := argMap["id"]
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			fmt.Println("  Invalid ID:", idString)
			return
		}

		amountString := argMap["amount"]
		amount, err := strconv.ParseFloat(amountString, 64)
		if err != nil {
			fmt.Println("  Invalid amount:", amountString)
			return
		}

		model := gorm.Model{ID: uint(id)}

		newTransactionService.Update(&database.Transaction{Model: model, Amount: amount})
		fmt.Println("  Updated transaction with ID:", id)
	case "insert":
		accountString := argMap["account"]
		account, err := strconv.ParseUint(accountString, 10, 32)
		if err != nil {
			fmt.Println("  Invalid Account ID:", accountString)
			return
		}

		amountString := argMap["amount"]
		amount, err := strconv.ParseFloat(amountString, 64)
		if err != nil {
			fmt.Println("  Invalid amount:", amountString)
			return
		}

		transactionType := argMap["type"]
		transaction := database.Transaction{AccountID: uint32(account), Amount: amount, Type: transactionType}

		if error := newTransactionService.Create(&transaction); error != nil {
			fmt.Println("  Error creating transaction:", error)
			return
		}
		fmt.Println("  Inserted new transaction with ID:", transaction.Model.ID)

	default:
		fmt.Println("  Unknown command.")
	}
}

func HandleCommands(cmd string, db *gorm.DB) {

	// Regex pattern to capture key-value pairs
	pattern := `-(\w+)\s+([^-\s]+(?:\s+[^-\s]+)*)?`

	re := regexp.MustCompile(pattern)

	matches := re.FindAllStringSubmatch(cmd, -1)
	argMap := make(map[string]string)
	for _, match := range matches {
		key := match[1]
		value := strings.Trim(match[2], `"`) // Remove surrounding quotes if present
		argMap[key] = value
	}

	command := strings.Fields(cmd)[0]
	entity := argMap["entity"]

	if entity != "Account" && entity != "Transaction" {
		fmt.Println("Invalid entity specified.")
		return
	}

	if entity == "Account" {
		handleAccountOperations(command, argMap, db)
	} else if entity == "Transaction" {
		handleTransactionOperations(command, argMap, db)
	}
}
