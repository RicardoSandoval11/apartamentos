package apartment

import (
	"context"
	"encoding/json"
	"net/http"
)

func DecodeGetApartmentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	publicId := r.URL.Query().Get("publicId")
	return GetApartmentRequest{
		PublicId: publicId,
	}, nil
}

func EncodeGetApartmentResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, ok := response.(GetApartmentResponse); !ok {
		w.WriteHeader(http.StatusBadRequest)
	}
	return json.NewEncoder(w).Encode(response)
}
