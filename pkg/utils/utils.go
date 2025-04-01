package utils

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
)

func ValidateDecimal(value decimal.Decimal, precision int) error {
	str := value.String()
	if strings.Contains(str, ".") {
		length := len(strings.Split(str, ".")[1])
		if length > precision {
			return fmt.Errorf("incorrect precision: %s", str)
		}
	}

	return nil
}
