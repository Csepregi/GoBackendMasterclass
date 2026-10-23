package http

import (
	"context"
	"fmt"
	"eats/backend/orders/adapters/db/dbmodels"
	"eats/backend/common"
	"eats/backend/common/shared"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	db *pgxpool.Pool
}

func NewHandler(
	db *pgxpool.Pool,
) Handler {
	if db == nil {
		panic("db cannot be nil")
	}

	return Handler{
		db: db,
	}
}

func (h Handler) RegisterCustomer(ctx context.Context, request RegisterCustomerRequestObject) (RegisterCustomerResponseObject, error) {
	customerUUID := common.NewUUIDv7()
	commonAddress, err := openapiAddressToSharedAddress(request.Body.Address)

	if err != nil {
		return nil, fmt.Errorf("failed to convert address: %w", err)
	}

	queries := dbmodels.New(h.db)

	err = queries.InsertCustomer(ctx, dbmodels.InsertCustomerParams{
		CustomerUuid: customerUUID,
		Name:         request.Body.Name,
		Email:        string(request.Body.Email),
		PhoneNumber:  request.Body.PhoneNumber,
		Address:      commonAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert customer: %w", err)
	}
	return RegisterCustomer201JSONResponse{
		CustomerUuid: customerUUID,
	}, nil
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

func Register(ctx context.Context, e EchoRouter, handler Handler) error {
	RegisterHandlers(e, NewStrictHandler(handler, nil))

	return nil
}
