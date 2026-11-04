package app

import (
	"context"
	"strings"
	"eats/backend/common"
)

type CourierUUID struct {
	common.UUID
}

type RegisterCourier struct {
	Name        string
	PhoneNumber string
	City        string
}

type CourierRepository interface {
	RegisterCourier(ctx context.Context, courierUUID CourierUUID, courier RegisterCourier) (error)
}

func (s *Service) RegisterCourier(ctx context.Context, courier RegisterCourier) (CourierUUID, error) {
	errDetails := []common.ErrorDetails{}
	courierUUID := CourierUUID{common.NewUUIDv7()}

	if courierUUID.IsZero() {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "courier",
			EntityID:   "",
			ErrorSlug:  "empty-uuid",
			Message:    "UUID cannot be empty",
		})
	}

	if strings.TrimSpace(courier.Name) == "" {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "courier",
			EntityID:   courierUUID.String(),
			ErrorSlug:  "empty-name",
			Message:    "Name cannot be empty",
		})
	}
	
	if strings.TrimSpace(courier.PhoneNumber) == "" {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "courier",
			EntityID:   courierUUID.String(),
			ErrorSlug:  "empty-phone-number",
			Message:    "Phone number cannot be empty",
		})
	}

	if strings.TrimSpace(courier.City) == "" {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "courier",
			EntityID:   courierUUID.String(),
			ErrorSlug:  "empty-city",
			Message:    "City cannot be empty",
		})
	}

	if len(errDetails) > 0 {
		return CourierUUID{}, common.NewInvalidInputError("invalid-courier-data", "invalid courier data").WithDetails(errDetails)
	}

	err := s.courierRepository.RegisterCourier(ctx, courierUUID, courier)
	if err != nil {
		return CourierUUID{}, err
	}
	return courierUUID, nil
}
