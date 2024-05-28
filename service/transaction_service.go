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

// implements the create a new transaction method of the transaction service interface
func (ts *TransactionService) Create(transaction *database.Transaction) error {
	//check that the account exists
	account, err := ts.accountService.FetchById(uint(transaction.AccountID))
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return database.ErrParentNotFound
		}
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
			return database.ErrInvalidType
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

// implements the FetchByID method of the transaction service interface
func (ts *TransactionService) FetchById(id int) (*database.Transaction, error) {
	var transaction database.Transaction
	if result := ts.db.First(&transaction, id); result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, database.ErrNotFound
		}

		return nil, result.Error
	}

	return &transaction, nil
}

// implements the List method of the transaction service interface
func (ts *TransactionService) List() (*[]database.Transaction, error) {
	return ts.list(nil)
}

// private of the transaction service
func (ts *TransactionService) list(criteria *database.Transaction) (*[]database.Transaction, error) {

	db := ts.db.Preload("Account") //preloads the account object

	if criteria != nil {
		db = db.Where(criteria)
	}

	var result []database.Transaction

	if resp := db.Find(&result); resp.Error != nil {
		return nil, resp.Error
	}

	return &result, nil
}

func (ts *TransactionService) ListByAccount(accountID uint) (*[]database.Transaction, error) {
	criteria := database.Transaction{AccountID: uint32(accountID)}
	return ts.list(&criteria)
}

//func (ts *TransactionService) List(accountID uint) (*[]database.Transaction, error) {
//	var transactions []database.Transaction
//	//SELECT * FROM transactions WHERE account = accountID
//	result := ts.db.Where("account = ?", accountID).Find(&transactions)
//	if result.Error != nil {
//		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
//			return nil, database.ErrParentNotFound
//		}
//		return nil, result.Error
//	}
//	return &transactions, nil
//}

// list all transactions that belong to a specific account
//func (ts *TransactionService) ListByAccount(accountID uint) (*[]database.Transaction, error) {
//	var transactions []database.Transaction
//
//	//SELECT * FROM transactions WHERE account = accountID
//	result := ts.db.Where("account = ?", accountID).Find(&transactions)
//
//	if result.Error != nil {
//		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
//			return nil, database.ErrParentNotFound
//		}
//		return nil, result.Error
//	}
//
//	return &transactions, nil
//}

// implements the Update method of the transaction service interface
func (ts *TransactionService) Update(transaction *database.Transaction) error {
	var t database.Transaction

	//TODO: validate transaction amount
	// - can't be <= 0 (transaction type determines if it subtracts or adds to the account balance)
	// - for a withdrawal, the amount can't be greater than the account balance

	if resp := ts.db.First(&t, transaction.Model.ID); resp.Error != nil {

		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			return database.ErrNotFound
		}

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

// implements the delete method of the transaction service interface
func (ts *TransactionService) Delete(id uint) error {
	var transaction database.Transaction

	if result := ts.db.First(&transaction, id); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return database.ErrNotFound
	}

	return ts.db.Delete(&transaction, id).Error
}
