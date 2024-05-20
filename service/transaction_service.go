package service

import ("gorm.io/gorm"
		"github.com/jobullo/go-api-example/database"
)

type TransactionService struct {
	db *db.gorm.DB
}

func NewTransactionService(db *db.gorm.DB) *TransactionService {
	return &TransactionService{db: db}
}

func (ts *TransactionService) Create(transaction *database.Transaction) (*database.Transaction, error) {
	result := ts.db.Create(&transaction)
	if result.Error != nil {
		return nil, result.Error
	}
	return transaction, nil
}

func (ts *TransactionService) FetchById(id int) (*database.Transaction, error) {	
	var transaction database.Transaction
	result := ts.db.First(&transaction, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &transaction, nil
}

func (ts *TransactionService) List() ([]*database.Transaction, error) {
	var transactions []database.Transaction
	result := ts.db.Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}
	return &transactions, nil
}

func (ts *TransactionService) Update(transaction *database.Transaction) (error) {
	var t database.Transaction
	if resp:= ts.db.First(&t, transaction.ID); resp.Error != nil {
		return nil, resp.Error
	}

	t.Amount = transaction.Amount
	t.Description = transaction.Description
	if resp := ts.db.Save(&t); resp.Error != nil {
		return nil, resp.Error
	}

	transaction.CreatedAt = t.CreatedAt
	transaction.UpdatedAt = t.UpdatedAt

	return nil
}