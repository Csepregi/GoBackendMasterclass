package http

import (
	"context"
	"github.com/google/uuid"
)

type Handler struct{
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

func Register(ctx context.Context, e EchoRouter, handler Handler) error {
    RegisterHandlers(e, NewStrictHandler(handler, nil))

    return nil
}
