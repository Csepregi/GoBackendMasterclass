package http

import (
	"context"
	"fmt"
	"eats/backend/orders/adapters/db/dbmodels"
	"eats/backend/common"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	db *pgxpool.Pool
}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) RegisterCustomer(ctx context.Context, request RegisterCustomerRequestObject) (RegisterCustomerResponseObject, error) {
	customerUUID := common.NewUUIDv7()
	commonAddress, err := openapiAddressToSharedAddress(customer.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to convert address: %w", err)
	}

	queries := dbmodels.NewQueries(h.db)

	err = queries.InsertCustomer(ctx, dbmodels.InsertCustomerParams{
		CustomerUuid: customerUUID,
		Name:         request.Body.Name,
		Email:        string(request.Body.Email),
		PhoneNumber:  customer.PhoneNumber,
		Address:      commonAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert customer: %w", err)
	}
	return RegisterCustomer201JSONResponse{
		CustomerUuid: customerUUID,
	}, nil
}

func Register(ctx context.Context, e EchoRouter, handler Handler) error {
	RegisterHandlers(e, NewStrictHandler(handler, nil))

	return nil
}
