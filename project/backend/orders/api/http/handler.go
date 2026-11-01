package http

import (
	"context"
	"fmt"

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
	var menuItems []app.MenuItem
	for _, item := range request.Body.MenuItems {
		menuItems = append(menuItems, app.MenuItem{
			MenuItemUUID: item.Uuid,
			Name:         item.Name,
			GrossPrice:   item.GrossPrice,
			Ordering:     float64(item.Ordering),
		})
	}
	addr, err := openapiAddressToSharedAddress(request.Body.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}
	restaurant := app.OnboardRestaurant{
		Name:        request.Body.Name,
		Description: request.Body.Description,
		MenuItems:   menuItems,
		Address:     addr,
		Currency:    request.Body.Currency,
	}

	err = h.service.OnboardRestaurant(ctx, request.RestaurantUuid, restaurant)
	if err != nil {
		return nil, fmt.Errorf("onboarding restaurant failed: %w", err)
	}

	return OnboardRestaurant204Response{}, nil
}

func (h Handler) CustomerCreateQuote(ctx context.Context, request CustomerCreateQuoteRequestObject) (CustomerCreateQuoteResponseObject, error) {
	var quoteItems []app.CreateQuoteItem
	for _, item := range request.Body.Items {
		quoteItems = append(quoteItems, app.CreateQuoteItem{
			MenuItemUUID: item.MenuItemUuid,
			Quantity:     item.Quantity,
		})
	}


	quote := app.CreateQuote{
		CustomerUUID: request.Params.CustomerUUID,
		RestaurantUUID: request.Body.RestaurantUuid,
		QuoteItems: quoteItems,
		DeliveryAddress: request.Body.DeliveryAddress,
	}
	quote, err := h.service.CreateQuote(ctx, quote)
	if err != nil {
		return nil, fmt.Errorf("creating quote failed: %w", err)
	}
	return CustomerCreateQuote201JSONResponse{
		QuoteUuid: app.QuoteUUID{UUID: quote.QuoteUUID},
		ItemsSubtotalGross: decimal.NewFromInt(quote.ItemsSubtotalGross),
		ServiceFeeGross: decimal.NewFromInt(quote.ServiceFeeGross),
		DeliveryFeeGross: decimal.NewFromInt(quote.DeliveryFeeGross),
		TotalAmountGross: decimal.NewFromInt(quote.TotalAmountGross),
		TotalTax: decimal.NewFromInt(quote.TotalTax),
		ExpiresAt: quote.ExpiresAt,
		Currency: shared.Currency(quote.Currency),
	}, nil
}