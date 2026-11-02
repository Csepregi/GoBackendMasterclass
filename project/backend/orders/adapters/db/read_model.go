package db

import (
	"context"

	"eats/backend/orders/adapters/db/dbmodels"
	"eats/backend/orders/api/http"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ReadModel provides read-optimized queries that return HTTP response types directly.
type ReadModel struct {
	db *pgxpool.Pool
}

func NewReadModel(db *pgxpool.Pool) *ReadModel {
	if db == nil {
		panic("db connection pool cannot be nil")
	}
	return &ReadModel{db: db}
}

// TODO: Implement this method using sqlc to query menu items joined with restaurants.
func (r ReadModel) ListMenuItemsWithRestaurant(ctx context.Context, filter http.ListMenuItemsFilter) ([]http.MenuItemWithRestaurant, error) {
	queries := dbmodels.New(r.db)
	rows, err := queries.ListMenuItems(ctx, dbmodels.ListMenuItemsParams{
		SearchTerm:           filter.Search,
		RestaurantNameFilter: filter.RestaurantName,
		OrderBy:              filter.OrderBy,
	})
	if err != nil {
		return nil, err
	}

	items := make([]http.MenuItemWithRestaurant, 0, len(rows))
	for _, row := range rows {
		items = append(items, http.MenuItemWithRestaurant{
			RestaurantName: row.RestaurantName,
			MenuItemUuid:   row.MenuItemUuid,
			MenuItemName:   row.MenuItemName,
			Currency:       row.Currency,
			RestaurantUuid: row.RestaurantUuid,
			GrossPrice:     row.GrossPrice,
		})
	}
	return items, nil
}
