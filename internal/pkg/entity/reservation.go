package entity

import (
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"time"
)

type Reservation struct {
	gorm.Model
	ID          ulid.ULID `gorm:"primaryKey; not null; unique"`
	UserID      ulid.ULID `gorm:"not null"`
	ServiceName string    `gorm:"type:varchar(50); not null"`
	StartTime   time.Time `gorm:"type:timestamp; not null"`
}
