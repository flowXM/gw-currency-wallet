package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gw-currency-wallet/internal/storages"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockWalletRepository struct {
	Wallets map[uuid.UUID]storages.Wallet
}

func (m mockWalletRepository) DepositAmount(walletId uuid.UUID, currency storages.Currency, amount decimal.Decimal) (*storages.Wallet, error) {
	wallet := m.Wallets[walletId]
	switch currency {
	case storages.RUB:
		wallet.RubAmount = wallet.RubAmount.Add(amount)
	case storages.USD:
		wallet.UsdAmount = wallet.UsdAmount.Add(amount)
	case storages.EUR:
		wallet.EurAmount = wallet.EurAmount.Add(amount)
	}

	m.Wallets[walletId] = wallet
	return &wallet, nil
}

func (m mockWalletRepository) WithdrawAmount(walletId uuid.UUID, currency storages.Currency, amount decimal.Decimal) (*storages.Wallet, error) {
	wallet := m.Wallets[walletId]
	switch currency {
	case storages.RUB:
		if wallet.RubAmount.LessThan(amount) {
			return nil, fmt.Errorf("")
		}
		wallet.RubAmount = wallet.RubAmount.Sub(amount)
	case storages.USD:
		wallet.UsdAmount = wallet.UsdAmount.Sub(amount)
	case storages.EUR:
		wallet.EurAmount = wallet.EurAmount.Sub(amount)
	}

	m.Wallets[walletId] = wallet
	return &wallet, nil
}

func (m mockWalletRepository) GetWalletByUser(userId uuid.UUID) (*storages.Wallet, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockWalletRepository) ExchangeAmount(walletId uuid.UUID, fromCurrency, toCurrency storages.Currency, amount decimal.Decimal) (*storages.Wallet, *decimal.Decimal, error) {
	//TODO implement me
	panic("implement me")
}

func NewWalletRepository() storages.WalletRepository {
	wallets := make(map[uuid.UUID]storages.Wallet)

	id, err := uuid.Parse("53b7e416-b5b7-4c75-a8da-920db89d8c65")
	if err != nil {
		panic(err)
	}
	wallets[id] = storages.Wallet{WalletId: id, UserId: id, RubAmount: decimal.NewFromInt(100), EurAmount: decimal.NewFromInt(200), UsdAmount: decimal.NewFromInt(300)}

	return &mockWalletRepository{
		Wallets: wallets,
	}
}

//func (w *mockWalletRepository) UpdateAmount(walletId uuid.UUID, amount decimal.Decimal) error {
//	wallet, exists := w.Wallets[walletId]
//	if !exists {
//		return errors.New("wallet not found")
//	}
//
//	wallet.Amount = amount
//	w.Wallets[walletId] = wallet
//	return nil
//}
//
//func (w *mockWalletRepository) GetAmount(walletId uuid.UUID) (decimal.Decimal, error) {
//	wallet, exists := w.Wallets[walletId]
//	if !exists {
//		return decimal.Decimal{}, errors.New("wallet not found")
//	}
//	return wallet.Amount, nil
//}

func TestGetWallet(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		handler := HandleGetWallet(NewWalletRepository())

		req, _ := http.NewRequest("GET", "/api/v1/wallets/{id}", nil)
		req.SetPathValue("id", "5b5220ed-8291-4f19-800c-2bcdf5edbc68")

		r := httptest.NewRecorder()
		handler.ServeHTTP(r, req)

		if status := r.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("unsuccessful retrieval of not existing wallet", func(t *testing.T) {
		handler := HandleGetWallet(NewWalletRepository())

		req, _ := http.NewRequest("GET", "/api/v1/wallets/{id}", nil)
		req.SetPathValue("id", "5d051b24-5ace-48cf-a369-faea1883c117")

		r := httptest.NewRecorder()
		handler.ServeHTTP(r, req)

		if status := r.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	t.Run("unsuccessful retrieval of incorrect UUID", func(t *testing.T) {
		handler := HandleGetWallet(NewWalletRepository())

		req, _ := http.NewRequest("GET", "/api/v1/wallets/{id}", nil)
		req.SetPathValue("id", "qwerty12345")

		r := httptest.NewRecorder()
		handler.ServeHTTP(r, req)

		if status := r.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}

func TestPostWallet(t *testing.T) {
	t.Run("successful send", func(t *testing.T) {
		walletRep := NewWalletRepository()
		handler := HandlePostWallet(walletRep)

		id, err := uuid.Parse("53b7e416-b5b7-4c75-a8da-920db89d8c65")
		if err != nil {
			t.Errorf("incorrect uuid")
		}

		model := requestmodels.Transaction{WalletId: id, OperationType: "DEPOSIT", Amount: decimal.NewFromInt(100)}
		requestByte, err := json.Marshal(model)
		if err != nil {
			t.Errorf("incorrect body")
		}

		req, _ := http.NewRequest("POST", "/api/v1/wallet", bytes.NewReader(requestByte))

		r := httptest.NewRecorder()
		handler.ServeHTTP(r, req)

		if status := r.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		amount, _ := walletRep.GetAmount(id)

		if amount.Cmp(decimal.NewFromInt(200)) == -1 {
			t.Errorf("incorrect amount")
		}
	})

	t.Run("successful send", func(t *testing.T) {
		walletRep := NewWalletRepository()
		handler := HandlePostWallet(walletRep)

		id, err := uuid.Parse("53b7e416-b5b7-4c75-a8da-920db89d8c65")
		if err != nil {
			t.Errorf("incorrect uuid")
		}

		model := requestmodels.Transaction{WalletId: id, OperationType: "WITHDRAW", Amount: decimal.NewFromInt(100)}
		requestByte, err := json.Marshal(model)
		if err != nil {
			t.Errorf("incorrect body")
		}

		req, _ := http.NewRequest("POST", "/api/v1/wallet", bytes.NewReader(requestByte))

		r := httptest.NewRecorder()
		handler.ServeHTTP(r, req)

		if status := r.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		amount, _ := walletRep.GetAmount(id)

		if amount.Cmp(decimal.NewFromInt(0)) == -1 {
			t.Errorf("incorrect amount")
		}
	})

	t.Run("unsuccessful send: not existing wallet", func(t *testing.T) {
		walletRep := NewWalletRepository()
		handler := HandlePostWallet(walletRep)

		id, err := uuid.Parse("e3c24d9e-958d-462f-b225-c62c278846de")
		if err != nil {
			t.Errorf("incorrect uuid")
		}

		model := requestmodels.Transaction{WalletId: id, OperationType: "DEPOSIT", Amount: decimal.NewFromInt(100)}
		requestByte, err := json.Marshal(model)
		if err != nil {
			t.Errorf("incorrect body")
		}

		req, _ := http.NewRequest("POST", "/api/v1/wallet", bytes.NewReader(requestByte))

		r := httptest.NewRecorder()
		handler.ServeHTTP(r, req)

		if status := r.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	t.Run("unsuccessful send: not valid wallet id", func(t *testing.T) {
		walletRep := NewWalletRepository()
		handler := HandlePostWallet(walletRep)

		model := struct {
			WalletId      string          `json:"walletId"`
			OperationType string          `json:"operationType"`
			Amount        decimal.Decimal `json:"amount"`
		}{
			WalletId: "abc", OperationType: "DEPOSIT", Amount: decimal.NewFromInt(100),
		}

		requestByte, err := json.Marshal(model)
		if err != nil {
			t.Errorf("incorrect body")
		}

		req, _ := http.NewRequest("POST", "/api/v1/wallet", bytes.NewReader(requestByte))

		r := httptest.NewRecorder()
		handler.ServeHTTP(r, req)

		if status := r.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("unsuccessful send: not valid operation type", func(t *testing.T) {
		walletRep := NewWalletRepository()
		handler := HandlePostWallet(walletRep)

		model := struct {
			WalletId      string          `json:"walletId"`
			OperationType string          `json:"operationType"`
			Amount        decimal.Decimal `json:"amount"`
		}{
			WalletId: "53b7e416-b5b7-4c75-a8da-920db89d8c65", OperationType: "NOT VALID OPERATION TYPE", Amount: decimal.NewFromInt(100),
		}

		requestByte, err := json.Marshal(model)
		if err != nil {
			t.Errorf("incorrect body")
		}

		req, _ := http.NewRequest("POST", "/api/v1/wallet", bytes.NewReader(requestByte))

		r := httptest.NewRecorder()
		handler.ServeHTTP(r, req)

		if status := r.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}
