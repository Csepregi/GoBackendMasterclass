package shared

import (
	"fmt"

	"eats/backend/common"
)


type CurrencyType string

type Currency struct {
	common.Enum[CurrencyType]
}


func (c CurrencyType) Values() []string {
	return []string{
		"USD",
		"EUR",
		"GBP",
		"JPY",
		"PLN",
	}
}

func MustNewCurrencyCode(value string) Currency {
	c := Currency{}
	err := c.UnmarshalText([]byte(value))
	if err != nil {
		panic(fmt.Errorf("error unmarshalling currency code value: %s", value))
	}

	return c
}
