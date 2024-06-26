package repository

import (
	"github.com/nathakusuma/sea-salon-be/internal/pkg/entity"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"time"
)

type IReservationRepository interface {
	Create(reservation *entity.Reservation) (ulid.ULID, error)
	FindByTimeRange(serviceName string, start time.Time, end time.Time) ([]entity.Reservation, error)
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) IReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) Create(reservation *entity.Reservation) (ulid.ULID, error) {
	if err := r.db.Debug().Create(reservation).Error; err != nil {
		return ulid.ULID{}, err
	}

	return reservation.ID, nil
}

func (r *reservationRepository) FindByTimeRange(serviceName string, start time.Time, finish time.Time) ([]entity.Reservation, error) {
	var reservations []entity.Reservation

	tx := r.db.Debug().Where("service_name = ? AND start_time BETWEEN ? AND ?", serviceName, start, finish)

	tx.Find(&reservations)

	return reservations, tx.Error
}
