package service

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jobullo/go-api-example/database"
	"github.com/jobullo/go-api-example/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// A new test suite is created by embedding
// the suite.Suite struct.
type TransactionServiceSuite struct {
	suite.Suite
	assert       *assert.Assertions
	sqlmock      sqlmock.Sqlmock
	transService *TransactionService
	acctService  *AccountService
	account      database.Account
}

// Invoke this function to run the test suite with "go test" at the CLI
func TestTransactionServiceSuite(t *testing.T) {
	suite.Run(t, new(TransactionServiceSuite))
}

func (ts *TransactionServiceSuite) SetupTest() {
	t := ts.T() //gets the current *testing.T context, used by testing framework to manage test state

	db, sql, err := mock.DB() //mock.DB() returns a gorm.DB and sqlmock.Sqlmock
	require.NoError(t, err)   //use require from testify to stop test if error is not nil

	ts.assert = assert.New(t) //create a new assert.Assertions object for use in tests
	ts.sqlmock = sql          // sql mock used to create expectations
	ts.acctService = NewAccountService(db)
	ts.transService = NewTransactionService(db, *ts.acctService)

	ts.account = database.Account{
		AccountHolder: "Foo Bar",
		AccountType:   "checking",
		Balance:       100,
	}

	// Set expectations on the mock to insert a new row into the accounts table
	ts.sqlmock.ExpectBegin()
	ts.sqlmock.ExpectExec("^INSERT INTO `accounts`").
		WithArgs(ts.account.AccountHolder, ts.account.AccountType, ts.account.Balance).
		WillReturnResult(sqlmock.NewResult(1, 1))
	ts.sqlmock.ExpectCommit()
}

func (ts *TransactionServiceSuite) TestCreate_ExecutesInsert() {
	// Set expectations on the mock for a SELECT query on the accounts table
	ts.sqlmock.ExpectBegin()
	ts.sqlmock.ExpectQuery("^SELECT (.+) FROM `accounts` WHERE `accounts`.`deleted_at` IS NULL AND \\(\\(`accounts`.`id` = \\?\\)\\) ORDER BY `accounts`.`id` ASC LIMIT 1").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "account_holder", "account_type", "balance"}).
			AddRow(1, time.Now(), time.Now(), nil, ts.account.AccountHolder, ts.account.AccountType, ts.account.Balance))
	ts.sqlmock.ExpectCommit()

	// Create a new transaction object pointer
	transaction := &database.Transaction{
		AccountID: 1,
		Type:      "deposit",
		Amount:    100,
		Account:   &ts.account,
	}

	// Set expectations on the mock for an INSERT query on the transactions table
	ts.sqlmock.ExpectBegin()
	ts.sqlmock.ExpectExec("^INSERT INTO `transactions`").
		WithArgs(mock.Any{}, mock.Any{}, mock.Any{}, mock.Any{}, 1, mock.Any{}, "deposit", 100).
		WillReturnResult(sqlmock.NewResult(1, 1))
	ts.sqlmock.ExpectCommit()

	// Call the method under test
	err := ts.transService.Create(transaction)

	// Check the result
	ts.assert.NoError(err)
	ts.assert.NoError(ts.sqlmock.ExpectationsWereMet())
}

//func (ts *TransactionServiceSuite) TestTransactionService_List() {
//	//TODO - implement
//}
//
//func (ts *TransactionServiceSuite) TestTransactionService_FetchById() {
//	//TODO - implement
//}
//
//func (ts *TransactionServiceSuite) TestTransactionService_Update() {
//	//TODO - implement
//}
//
//func (ts *TransactionServiceSuite) TestTransactionService_Delete() {
//	//TODO - implement
//}
