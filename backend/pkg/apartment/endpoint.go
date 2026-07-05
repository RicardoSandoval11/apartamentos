package apartment

import (
	"context"

	"github.com/RicardoSandoval11/apartamentos/backend/entities"
	"github.com/go-kit/kit/endpoint"
)

type GetApartmentRequest struct {
	PublicId string `json:"publicId"`
}

type GetApartmentResponse struct {
	Data    *entities.Apartment `json:"data"`
	Success bool                `json:"success"`
	Errors  []string            `json:"errors,omitempty"`
}

func MakeGetApartmentsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetApartmentRequest)
		apt, err := svc.GetApartment(ctx, req.PublicId)
		if err != nil {
			return GetApartmentResponse{
				Data:    nil,
				Success: false,
				Errors:  []string{err.Error()},
			}, err
		}

		return GetApartmentResponse{
			Data:    apt,
			Success: true,
			Errors:  []string{},
		}, nil
	}
}
