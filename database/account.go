package database

type Account struct {
	Model          Model
	AccountHolder  string        `json:"accountHolder"`
	AccountNumber  uint          `json:"accountNumber"`
	AccountHistory []Transaction `json:"accountHistory"`
	Balance        float64       `json:"balance"`
}

type AccountService interface {
	Service[Account]
}
