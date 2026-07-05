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
}

func MakeGetApartmentsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetApartmentRequest)
		apt, err := svc.GetApartment(ctx, req.PublicId)
		if err != nil {
			return GetApartmentResponse{
				Data:    nil,
				Success: false,
			}, err
		}

		return GetApartmentResponse{
			Data:    apt,
			Success: true,
		}, nil
	}
}
