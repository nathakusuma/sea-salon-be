package repository

import (
	"github.com/nathakusuma/sea-salon-be/internal/pkg/entity"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type IServiceRepository interface {
	Create(service *entity.Service) (ulid.ULID, error)
	FindAll() ([]entity.Service, error)
	FindByID(id ulid.ULID) (entity.Service, error)
}

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) IServiceRepository {
	return &serviceRepository{db: db}
}

func (s *serviceRepository) Create(service *entity.Service) (ulid.ULID, error) {
	if err := s.db.Debug().Create(service).Error; err != nil {
		return ulid.ULID{}, err
	}

	return service.ID, nil
}

func (s *serviceRepository) FindAll() ([]entity.Service, error) {
	var services []entity.Service

	tx := s.db.Debug()

	tx.Find(&services)

	return services, tx.Error
}

func (s *serviceRepository) FindByID(id ulid.ULID) (entity.Service, error) {
	var service entity.Service

	tx := s.db.Debug().Where("id = ?", id)

	tx.First(&service)

	return service, tx.Error
}
