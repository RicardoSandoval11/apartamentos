package apartment

import (
	"context"
	"errors"

	"github.com/RicardoSandoval11/apartamentos/apartments-service/entities"
	"github.com/google/uuid"
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
		return nil, errors.New("invalid public id")
	}

	if err := uuid.Validate(publicId); err != nil {
		return nil, errors.New("invalid uuid")
	}

	result, err := s.repository.GetById(ctx, publicId)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
