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
	err := r.db.InsertCustomer(ctx, r.db.InsertCustomerParams{
		CustomerUuid: customerUUID,
		Name:         customer.Name,
		Email:        string(customer.Email),
		PhoneNumber:  customer.PhoneNumber,
		Address:      customer.sharedAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert customer: %w", err)
	}
	return nil
}
