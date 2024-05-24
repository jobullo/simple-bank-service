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
		return resp.Error
	}

	// update the account holder
	t.AccountHolder = account.AccountHolder
	t.Balance = account.Balance

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

	if resp := s.db.First(&account, id); errors.Is(resp.Error, gorm.ErrRecordNotFound) {
		return gorm.ErrRecordNotFound
	}

	return s.db.Delete(&account).Error
}
