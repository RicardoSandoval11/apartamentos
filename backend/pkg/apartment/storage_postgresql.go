package apartment

import (
	"context"
	"errors"

	"github.com/RicardoSandoval11/apartamentos/backend/entities"
	"gorm.io/gorm"
)

type postgresqlRepository struct {
	db *gorm.DB
}

func NewPostgresqlRepository(db *gorm.DB) Repository {
	return &postgresqlRepository{db: db}
}

func (r *postgresqlRepository) GetById(ctx context.Context, publicId string) (entities.Apartment, error) {
	var result entities.Apartment

	err := r.db.WithContext(ctx).Where("public_id = ?", publicId).First(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Apartment{}, errors.New("apartment not found")
		}

		return entities.Apartment{}, err
	}

	return result, nil
}
