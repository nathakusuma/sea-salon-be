package repository

import (
	"github.com/nathakusuma/sea-salon-be/internal/pkg/entity"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type IBranchRepository interface {
	Create(branch *entity.Branch) (ulid.ULID, error)
	FindAll() ([]entity.Branch, error)
	FindByID(id ulid.ULID) (entity.Branch, error)
	SetServices(branchID ulid.ULID, serviceID []ulid.ULID) error
}

type branchRepository struct {
	db *gorm.DB
}

func NewBranchRepository(db *gorm.DB) IBranchRepository {
	return &branchRepository{db: db}
}

func (b *branchRepository) Create(branch *entity.Branch) (ulid.ULID, error) {
	if err := b.db.Debug().Create(branch).Error; err != nil {
		return ulid.ULID{}, err
	}

	return branch.ID, nil
}

func (b *branchRepository) FindAll() ([]entity.Branch, error) {
	var branches []entity.Branch

	tx := b.db.Debug().Preload("Services")

	tx.Find(&branches)

	return branches, tx.Error
}

func (b *branchRepository) FindByID(id ulid.ULID) (entity.Branch, error) {
	var branch entity.Branch

	tx := b.db.Debug().Preload("Services").Where("id = ?", id)

	tx.First(&branch)

	return branch, tx.Error
}

func (b *branchRepository) SetServices(branchID ulid.ULID, serviceID []ulid.ULID) error {
	branch := entity.Branch{ID: branchID}
	if err := b.db.Debug().Preload("Services").First(&branch).Error; err != nil {
		return err
	}

	services := make([]entity.Service, len(serviceID))
	for i, id := range serviceID {
		services[i] = entity.Service{ID: id}
	}

	if err := b.db.Debug().Model(&branch).Association("Services").Replace(services); err != nil {
		return err
	}

	return nil
}
