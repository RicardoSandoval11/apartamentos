package apartment

import (
	"context"
	"errors"

	"github.com/RicardoSandoval11/apartamentos/backend/entities"
)

type Service interface {
	GetApartment(ctx context.Context, publicId string) (*entities.Apartment, error)
}

type apartmentService struct {
	repository Repository
}

func NewApartmentService(repository Repository) Service {
	return &apartmentService{
		repository: repository,
	}
}

func (s *apartmentService) GetApartment(ctx context.Context, publicId string) (*entities.Apartment, error) {
	if publicId == "" {
		return nil, errors.New("Invalid public id")
	}

	result, err := s.repository.GetById(ctx, publicId)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
