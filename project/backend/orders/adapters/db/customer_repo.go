package db

import (
	"context"
	"fmt"
	"eats/backend/orders/adapters/db/dbmodels"
	"github.com/jackc/pgx/v5/pgxpool"
	"eats/backend/common/shared"
	"eats/backend/common"
	"eats/backend/orders/api/http"
)

type CustomerRepository struct {
	db *pgxpool.Pool
}

func NewCustomerRepository(db *pgxpool.Pool) *CustomerRepository {
	if db == nil {
		panic("db connection pool cannot be nil")
	}

	return &CustomerRepository{
		db: db,
	}
}

func (r *CustomerRepository) RegisterCustomer(ctx context.Context, customerUUID common.UUID, customer http.RegisterCustomer) error {

	commonAddress, err := openapiAddressToSharedAddress(customer.Address)

	if err != nil {
		return fmt.Errorf("failed to convert address: %w", err)
	}

	queries := dbmodels.New(r.db)

	err = queries.InsertCustomer(ctx, dbmodels.InsertCustomerParams{
		CustomerUuid: customerUUID,
		Name:         customer.Name,
		Email:        string(customer.Email),
		PhoneNumber:  customer.PhoneNumber,
		Address:  commonAddress,
	})
	if err != nil {
		return fmt.Errorf("failed to insert customer: %w", err)
	}
	return nil
}

func openapiAddressToSharedAddress(address http.Address) (shared.Address, error) {
	return shared.NewAddress(
		address.Line1,
		address.Line2,
		address.PostalCode,
		address.City,
		address.CountryCode,
	)
}
