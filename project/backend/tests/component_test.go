package tests_test

import (
	"testing"
	"eats/backend/common/testutils"
	"github.com/stretchr/testify/assert"
)

func TestComponent_CriticalFlow(t *testing.T) {
	t.Parallel()

	clients := newTestClients(t)
	ctx := t.Context()
	country := testutils.GenerateRandomCountry()
	customerUUID := registerCustomerInCity(ctx, t, clients, country, "San Francisco")
	assert.NotEmpty(t, customerUUID)

}
