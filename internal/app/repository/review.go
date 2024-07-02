package repository

import (
	"github.com/nathakusuma/sea-salon-be/internal/pkg/entity"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type IReviewRepository interface {
	Create(review *entity.Review) (ulid.ULID, error)
	FindByLazyLoad(idPivot ulid.ULID, action string, limit int) ([]entity.Review, error)
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) IReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(review *entity.Review) (ulid.ULID, error) {
	if err := r.db.Create(review).Error; err != nil {
		return ulid.ULID{}, err
	}

	return review.ID, nil
}

func (r *reviewRepository) FindByLazyLoad(idPivot ulid.ULID, action string, limit int) ([]entity.Review, error) {
	var reviews []entity.Review

	tx := r.db.Debug().Preload("User").Limit(limit)

	switch {
	case action == "top":
		tx = tx.Where("star_rating = 5").Order("id DESC")
	case idPivot != (ulid.ULID{}):
		if action == "prev" {
			tx = tx.Where("id > ?", idPivot).Order("id ASC")
		} else {
			tx = tx.Where("id < ?", idPivot).Order("id DESC")
		}
	default:
		tx = tx.Order("id DESC")
	}

	tx.Find(&reviews)

	if action == "prev" {
		for i, j := 0, len(reviews)-1; i < j; i, j = i+1, j-1 {
			reviews[i], reviews[j] = reviews[j], reviews[i]
		}
	}

	return reviews, tx.Error
}
