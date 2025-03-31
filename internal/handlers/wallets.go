package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gw-currency-wallet/internal/storages"
	"gw-currency-wallet/pkg/utils"
	"net/http"
)

func HandleGetWallet(db storages.WalletRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := utils.GetUserId(r)

		wallet, err := db.GetWalletByUser(uuid.MustParse(userId))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		res := struct {
			Balance interface{} `json:"balance"`
		}{
			Balance: struct {
				USD decimal.Decimal
				RUB decimal.Decimal
				EUR decimal.Decimal
			}{
				USD: wallet.UsdAmount,
				RUB: wallet.RubAmount,
				EUR: wallet.EurAmount,
			},
		}

		jsonBytes, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	}
}

func HandleDepositWallet(db storages.WalletRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := utils.GetUserId(r)

		body := struct {
			Amount   decimal.Decimal   `json:"amount"`
			Currency storages.Currency `json:"currency"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		wallet, err := db.GetWalletByUser(uuid.MustParse(userId))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		wallet, err = db.DepositAmount(wallet.WalletId, body.Currency, body.Amount)
		if err != nil {
			jsonBytes, _ := json.Marshal(
				struct {
					Error string `json:"error"`
				}{Error: "Invalid amount or currency"})

			http.Error(w, string(jsonBytes), http.StatusBadRequest)
			return
		}

		res := struct {
			Message string      `json:"message"`
			Balance interface{} `json:"new_balance"`
		}{
			Message: "Account topped up successfully",
			Balance: struct {
				USD decimal.Decimal
				RUB decimal.Decimal
				EUR decimal.Decimal
			}{
				USD: wallet.UsdAmount,
				RUB: wallet.RubAmount,
				EUR: wallet.EurAmount,
			},
		}

		jsonBytes, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	}
}

func HandleWithdrawWallet(db storages.WalletRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := utils.GetUserId(r)

		body := struct {
			Amount   decimal.Decimal   `json:"amount"`
			Currency storages.Currency `json:"currency"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		wallet, err := db.GetWalletByUser(uuid.MustParse(userId))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		wallet, err = db.WithdrawAmount(wallet.WalletId, body.Currency, body.Amount)
		if err != nil {
			jsonBytes, _ := json.Marshal(
				struct {
					Error string `json:"error"`
				}{Error: "Insufficient funds or invalid amount"})

			http.Error(w, string(jsonBytes), http.StatusBadRequest)
			return
		}

		res := struct {
			Message string      `json:"message"`
			Balance interface{} `json:"new_balance"`
		}{
			Message: "Withdrawal successful",
			Balance: struct {
				USD decimal.Decimal
				RUB decimal.Decimal
				EUR decimal.Decimal
			}{
				USD: wallet.UsdAmount,
				RUB: wallet.RubAmount,
				EUR: wallet.EurAmount,
			},
		}

		jsonBytes, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	}
}

func HandleExchangeRates(db storages.ExchangeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rates, err := db.GetExchangeRates()
		if err != nil {
			jsonBytes, _ := json.Marshal(
				struct {
					Error string `json:"error"`
				}{Error: "Failed to retrieve exchange rates"})

			http.Error(w, string(jsonBytes), http.StatusInternalServerError)
			return
		}

		res := struct {
			Rates interface{} `json:"rates"`
		}{
			Rates: struct {
				USD decimal.Decimal
				RUB decimal.Decimal
				EUR decimal.Decimal
			}{
				USD: decimal.NewFromFloat(float64(rates.Rates["USD"])),
				RUB: decimal.NewFromFloat(float64(rates.Rates["RUB"])),
				EUR: decimal.NewFromFloat(float64(rates.Rates["EUR"])),
			},
		}

		jsonBytes, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	}
}

func HandleExchangeWallet(db storages.WalletRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := utils.GetUserId(r)

		body := struct {
			Amount       decimal.Decimal   `json:"amount"`
			FromCurrency storages.Currency `json:"from_currency"`
			ToCurrency   storages.Currency `json:"to_currency"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		wallet, err := db.GetWalletByUser(uuid.MustParse(userId))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		wallet, exchangedAmount, err := db.ExchangeAmount(wallet.WalletId, body.FromCurrency, body.ToCurrency, body.Amount)
		if err != nil {
			jsonBytes, _ := json.Marshal(
				struct {
					Error string `json:"error"`
				}{Error: "Insufficient funds or invalid currencies"})

			http.Error(w, string(jsonBytes), http.StatusBadRequest)
			return
		}

		res := struct {
			Message         string          `json:"message"`
			ExchangedAmount decimal.Decimal `json:"exchanged_amount"`
			Balance         interface{}     `json:"new_balance"`
		}{
			Message:         "Exchange successful",
			ExchangedAmount: *exchangedAmount,
			Balance: struct {
				USD decimal.Decimal
				RUB decimal.Decimal
				EUR decimal.Decimal
			}{
				USD: wallet.UsdAmount,
				RUB: wallet.RubAmount,
				EUR: wallet.EurAmount,
			},
		}

		jsonBytes, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	}
}
