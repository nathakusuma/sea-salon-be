package entity

import (
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	ID           ulid.ULID `gorm:"primaryKey; not null; unique"`
	CustomerName string    `gorm:"type:varchar(50); not null"`
	StarRating   int       `gorm:"type:smallint; not null"`
	Comment      string    `gorm:"type:varchar(255); not null"`
}
