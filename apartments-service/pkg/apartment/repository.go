package apartment

import (
	"context"

	"github.com/RicardoSandoval11/apartamentos/apartments-service/entities"
)

type Repository interface {
	GetById(ctx context.Context, publicId string) (entities.Apartment, error)
}
