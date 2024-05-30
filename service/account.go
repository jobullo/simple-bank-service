package service

import (
	"errors"

	"github.com/jobullo/go-api-example/database"

	gorm "github.com/jinzhu/gorm"
)

type AccountService struct {
	db *gorm.DB
}

// create a new account service
func NewAccountService(db *gorm.DB) *AccountService {
	return &AccountService{db: db}
}

// implement the create a new account method of the account service interface
func (ts *AccountService) Create(account *database.Account) error {

	//TODO - add validation
	// - balance is not negative
	// - account name is not empty
	// - account type is not empty and is a valid type

	//create from provided account struct object
	if result := ts.db.Create(account); result.Error != nil {
		return result.Error
	}
	return nil
}

// implement the FetchByID method of the account service interface
func (ts *AccountService) FetchById(id uint) (*database.Account, error) {
	var account database.Account
	//account ids are unique, so we can use First
	if result := ts.db.First(&account, id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, database.ErrNotFound
		}

		return nil, result.Error
	}

	return &account, nil
}

// implement the List method of the account service interface
func (ts *AccountService) List() (*[]database.Account, error) {
	var accounts []database.Account
	result := ts.db.Find(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}
	return &accounts, nil
}

// implement the Update method of the account service interface
func (ts *AccountService) Update(account *database.Account) error {
	var t database.Account
	// fetch the account by id
	if resp := ts.db.First(&t, account.Model.ID); resp.Error != nil {
		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			return database.ErrNotFound
		}
		return resp.Error
	}

	//TODO: Need to check that the account holder is not empty or is not the same as the current account holder

	// update the account holder
	t.AccountHolder = account.AccountHolder
	t.Balance = account.Balance //-- should updates to balance only be done through transactions?

	if resp := ts.db.Save(&t); resp.Error != nil {
		return resp.Error
	}

	//update the account timestamp
	account.Model.CreatedAt = t.Model.CreatedAt
	account.Model.UpdatedAt = t.Model.UpdatedAt

	return nil
}

// implement the Delete method of the account service interface
func (s *AccountService) Delete(id uint) error {
	var account database.Account

	//first check that the account exists
	if resp := s.db.First(&account, id); errors.Is(resp.Error, gorm.ErrRecordNotFound) {
		return database.ErrNotFound
	}

	//fetch associated transactions
	var transactions []database.Transaction
	if err := s.db.Model(&account).Association("Transactions").Find(&transactions).Error; err != nil {
		return err
	}

	//inline function to handle the deletion of the account and its associated transactions
	deleteAllAccountData := func(db *gorm.DB) error {

		// Delete each transaction
		for _, transaction := range transactions {
			if err := s.db.Delete(&transaction).Error; err != nil {
				return err
			}
		}

		// Delete the account
		if err := s.db.Delete(&account).Error; err != nil {
			return err
		}

		return nil
	}

	//will roll back the transaction if an error is returned by deleteAllAccountData
	if err := s.db.Transaction(deleteAllAccountData); err != nil {
		return err
	}

	return nil
}
