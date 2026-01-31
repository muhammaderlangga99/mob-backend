package config

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewMySQLDB creates a GORM connection for the onboarding database.
// It is called from main.go to keep database wiring isolated from routing logic.
func NewMySQLDB(dsn string) (*gorm.DB, error) {
	// GORM handles connection pooling internally; we only provide DSN + config.
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("connected to MySQL")
	return db, nil
}
