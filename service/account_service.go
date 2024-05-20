package service

import (
	"github.com/jobullo/go-api-example/database"

	"gorm.io/gorm"
)

type AccountService struct {
	db *gorm.DB
}

func NewAccountService(db *gorm.DB) *AcccountService {
	return &AccountService{db: db}
}

func (ts *AccountService) Create(account *database.Account) (*database.Account, error) {

	if result := ts.db.Create(account); result.Error != nil {
		return nil, result.Error
	}
	return account, nil
}

func (ts *AccountService) FetchById(id int) (*database.Account, error) {
	var account database.Account
	if result := ts.db.First(&account, id); result.Error != nil {
		return nil, result.Error
	}
	return &account, nil
}

func (ts *AccountService) List() ([]*database.Account, error) {
	var accounts []database.Account
	result := ts.db.Find(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}
	return &accounts, nil
}

func (ts *AccountService) Update(account *database.Account) error {
	var t database.Account
	if resp := ts.db.First(&t, account.ID); resp.Error != nil {
		return resp.Error
	}

	t.Amount = account.Amount
	t.Description = account.Description

	if resp := ts.db.Save(&t); resp.Error != nil {
		return resp.Error
	}

	account.CreatedAt = t.CreatedAt
	account.UpdatedAt = t.UpdatedAt

	return nil
}
