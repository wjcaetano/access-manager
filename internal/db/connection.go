package db

import (
	"access-manager/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDatabase(cfg config.Configuration) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil

}

func BuildDSN(cfg config.Configuration) string {
	return cfg.Database.DSN
}
