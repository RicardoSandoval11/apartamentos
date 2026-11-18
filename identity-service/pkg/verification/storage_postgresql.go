package verification

import (
	"context"
	"time"

	"github.com/RicardoSandoval11/apartamentos/identity-service/entities"
	randstr "go.rtnl.ai/x/randstr"
	"gorm.io/gorm"
)

const validChars = "0123456789"
const defaultExpiration = 60 * time.Second

type postgresqlRepository struct {
	db *gorm.DB
}

func NewVerificationRepository(db *gorm.DB) Repository {
	return &postgresqlRepository{
		db: db,
	}
}

func (r *postgresqlRepository) SaveVerification(ctx context.Context, identityId uint64) (string, error) {
	var generatedCode string = randstr.Generate(8, validChars)

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		verification := entities.Verification{
			Code:       generatedCode,
			Expiration: time.Now().Add(defaultExpiration),
			IdentityId: identityId,
			Verified:   false,
		}

		result := tx.Create(&verification)

		if result.Error != nil {
			return result.Error
		}

		attempt := entities.Attempt{
			VerificationId: verification.Id,
			Times:          0,
		}

		result = tx.Create(&attempt)

		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return generatedCode, nil
}

func (r *postgresqlRepository) IncreaseVerificationAttempts(ctx context.Context, verificationId uint64) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var attempt entities.Attempt

		result := tx.Where("verification_id = ?", verificationId).First(&attempt)

		if result.Error != nil {
			return result.Error
		}

		result = tx.Model(&entities.Attempt{}).Where("Id = ?", attempt.Id).Update("times", attempt.Times+1)

		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}

func (r *postgresqlRepository) GetVerificationByIdentityId(ctx context.Context, identityId uint64) (*entities.Verification, error) {
	var verification entities.Verification

	result := r.db.WithContext(ctx).Where("identity_id = ?", identityId).First(&verification)

	if result.Error != nil {
		return nil, result.Error
	}

	return &verification, nil
}

func (r *postgresqlRepository) GetAttemptsByVerificationId(ctx context.Context, verificationId uint64) (*entities.Attempt, error) {
	var attempt entities.Attempt

	result := r.db.WithContext(ctx).Where("verification_id = ?", verificationId).First(&attempt)

	if result.Error != nil {
		return nil, result.Error
	}

	return &attempt, nil
}

func (r *postgresqlRepository) MarkVerificationAsVerified(ctx context.Context, verificationId uint64) error {
	result := r.db.WithContext(ctx).Where("id = ?", verificationId).Update("", true)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
