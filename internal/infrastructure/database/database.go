package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(conf *DatabaseConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch conf.DbType {
	case "mysql":
		db, err = gorm.Open(mysql.Open(conf.DSN()), &gorm.Config{})
	case "postgres":
		db, err = gorm.Open(postgres.Open(conf.DSN()), &gorm.Config{})
	default:
		log.Fatalf("unsupported database type: %s", conf.DbType)
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}
