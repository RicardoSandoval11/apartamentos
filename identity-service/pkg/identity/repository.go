package identity

import (
	"context"

	"github.com/RicardoSandoval11/apartamentos/identity-service/entities"
)

type Repository interface {
	GetByEmail(ctx context.Context, email string) (*entities.Identity, error)
	Create(ctx context.Context, identity entities.Identity) (*uint64, error)
	Update(ctx context.Context, key string, identity entities.Identity) (bool, error)
}
