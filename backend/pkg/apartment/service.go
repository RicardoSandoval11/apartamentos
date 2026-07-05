package apartment

import (
	"context"
	"errors"

	"github.com/RicardoSandoval11/apartamentos/backend/entities"
	"github.com/google/uuid"
)

type Service interface {
	GetApartment(ctx context.Context, publicId string) (*entities.Apartment, error)
}

type apartmentService struct{}

func NewApartmentService() Service {
	return &apartmentService{}
}

func (s *apartmentService) GetApartment(ctx context.Context, publicId string) (*entities.Apartment, error) {
	if publicId == "" {
		return nil, errors.New("Invalid public id")
	}

	result := entities.Apartment{
		Id:       1,
		PublicId: uuid.New(),
		Title:    "Apartamento Zona 12",
	}

	return &result, nil
}
