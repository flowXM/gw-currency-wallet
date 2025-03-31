package storages

import (
	proto_exchange "github.com/flowXM/proto-exchange/exchange"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type WalletRepository interface {
	DepositAmount(walletId uuid.UUID, currency Currency, amount decimal.Decimal) (*Wallet, error)
	WithdrawAmount(walletId uuid.UUID, currency Currency, amount decimal.Decimal) (*Wallet, error)
	GetWalletByUser(userId uuid.UUID) (*Wallet, error)
	ExchangeAmount(walletId uuid.UUID, fromCurrency, toCurrency Currency, amount decimal.Decimal) (*Wallet, *decimal.Decimal, error)
}

type ExchangeRepository interface {
	GetExchangeRates() (*proto_exchange.ExchangeRatesResponse, error)
	GetExchangeRateForCurrency(from, to Currency) (*proto_exchange.ExchangeRateResponse, error)
}

type UserRepository interface {
	Register(username, password, email string) error
	Login(username, password string) (token string, err error)
}
