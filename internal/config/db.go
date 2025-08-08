package config

import (
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewMySQLDB() *sqlx.DB {
	db, err := Connect()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}

func Connect() (*sqlx.DB, error) {
	config := mysql.Config{
		User:                 os.Getenv("MYSQL_ROOT_USER"),
		Passwd:               os.Getenv("MYSQL_ROOT_PASSWORD"),
		Addr:                 os.Getenv("MYSQL_HOST"),
		DBName:               os.Getenv("MYSQL_DATABASE"),
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := sqlx.Connect("mysql", config.FormatDSN())
	return db, err
}
