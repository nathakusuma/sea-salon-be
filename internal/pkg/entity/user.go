package entity

import (
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           ulid.ULID `gorm:"primaryKey; not null; unique"`
	FullName     string    `gorm:"type:varchar(50); not null"`
	Email        string    `gorm:"type:varchar(320); not null; unique"`
	PhoneNumber  string    `gorm:"type:varchar(15); not null; unique"`
	Password     string    `gorm:"type:varchar(255); not null"`
	IsAdmin      bool      `gorm:"type:boolean; not null; default:false"`
	Reservations []Reservation
}
