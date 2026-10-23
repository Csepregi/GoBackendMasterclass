package http

import (
	"context"
	uuid "github.com/google/uuid"
	"eats/backend/common"
)

type Handler struct{
	RegisterCustomer StrictServerInterface
}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) RegisterCustomer(ctx context.Context, request RegisterCustomerRequestObject) (RegisterCustomerResponseObject, error) {
    customerUUID := uuid.New()

    return RegisterCustomer201JSONResponse{
        CustomerUuid: customerUUID,
    }, nil
}

func Register(ctx context.Context, e common.EchoRouter, handler Handler) error {
    RegisterHandlers(e, NewStrictHandler(handler, nil))

    return nil
}
