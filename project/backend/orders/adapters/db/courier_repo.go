package db

import (
	"context"
	"fmt"

 	"eats/backend/common"
	"eats/backend/orders/adapters/db/dbmodels"
	"eats/backend/orders/app"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CourierRepository struct {
	db *pgxpool.Pool
}

func NewCourierRepository(db *pgxpool.Pool) *CourierRepository {
	if db == nil {
		panic("db connection pool cannot be nil")
	}

	return &CourierRepository{
		db: db,
	}
}

func (r *CourierRepository) RegisterCourier(ctx context.Context, courierUUID app.CourierUUID, courier app.RegisterCourier) error {
	return common.UpdateInTx(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		queries := dbmodels.New(tx)
		err := queries.AddCourier(ctx, dbmodels.AddCourierParams{
			CourierUuid: courierUUID.UUID,
			Name:        courier.Name,
			PhoneNumber: courier.PhoneNumber,
			City:        courier.City,
		})
		if err != nil {
			return fmt.Errorf("failed to insert courier: %w", err)
		}
		return nil
	})
}