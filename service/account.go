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
func (as *AccountService) Create(account *database.Account) error {

	//TODO - add validation
	// - balance is not negative
	// - account name is not empty
	// - account type is not empty and is a valid type

	//create from provided account struct object
	if result := as.db.Create(account); result.Error != nil {
		return result.Error
	}
	return nil
}

// implement the FetchByID method of the account service interface
func (as *AccountService) FetchById(id uint) (*database.Account, error) {
	var account database.Account
	//account ids are unique, so we can use First
	if result := as.db.First(&account, id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, database.ErrNotFound
		}

		return nil, result.Error
	}

	return &account, nil
}

// implement the List method of the account service interface
func (as *AccountService) List() (*[]database.Account, error) {
	var account []database.Account
	result := as.db.Find(&account)
	if result.Error != nil {
		return nil, result.Error
	}
	return &account, nil
}

// implement the Update method of the account service interface
func (as *AccountService) Update(account *database.Account) error {
	var acc database.Account
	// fetch the account by id
	if resp := as.db.First(&acc, account.Model.ID); resp.Error != nil {
		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			return database.ErrNotFound
		}
		return resp.Error
	}

	//TODO: Need to check that the account holder is not empty or if it is different from the current account holder

	// update the account holder
	acc.AccountHolder = account.AccountHolder
	acc.Balance = account.Balance //-- should updates to balance only be done through transactions?

	if resp := as.db.Save(&acc); resp.Error != nil {
		return resp.Error
	}

	//update the account timestamp
	account.Model.CreatedAt = acc.Model.CreatedAt
	account.Model.UpdatedAt = acc.Model.UpdatedAt

	return nil
}

// implement the Delete method of the account service interface
func (as *AccountService) Delete(id uint) error {
	var account database.Account

	//first check that the account exisas
	if resp := as.db.First(&account, id); errors.Is(resp.Error, gorm.ErrRecordNotFound) {
		return database.ErrNotFound
	}

	//fetch associated transactions
	var transactions []database.Transaction
	if err := as.db.Model(&account).Association("Transactions").Find(&transactions).Error; err != nil {
		return err
	}

	//inline function to handle the deletion of the account and its associated transactions
	deleteAllAccountData := func(db *gorm.DB) error {

		// Delete each transaction
		for _, transaction := range transactions {
			if err := as.db.Delete(&transaction).Error; err != nil {
				return err
			}
		}

		// Delete the account
		if err := as.db.Delete(&account).Error; err != nil {
			return err
		}

		return nil
	}

	//will roll back the transaction if an error is returned by deleteAllAccountData
	if err := as.db.Transaction(deleteAllAccountData); err != nil {
		return err
	}

	return nil
}
