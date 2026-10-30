package db

import (
	"context"
	"fmt"

	"eats/backend/orders/adapters/db/dbmodels"
	"eats/backend/orders/app"

	"eats/backend/common"

	"github.com/jackc/pgx/v5"
	
	"github.com/jackc/pgx/v5/pgxpool"
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

func (r *CustomerRepository) RegisterCustomer(ctx context.Context, customer app.Customer) error {
	return common.UpdateInTx(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		queries := dbmodels.New(tx)
		err := queries.InsertCustomer(ctx, dbmodels.InsertCustomerParams{
			CustomerUuid: customer.CustomerUUID,
			Name:         customer.Name,
			Email:        customer.Email,
			PhoneNumber:  customer.PhoneNumber,
			Address:      customer.Address,
		})
		if err != nil {
			return fmt.Errorf("failed to insert customer: %w", err)
		}
		return nil
	})
}
