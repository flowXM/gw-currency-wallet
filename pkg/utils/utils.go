package utils

import (
	"fmt"
	"github.com/shopspring/decimal"
	"os"
	"strconv"
	"strings"
)

const saltSize = 22
const secretKey = "b17a8f413497b3715b328edc6ab81b634d92b355eb279166f050b5d0a71bc552"

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func GetEnvUint16(key string, fallback uint16) uint16 {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	res, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return fallback
	}

	return uint16(res)
}

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
