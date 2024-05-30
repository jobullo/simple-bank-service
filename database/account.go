package database

import "github.com/jinzhu/gorm"

type Account struct {
	gorm.Model                  //leaving this ananymous field here so gorm:embedded tag isn't necessary
	AccountHolder string        `json:"accountHolder" binding:"required"`
	AccountType   string        `json:"accountType" binding:"required"`
	Balance       float64       `json:"balance" binding:"required"`
	Transactions  []Transaction `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE;"`
}

type AccountService interface {
	Service[Account]
}
