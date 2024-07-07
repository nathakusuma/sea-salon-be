package repository

import (
	"github.com/nathakusuma/sea-salon-be/internal/pkg/entity"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"time"
)

type IReservationRepository interface {
	Create(reservation *entity.Reservation) (ulid.ULID, error)
	FindByTimeRange(branchID ulid.ULID, serviceID ulid.ULID, start time.Time, end time.Time) ([]entity.Reservation, error)
	FindByUserID(userID ulid.ULID) ([]entity.Reservation, error)
	FindByDateAndBranch(date time.Time, branchID ulid.ULID) ([]entity.Reservation, error)
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

func (r *reservationRepository) FindByTimeRange(branchID ulid.ULID, serviceID ulid.ULID, start time.Time, finish time.Time) ([]entity.Reservation, error) {
	var reservations []entity.Reservation

	tx := r.db.Debug().Where("branch_id = ? AND service_id = ? AND start_time BETWEEN ? AND ?", branchID, serviceID, start, finish)

	tx.Find(&reservations)

	return reservations, tx.Error
}

func (r *reservationRepository) FindByUserID(userID ulid.ULID) ([]entity.Reservation, error) {
	var reservations []entity.Reservation

	tx := r.db.Debug().Preload("Service").Preload("Branch").Where("user_id = ?", userID).Order("start_time DESC")

	tx.Find(&reservations)

	return reservations, tx.Error
}

func (r *reservationRepository) FindByDateAndBranch(date time.Time, branchID ulid.ULID) ([]entity.Reservation, error) {
	var reservations []entity.Reservation

	tx := r.db.Debug().Preload("Service").Preload("User").Preload("Branch").Where("DATE(start_time) = ? AND branch_id = ?", date, branchID).Order("start_time")

	tx.Find(&reservations)

	return reservations, tx.Error
}
