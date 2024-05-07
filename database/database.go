package database

import (
	"fmt"
	"log"
	"strconv"

	gorm "github.com/jinzhu/gorm"

	// Get the guts for postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

// Hold on to the database connection for all queries
var db *gorm.DB

// Hold on to the error
var err error

// Init creates a connection to a postgres database and returns the DB instance
func Init(user string, password string, host string, port int, database string) {
	dbinfo := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		user,
		password,
		host,
		strconv.Itoa(port),
		database,
	)
	Connect(dbinfo)
}

// GetDatabase will return the database instance
func GetDatabase() *gorm.DB {
	return db
}

// Connect will connect to the database and return a DB instance
func Connect(dbinfo string) {
	db, err = gorm.Open("postgres", dbinfo)
	if err != nil {
		log.Println("Failed to connect to database.")
		// panic(err)
	}
	log.Println("Database connected")
}

// BuildDatabase sets up the database if there are no tables
func BuildDatabase() {
	if !db.HasTable(SampleEntity{}) {
		err := db.CreateTable(SampleEntity{})
		if err != nil {
			log.Println("Table already exists")
		}
	}
	db.AutoMigrate(SampleEntity{})
}

// Close closes the database connection
func Close() {
	db.Close()
}
