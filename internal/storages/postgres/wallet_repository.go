package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gw-currency-wallet/internal/storages"
	"gw-currency-wallet/internal/storages/grpc"
	"gw-currency-wallet/pkg/client/postgresql"
	"gw-currency-wallet/pkg/logger"
	"gw-currency-wallet/pkg/utils"
)

type walletRepository struct{}

func NewWalletRepository() storages.WalletRepository {
	return &walletRepository{}
}

func (w *walletRepository) CreateWallet(userId uuid.UUID) error {
	db, err := postgresql.NewClient()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO wallets VALUES ($1, $2, 100, 100, 100)", uuid.New(), userId)
	if err != nil {
		return err
	}

	return nil
}

func (w *walletRepository) GetWalletByUser(userId uuid.UUID) (*storages.Wallet, error) {
	db, err := postgresql.NewClient()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var wallet storages.Wallet

	result := db.QueryRow("SELECT * FROM wallets WHERE user_id = $1", userId)
	err = result.Scan(&wallet.WalletId, &wallet.UserId, &wallet.RubAmount, &wallet.UsdAmount, &wallet.EurAmount)
	if err != nil {
		return nil, err
	}

	logger.Log.Info("Successfully got amount on wallet", "Wallet", wallet)

	return &wallet, nil
}

func (w *walletRepository) DepositAmount(walletId uuid.UUID, currency storages.Currency, amount decimal.Decimal) (*storages.Wallet, error) {
	logger.Log.Debug("Trying update amount on wallet", "Wallet ID", walletId, "Amount", amount)

	if err := utils.ValidateDecimal(amount, 2); err != nil {
		return nil, err
	}

	if amount.IsNegative() {
		return nil, fmt.Errorf("amount is negative")
	}

	db, err := postgresql.NewClient()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var wallet storages.Wallet

	query := fmt.Sprintf("UPDATE wallets SET %s = %s + $1 WHERE wallet_id = $2 RETURNING *", currency.GetDBField(), currency.GetDBField())
	result := db.QueryRow(query, amount, walletId)
	err = result.Scan(&wallet.WalletId, &wallet.UserId, &wallet.RubAmount, &wallet.UsdAmount, &wallet.EurAmount)
	if err != nil {
		return nil, err
	}
	logger.Log.Info("Successfully updated amount on wallet", "Wallet ID", walletId, "Amount", amount)

	return &wallet, nil
}

func (w *walletRepository) WithdrawAmount(walletId uuid.UUID, currency storages.Currency, amount decimal.Decimal) (*storages.Wallet, error) {
	logger.Log.Debug("Trying update amount on wallet", "Wallet ID", walletId, "Amount", amount)

	if err := utils.ValidateDecimal(amount, 2); err != nil {
		return nil, err
	}

	if amount.IsNegative() {
		return nil, fmt.Errorf("amount is negative")
	}

	db, err := postgresql.NewClient()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var wallet storages.Wallet

	query := fmt.Sprintf("UPDATE wallets SET %s = %s - $1 WHERE wallet_id = $2 RETURNING *", currency.GetDBField(), currency.GetDBField())
	result := db.QueryRow(query, amount, walletId)
	err = result.Scan(&wallet.WalletId, &wallet.UserId, &wallet.RubAmount, &wallet.UsdAmount, &wallet.EurAmount)
	if err != nil {
		return nil, err
	}
	logger.Log.Info("Successfully updated amount on wallet", "Wallet ID", walletId, "Amount", amount)

	return &wallet, nil
}

func (w *walletRepository) ExchangeAmount(walletId uuid.UUID, fromCurrency, toCurrency storages.Currency, amount decimal.Decimal) (*storages.Wallet, *decimal.Decimal, error) {
	logger.Log.Debug("Trying exchange amount on wallet", "Wallet ID", walletId, "fromCurrency", fromCurrency, "toCurrency", toCurrency, "Amount", amount)

	if err := utils.ValidateDecimal(amount, 2); err != nil {
		return nil, nil, err
	}

	if amount.IsNegative() {
		return nil, nil, fmt.Errorf("amount is negative")
	}

	db, err := postgresql.NewClient()
	if err != nil {
		return nil, nil, err
	}
	defer db.Close()

	rr := grpc.NewExchangeRepository()
	rate, err := rr.GetExchangeRateForCurrency(fromCurrency, toCurrency)
	if err != nil {
		return nil, nil, err
	}

	logger.Log.Info("", rate)

	var wallet storages.Wallet

	exchangedAmount := amount.Mul(decimal.NewFromFloat(float64(rate.Rate)))

	query := fmt.Sprintf("UPDATE wallets SET %s = %s - $1, %s = %s + $2  WHERE wallet_id = $3 RETURNING *",
		fromCurrency.GetDBField(), fromCurrency.GetDBField(), toCurrency.GetDBField(), toCurrency.GetDBField())
	result := db.QueryRow(query, amount, exchangedAmount, walletId)
	err = result.Scan(&wallet.WalletId, &wallet.UserId, &wallet.RubAmount, &wallet.UsdAmount, &wallet.EurAmount)
	if err != nil {
		return nil, nil, err
	}

	logger.Log.Info("Successfully updated amount on wallet", "Wallet ID", walletId, "Amount", amount)

	return &wallet, &exchangedAmount, nil
}
