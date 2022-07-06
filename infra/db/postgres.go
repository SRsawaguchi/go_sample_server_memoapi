package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var psqlDB *gorm.DB = nil

func NewPostgresDB(dsn string) (*gorm.DB, error) {
	if psqlDB != nil {
		return psqlDB, nil
	}

	var err error
	psqlDB, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}

	return psqlDB, nil
}
