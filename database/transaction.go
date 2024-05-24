package database

import (
	"github.com/jinzhu/gorm"
)

type Transaction struct {
	gorm.Model          //leaving this ananymous field here so gorm:embedded tag isn't necessary
	AccountID  uint32   `json:"accountID" binding:"required"`
	Account    *Account `json:"account"`
	Type       string   `json:"transactionType" binding:"required"`
	Amount     float64  `json:"transactionAmount" binding:"required"`
}

type TransactionService interface {
	Service[Transaction]
	ListByAccount(accountID uint) (*[]Transaction, error)
}
