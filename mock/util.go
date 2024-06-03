package mock

import (
	"errors"
	"math/rand"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-faker/faker/v4"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func DB() (*gorm.DB, sqlmock.Sqlmock, error) {
	conn, sql, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	db, err := gorm.Open("postgres", conn)
	if err != nil {
		return nil, nil, err
	}

	db.LogMode(false) // Disable logging

	return db, sql, nil
}

func Error() error {
	return errors.New(faker.Sentence())
}

func ID() uint {
	return uint(rand.Uint32())
}
