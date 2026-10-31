package http

import (
	"context"
	"fmt"
	"maps"
	"github.com/shopspring/decimal"

	"eats/backend/common"
	"eats/backend/common/shared"
	"eats/backend/orders/app"
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
	customerUUID := app.CustomerUUID{UUID: common.NewUUIDv7()}
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


func (h Handler) OnboardRestaurant(ctx context.Context, request OnboardRestaurantRequestObject) (OnboardRestaurantResponseObject, error) {
	restaurant := app.OnboardRestaurant{
		Name:           request.Body.Name,
		Description:    request.Body.Description,
		MenuItems: func() []app.MenuItem {
			return maps.Map(request.Body.MenuItems, func(item MenuItem) app.MenuItem {
				menuItemUUID := app.RestaurantMenuItemUUID{UUID: common.NewUUIDv7()}
				return app.MenuItem{
					MenuItemUUID: menuItemUUID,
					Name:         item.Name,
					GrossPrice:   decimal.NewFromInt(item.GrossPrice),
					Ordering:     item.Ordering.float64(),
				}
			})
		}(),
		Address:        request.Body.Address,
		Currency:       request.Body.Currency,
	}

	err := h.service.OnboardRestaurant(ctx, restaurant)
	if err != nil {
		return nil, err
	}

	return OnboardRestaurant201JSONResponse{
	}, nil
}