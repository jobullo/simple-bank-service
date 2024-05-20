package database

type Transaction struct {
	Model             Model
	AccountNumber     uint    `json:"accountNumber"`
	TransactionType   string  `json:"transactionType"`
	TransactionAmount float64 `json:"transactionAmount"`
}

type TransactionService interface {
	Service[Transaction]
}
