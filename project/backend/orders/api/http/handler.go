package http

import (
	"context"
	"eats/backend/common"
	"eats/backend/orders/app"
	"fmt"
	"eats/backend/common/shared"
)


type Handler struct {
	service *app.Service
}

func NewHandler(
	service *app.Service,
) Handler {
	if service == nil {
		panic("service cannot be nil")
	}
	return Handler{
		service: service,
	}
}




func (h Handler) RegisterCustomer(ctx context.Context, request RegisterCustomerRequestObject) (RegisterCustomerResponseObject, error) {
	customerUUID := common.NewUUIDv7()
	addr, err := openapiAddressToSharedAddress(request.Body.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}
	err = h.service.RegisterCustomer(ctx, app.Customer{
		CustomerUUID: customerUUID,
		Name:         request.Body.Name,
		Email:        string(request.Body.Email),
		PhoneNumber:  request.Body.PhoneNumber,
		Address:      addr,
	})
	if err != nil {
		return nil, err
	}
	
	return RegisterCustomer201JSONResponse{
		CustomerUuid: customerUUID,
	}, nil
}


func Register(ctx context.Context, e EchoRouter, handler Handler) error {
	RegisterHandlers(e, NewStrictHandler(handler, nil))

	return nil
}


func openapiAddressToSharedAddress(address Address) (shared.Address, error) {
	return shared.NewAddress(
		address.Line1,
		address.Line2,
		address.PostalCode,
		address.City,
		address.CountryCode,
	)
}
