package service

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-faker/faker/v4"
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
		AccountHolder: faker.Name(),
		AccountType:   "savings",
		Balance:       100,
	}

	queryPattern := `(?i)INSERT INTO "accounts" \("created_at","updated_at","deleted_at","account_holder","account_type","balance"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\) RETURNING "accounts"\."id"`

	// Set expectations on the mock for an INSERT query on the transactions table
	as.sqlmock.ExpectBegin()
	as.sqlmock.ExpectQuery(queryPattern).
		WithArgs(mock.Any{}, mock.Any{}, mock.Any{}, account.AccountHolder, account.AccountType, account.Balance).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	as.sqlmock.ExpectCommit()

	// Call the method under test
	err := as.service.Create(account)

	// Check the result
	as.assert.NoError(err)
	as.assert.NoError(as.sqlmock.ExpectationsWereMet())
}

func (as *AccountServiceSuite) TestAccountService_FetchById() {
	id := mock.ID()
	rows := as.newRows()
	as.addRow(rows, id, time.Now(), time.Now(), nil, faker.Name(), "savings", 100)

	// Formulating the regex pattern
	queryPattern := `(?i)SELECT\s+\*\s+FROM\s+"accounts"\s+WHERE\s+"accounts"\."deleted_at"\s+IS\s+NULL\s+AND\s+\(\("accounts"\."id"\s+=\s+\d+\)\)\s+ORDER\s+BY\s+"accounts"\."id"\s+ASC\s+LIMIT\s+1`

	// Using sqlmock to expect the query
	as.sqlmock.ExpectQuery(queryPattern).WillReturnRows(rows)

	_, err := as.service.FetchById(id)

	if as.assert.NoError(err) {
		as.assert.NoError(as.sqlmock.ExpectationsWereMet())
	}
}

func (as *AccountServiceSuite) TestAccountService_FetchById_NotFound() {
	id := mock.ID()
	rows := as.newRows()

	// Formulating the regex pattern
	queryPattern := `(?i)SELECT\s+\*\s+FROM\s+"accounts"\s+WHERE\s+"accounts"\."deleted_at"\s+IS\s+NULL\s+AND\s+\(\("accounts"\."id"\s+=\s+\d+\)\)\s+ORDER\s+BY\s+"accounts"\."id"\s+ASC\s+LIMIT\s+1`

	// Using sqlmock to expect the query
	as.sqlmock.ExpectQuery(queryPattern).WillReturnRows(rows)

	_, err := as.service.FetchById(id)

	if as.assert.Error(err) {
		as.assert.NoError(as.sqlmock.ExpectationsWereMet())
	}
}

func (as *AccountServiceSuite) TestAccountService_Update() {
	id := mock.ID()
	originalAccountHolder := faker.Name()
	account := &database.Account{
		AccountHolder: originalAccountHolder,
		AccountType:   "savings",
		Balance:       100,
	}

	// set the schema for the account table in the mock database
	rows := as.newRows()
	// add a row to the account table
	as.addRow(rows, id, time.Now(), time.Now(), nil, account.AccountHolder, account.AccountType, account.Balance)

	// generate a new account holder name and set it on the account object
	newAccountHolder := faker.Name()
	account.AccountHolder = newAccountHolder

	// set expectations on the mock for a SELECT query on the accounts table
	as.sqlmock.ExpectQuery("^SELECT .* FROM \"accounts\".*").WillReturnRows(rows)

	// account exists, so we expect a transaction to update the account holder
	as.sqlmock.ExpectBegin()
	as.sqlmock.ExpectExec("UPDATE .*").WillReturnResult(sqlmock.NewResult(1, 1))
	as.sqlmock.ExpectCommit()

	// update with new account holder
	as.service.Update(account)

	// assert that all of data was updated correctly
	as.assert.NoError(as.sqlmock.ExpectationsWereMet())
	as.assert.Equal(account.AccountHolder, newAccountHolder)
	as.assert.Equal("savings", account.AccountType)
	as.assert.Equal(float64(100), account.Balance)
}

// test that the Update method returns an error when the account is not found
func (as *AccountServiceSuite) TestAccountService_Update_NotFound() {
	originalAccountHolder := faker.Name()
	account := &database.Account{
		AccountHolder: originalAccountHolder,
		AccountType:   "savings",
		Balance:       100,
	}

	//row not added so account won't be found
	rows := as.newRows()

	// select query will return no rows and no update query will be executed
	as.sqlmock.ExpectQuery("^SELECT .* FROM \"accounts\".*").WillReturnRows(rows)

	err := as.service.Update(account)

	if as.assert.Error(err) {
		as.assert.NoError(as.sqlmock.ExpectationsWereMet())
	}
}

func (as *AccountServiceSuite) TestAccountService_List() {
	// set the schema for the account table in the mock database
	rows := as.newRows()
	// add two rows to the account table
	as.addRow(rows, mock.ID(), time.Now(), time.Now(), nil, faker.Name(), "savings", 100)
	as.addRow(rows, mock.ID(), time.Now(), time.Now(), nil, faker.Name(), "checking", 200)

	// Formulating the regex pattern
	queryPattern := `^SELECT\s+\*\s+FROM\s+"accounts"\s+WHERE\s+"accounts"\."deleted_at"\s+IS\s+NULL`

	// Using sqlmock to expect the query
	as.sqlmock.ExpectQuery(queryPattern).WillReturnRows(rows)

	_, err := as.service.List()

	if as.assert.NoError(err) {
		as.assert.NoError(as.sqlmock.ExpectationsWereMet())
	}
}

//TODO: Fix the Delete test
//func (as *AccountServiceSuite) TestAccountService_Delete() {
//	id := mock.ID()
//	originalAccountHolder := faker.Name()
//	account := &database.Account{
//		AccountHolder: originalAccountHolder,
//		AccountType:   "savings",
//		Balance:       100,
//	}
//
//	// set the schema for the account table in the mock database
//	accountRows := as.newRows()
//	// add a row to the account table
//	as.addRow(accountRows, id, time.Now(), time.Now(), nil, account.AccountHolder, account.AccountType, account.Balance)
//
//	// set the schema for the transactions table in the mock database
//	transactionRows := as.newTransactionRows()
//	// add a row to the transactions table
//	as.addTransactionRow(transactionRows, mock.ID(), time.Now(), time.Now(), nil, id, 100, "deposit")
//
//	// set expectations on the mock for a SELECT query on the accounts table
//	as.sqlmock.ExpectQuery("^SELECT .* FROM \"accounts\".*").WillReturnRows(accountRows)
//
//	// expect a SELECT query on the transactions table to find associated transactions
//	as.sqlmock.ExpectQuery("^SELECT .* FROM \"transactions\".*").WillReturnRows(sqlmock.NewRows([]string{"id"}))
//
//	// expect a transaction to begin
//	as.sqlmock.ExpectBegin()
//
//	// expect transactions to be deleted
//	as.sqlmock.ExpectExec("DELETE .*").WillReturnResult(sqlmock.NewResult(1, 1))
//
//	// expect account to be deleted
//	as.sqlmock.ExpectExec("DELETE .*").WillReturnResult(sqlmock.NewResult(1, 1))
//
//	// expect the transaction to commit
//	as.sqlmock.ExpectCommit()
//
//	// delete the account
//	err := as.service.Delete(id)
//
//	// assert that all of data was updated correctly
//	as.assert.NoError(err)
//	as.assert.NoError(as.sqlmock.ExpectationsWereMet())
//}

// creates the rows object for use in tests
func (s *AccountServiceSuite) newRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "account_holder", "account_type", "balance"})
}

// creates the rows object for use in tests for the "transactions" table
func (s *AccountServiceSuite) newTransactionRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "account_id", "amount", "transaction_type"})
}

// populates the transaction rows object with a row of data
func (s *AccountServiceSuite) addTransactionRow(
	rows *sqlmock.Rows,
	id uint,
	createdAt time.Time,
	updatedAt time.Time,
	deletedAt *time.Time,
	accountID uint,
	amount float64,
	transactionType string,
) {
	rows.AddRow(id, createdAt, updatedAt, deletedAt, accountID, amount, transactionType)
}

// populates the rows object with a row of data
func (s *AccountServiceSuite) addRow(
	rows *sqlmock.Rows,
	id uint,
	createdAt time.Time,
	updatedAt time.Time,
	deletedAt *time.Time,
	accountHolder string,
	accountType string,
	balance float64,
) {
	rows.AddRow(id, createdAt, updatedAt, deletedAt, accountHolder, accountType, balance)
}
