package verification

import (
	"context"

	"github.com/RicardoSandoval11/apartamentos/identity-service/entities"
)

type Repository interface {
	SaveVerification(ctx context.Context, identityId uint64) (string, error)
	IncreaseVerificationAttempts(ctx context.Context, verificationId uint64) error
	GetVerificationByIdentityId(ctx context.Context, identityId uint64) (*entities.Verification, error)
	GetAttemptsByVerificationId(ctx context.Context, verificationId uint64) (*entities.Attempt, error)
	MarkVerificationAsVerified(ctx context.Context, verificationId uint64) error
}
