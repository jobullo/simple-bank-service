package service

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jobullo/go-api-example/database"
	"github.com/jobullo/go-api-example/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// A new test suite is created by embedding
// the suite.Suite struct.
type AccountServiceSuite struct {
	suite.Suite
	assert  *assert.Assertions
	sqlmock sqlmock.Sqlmock
	service *AccountService
}

// Invoke this function to run the test suite with "go test" at the CLI
func TestAccountServiceSuite(t *testing.T) {
	suite.Run(t, new(AccountServiceSuite))
}

func (as *AccountServiceSuite) SetupTest() {
	t := as.T() //gets the current *testing.T context, used by testing framework to manage test state

	db, sql, err := mock.DB() //mock.DB() returns a gorm.DB and sqlmock.Sqlmock
	require.NoError(t, err)   //use require from testify to stop test if error is not nil

	as.assert = assert.New(t) //create a new assert.Assertions object for use in tests
	as.sqlmock = sql          // sql mock used to create expectations
	as.service = NewAccountService(db)
}

func (as *AccountServiceSuite) TestAccountService_Create() {
	// Create a new account object pointer
	account := &database.Account{
		AccountHolder: "Foo Bar",
		AccountType:   "savings",
		Balance:       100,
	}

	// Set expectations on the mock for an INSERT query on the transactions table
	as.sqlmock.ExpectBegin()
	as.sqlmock.ExpectExec("^INSERT INTO `accounts`").
		WithArgs(mock.Any{}, mock.Any{}, mock.Any{}, mock.Any{}, "Foo Bar", "deposit", 100).
		WillReturnResult(sqlmock.NewResult(1, 1))
	as.sqlmock.ExpectCommit()

	// Call the method under test
	err := as.service.Create(account)

	// Check the result
	as.assert.NoError(err)
	as.assert.NoError(as.sqlmock.ExpectationsWereMet())
}
