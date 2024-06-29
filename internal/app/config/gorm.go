package config

import (
	"github.com/nathakusuma/sea-salon-be/internal/pkg/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func NewDatabase() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	if err := migrateTables(db); err != nil {
		log.Fatalln(err)
	}

	return db
}

func migrateTables(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&entity.User{},
		&entity.Review{},
		&entity.Reservation{},
	); err != nil {
		return err
	}

	return nil
}
