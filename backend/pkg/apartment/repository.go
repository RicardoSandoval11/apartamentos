package apartment

import (
	"context"

	"github.com/RicardoSandoval11/apartamentos/backend/entities"
)

type Repository interface {
	GetById(ctx context.Context, publicId string) (entities.Apartment, error)
}
