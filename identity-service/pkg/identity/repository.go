package identity

import (
	"context"

	"github.com/RicardoSandoval11/apartamentos/identity-service/entities"
	"github.com/google/uuid"
)

type Repository interface {
	GetByEmail(ctx context.Context, email string) (*entities.Identity, error)
	GetByPublicId(ctx context.Context, publicId uuid.UUID) (*entities.Identity, error)
	Create(ctx context.Context, identity entities.Identity) (*uint64, error)
	Update(ctx context.Context, key string, identity entities.Identity) (bool, error)
}
