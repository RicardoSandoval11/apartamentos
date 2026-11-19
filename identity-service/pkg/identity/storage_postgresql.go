package identity

import (
	"context"
	"errors"

	"github.com/RicardoSandoval11/apartamentos/identity-service/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresqlRepository struct {
	db *gorm.DB
}

func NewIdentityRepository(db *gorm.DB) Repository {
	return &postgresqlRepository{db: db}
}

func (r *postgresqlRepository) GetByEmail(ctx context.Context, email string) (*entities.Identity, error) {

	var identity entities.Identity

	result := r.db.WithContext(ctx).Where("email = ?", email).First(&identity)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &identity, nil
}

func (r *postgresqlRepository) Create(ctx context.Context, identity entities.Identity) (*uint64, error) {
	result := r.db.WithContext(ctx).Create(&identity)

	if result.Error != nil {
		return nil, result.Error
	}

	return &identity.Id, nil
}

func (r *postgresqlRepository) Update(ctx context.Context, key string, identity entities.Identity) (bool, error) {
	result := r.db.WithContext(ctx).Model(&identity).Select(key).Updates(identity)

	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

func (r *postgresqlRepository) GetByPublicId(ctx context.Context, publicId uuid.UUID) (*entities.Identity, error) {
	var identity entities.Identity

	result := r.db.WithContext(ctx).Where("public_id = ?", publicId).First(&identity)

	if result.Error != nil {
		return nil, result.Error
	}

	return &identity, nil
}
