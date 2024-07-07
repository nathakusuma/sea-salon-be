package entity

import (
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"time"
)

type Branch struct {
	gorm.Model
	ID             ulid.ULID  `gorm:"primaryKey; not null; unique"`
	Name           string     `gorm:"type:varchar(50); not null"`
	Address        string     `gorm:"not null"`
	MapsURL        string     `gorm:"not null"`
	Phone          string     `gorm:"type:varchar(15); not null"`
	OpenTime       time.Time  `gorm:"type:timestamp; not null"`
	CloseTime      time.Time  `gorm:"type:timestamp; not null"`
	TimeZoneName   string     `gorm:"type:varchar(50); not null"`
	TimeZoneOffset int        `gorm:"type:int; not null"`
	Services       []*Service `gorm:"many2many:branch_services;"`
}
