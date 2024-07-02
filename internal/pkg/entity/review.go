package entity

import (
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	ID         ulid.ULID `gorm:"primaryKey; not null; unique"`
	UserID     ulid.ULID `gorm:"not null"`
	User       User
	StarRating int    `gorm:"type:smallint; not null"`
	Comment    string `gorm:"type:varchar(255); not null"`
}
