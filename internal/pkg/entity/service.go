package entity

import (
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	ID             ulid.ULID `gorm:"primaryKey; not null; unique"`
	Name           string    `gorm:"type:varchar(50); not null"`
	Description    string    `gorm:"type:varchar(255); not null"`
	Price          int       `gorm:"type:int; not null"`
	DurationMinute int       `gorm:"type:smallint; not null"`
}
