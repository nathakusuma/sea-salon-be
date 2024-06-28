package repository

import (
	"github.com/nathakusuma/sea-salon-be/internal/pkg/entity"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(user *entity.User) (ulid.ULID, error)
	FindByEmail(email string) (entity.User, error)
	FindByID(id ulid.ULID) (entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *entity.User) (ulid.ULID, error) {
	if err := r.db.Create(user).Error; err != nil {
		return ulid.ULID{}, err
	}

	return user.ID, nil
}

func (r *userRepository) FindByID(id ulid.ULID) (entity.User, error) {
	var user entity.User
	if err := r.db.First(&user, id).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) FindByEmail(email string) (entity.User, error) {
	var user entity.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}
