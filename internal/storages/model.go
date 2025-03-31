package storages

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Wallet struct {
	WalletId  uuid.UUID
	UserId    uuid.UUID
	RubAmount decimal.Decimal
	UsdAmount decimal.Decimal
	EurAmount decimal.Decimal
}

type Currency string

const (
	RUB Currency = "RUB"
	USD Currency = "USD"
	EUR Currency = "EUR"
)

func (c Currency) GetDBField() string {
	switch c {
	case RUB:
		return "rub_amount"
	case EUR:
		return "eur_amount"
	case USD:
		return "usd_amount"
	}
	return ""
}
