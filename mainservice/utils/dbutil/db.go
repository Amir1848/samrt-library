package dbutil

import (
	"context"
	"errors"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDBConnection(ctx context.Context) (*gorm.DB, error) {
	dbConnection := os.Getenv("DB")

	db, err := gorm.Open(postgres.Open(dbConnection), &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to connect to the database")
	}

	return db, nil
}

func CreateDatabaseTransaction(db *gorm.DB, transactionFunc func(*gorm.DB) error) error {
	tx := db.Begin()

	err := transactionFunc(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
