package tests_test

import (
	"testing"
	"eats/backend/common/testutils"
	"github.com/stretchr/testify/assert"
)

func TestComponent_CriticalFlow(t *testing.T) {
	t.Parallel()

	clients := newTestClients()
	ctx := t.Context()
	country := testutils.GenerateRandomCountry()
	customerUUID := registerCustomerInCity(ctx, t, clients, country, testutils.GenerateRandomCity())
	customerUUID.assert(t, customerUUID)

}
