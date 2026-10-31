package db

import (
	"context"
	"errors"
	"fmt"

	"eats/backend/common"
	"eats/backend/common/log"
	"eats/backend/orders/adapters/db/dbmodels"
	"eats/backend/orders/app"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/slices"
)

type RestaurantRepository struct {
	db *pgxpool.Pool
}

func NewRestaurantRepository(db *pgxpool.Pool) *RestaurantRepository {
	if db == nil {
		panic("db connection pool cannot be nil")
	}

	return &RestaurantRepository{
		db: db,
	}
}

func (r *RestaurantRepository) UpsertRestaurant(ctx context.Context, restaurantUUID app.RestaurantUUID, restaurant app.OnboardRestaurant) error {
	return common.UpdateInTx(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		queries := dbmodels.New(tx)

		log.FromContext(ctx).With("restaurant_uuid", restaurantUUID).Info("Upserting restaurant")

		dbRestaurant, err := queries.UpsertRestaurant(ctx, dbmodels.UpsertRestaurantParams{
			restaurantUUID,
			restaurant.Name,
			restaurant.Description,
			restaurant.Address,
			restaurant.Currency,
		})
		if err != nil {
			return fmt.Errorf("upsert restaurant failed: %w", err)
		}

		// Currency is immutable after creation - the upsert doesn't update it.
		// Check here catches attempts to change it and returns a clear error.
		if dbRestaurant.Currency != restaurant.Currency {
			return common.NewInvalidInputError("cannot-change-currency", "cannot change restaurant currency once set")
		}

		currentMenuItems, err := queries.GetRestaurantMenu(ctx, restaurantUUID)
		if err != nil {
			return fmt.Errorf("get restaurant menu failed: %w", err)
		}
		currentMenuItemsUUIDs := make([]app.RestaurantMenuItemUUID, 0, len(currentMenuItems))
		for _, item := range currentMenuItems {
			currentMenuItemsUUIDs = append(currentMenuItemsUUIDs, item.OrdersRestaurantMenuItem.RestaurantMenuItemUuid)
		}

		newUUIDs := make([]app.RestaurantMenuItemUUID, 0, len(restaurant.MenuItems))
		for _, item := range restaurant.MenuItems {
			newUUIDs = append(newUUIDs, item.MenuItemUUID)
		}

		menuItemsToArchive := make([]common.UUID, 0)
		for _, currentUUID := range currentMenuItemsUUIDs {
			if !slices.Contains(newUUIDs, currentUUID) {
				menuItemsToArchive = append(menuItemsToArchive, currentUUID.UUID)
			}
		}

		if len(menuItemsToArchive) > 0 {
			err := queries.ArchiveMenuItem(ctx, menuItemsToArchive)
			if err != nil {
				return fmt.Errorf("archive menu items failed: %w", err)
			}
		}

		for _, item := range restaurant.MenuItems {
			affectedRows, err := queries.UpsertRestaurantMenuItems(ctx, dbmodels.UpsertRestaurantMenuItemsParams{
				RestaurantMenuItemUuid: item.MenuItemUUID,
				RestaurantUuid:         restaurantUUID,
				Name:                   item.Name,
				GrossPrice:             item.GrossPrice,
				Ordering:               item.Ordering,
				IsArchived:             item.IsArchived,
			})
			if err != nil {
				return fmt.Errorf("upsert restaurant menu item failed: %w", err)
			}

			if affectedRows == 0 {
				return common.NewInvalidInputError(
					"menu-item-of-another-restaurant",
					"menu position s% belongs to another restaurant",
					item.MenuItemUUID,
				)
			}
		}

		return nil
	})
}

func (r *RestaurantRepository) GetRestaurantMenu(ctx context.Context, restaurantUUID app.RestaurantUUID) (app.RestaurantMenu, error) {
	queries := dbmodels.New(r.db)

	dbItems, err := queries.GetRestaurantMenu(ctx, restaurantUUID)
	if err != nil {
		return app.RestaurantMenu{}, fmt.Errorf("get restaurant menu failed: %w", err)
	}

	log.FromContext(ctx).With("restaurant_uuid", restaurantUUID, "count", len(dbItems)).Info("Fetched menu items")

	items := make([]app.MenuItem, len(dbItems))
	for i, dbItem := range dbItems {
		items[i] = app.MenuItem{
			MenuItemUUID: dbItem.OrdersRestaurantMenuItem.RestaurantMenuItemUuid,
			Name:         dbItem.OrdersRestaurantMenuItem.Name,
			GrossPrice:   dbItem.OrdersRestaurantMenuItem.GrossPrice,
			Ordering:     dbItem.OrdersRestaurantMenuItem.Ordering,
		}
	}

	restaurant, err := queries.GetRestaurant(ctx, restaurantUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return app.RestaurantMenu{}, common.NewNotFoundError("restaurant-not-found", "restaurant not found")
		}
		return app.RestaurantMenu{}, fmt.Errorf("get restaurant %s failed: %w", restaurantUUID, err)
	}

	return app.RestaurantMenu{
		RestaurantName: restaurant.Name,
		Address:        restaurant.Address,
		Description:    restaurant.Description,
		Currency:       restaurant.Currency,
		Positions:      items,
	}, nil
}
