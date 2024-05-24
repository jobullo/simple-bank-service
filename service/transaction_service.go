package service

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/jobullo/go-api-example/database"
)

type TransactionService struct {
	db             *gorm.DB
	accountService AccountService //dependency injection - keep balance in sync
}

func NewTransactionService(db *gorm.DB, accountService AccountService) *TransactionService {
	return &TransactionService{
		db:             db,
		accountService: accountService,
	}
}

// implement the create a new transaction method of the transaction service interface
func (ts *TransactionService) Create(transaction *database.Transaction) error {
	//check that the account exists
	account, err := ts.accountService.FetchById(uint(transaction.AccountID))
	if err != nil {
		return err
	}

	//populate account before insertion into DB
	transaction.Account = account

	//inline function to handle to pass to db.Transaction
	performTransaction := func(db *gorm.DB) error {
		//create from provided transaction struct object
		if result := ts.db.Create(&transaction); result.Error != nil {
			return result.Error
		}

		//compute the new account balance
		if transaction.Type == "deposit" {
			account.Balance += transaction.Amount
		} else if transaction.Type == "withdrawal" {
			account.Balance -= transaction.Amount
		} else {
			return errors.New("invalid transaction type")
		}

		//now update the account
		if err := ts.accountService.Update(account); err != nil {
			return err
		}

		return nil
	}

	//will roll back the transaction if an error is returned by accountTransaction
	if err := ts.accountService.db.Transaction(performTransaction); err != nil {
		return err
	}

	return nil
}

// implement the FetchByID method of the transaction service interface
func (ts *TransactionService) FetchById(id int) (*database.Transaction, error) {
	var transaction database.Transaction
	result := ts.db.First(&transaction, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &transaction, nil
}

// implement the List method of the transaction service interface
func (ts *TransactionService) List() (*[]database.Transaction, error) {
	var transactions []database.Transaction
	//only return transactions associated with the account id
	result := ts.db.Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}
	return &transactions, nil
}

// list all transactions that belong to a specific account
func (ts *TransactionService) ListByAccount(accountID uint) (*[]database.Transaction, error) {
	var transactions []database.Transaction
	result := ts.db.Where("account_id = ?", accountID).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}
	return &transactions, nil
}

// implement the Update method of the transaction service interface
func (ts *TransactionService) Update(transaction *database.Transaction) error {
	var t database.Transaction

	if resp := ts.db.First(&t, transaction.Model.ID); resp.Error != nil {
		return resp.Error
	}

	t.Amount = transaction.Amount
	if resp := ts.db.Save(&t); resp.Error != nil {
		return resp.Error
	}

	transaction.Type = t.Type
	transaction.AccountID = t.AccountID
	transaction.Model.CreatedAt = t.Model.CreatedAt
	transaction.Model.UpdatedAt = t.Model.UpdatedAt

	return nil
}

// implement the delete method of the transaction service interface
func (ts *TransactionService) Delete(id uint) error {
	var transaction database.Transaction
	result := ts.db.Delete(&transaction, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
