package handlers

import (
	"encoding/json"
	"gw-currency-wallet/internal/storages"
	"gw-currency-wallet/pkg/logger"
	"net/http"
)

func HandleRegister(db storages.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Email    string `json:"email"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if db.Register(body.Username, body.Password, body.Email) != nil {
			jsonBytes, _ := json.Marshal(
				struct {
					Error string `json:"error"`
				}{Error: "Username or email already exists"})

			logger.Log.Error("Error on register", "error", string(jsonBytes))
			http.Error(w, string(jsonBytes), http.StatusBadRequest)
			return
		}

		jsonBytes, _ := json.Marshal(
			struct {
				Message string `json:"message"`
			}{Message: "User registered successfully"})

		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)

		logger.Log.Debug("User registered successfully", "user", body.Username)
	}
}

func HandleLogin(db storages.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		token, err := db.Login(body.Username, body.Password)
		if err != nil {
			jsonBytes, _ := json.Marshal(
				struct {
					Error string `json:"error"`
				}{Error: "Invalid username or password"})

			logger.Log.Error("Error on login", "error", string(jsonBytes))
			http.Error(w, string(jsonBytes), http.StatusUnauthorized)
			return
		}

		jsonBytes, _ := json.Marshal(
			struct {
				Token string `json:"token"`
			}{Token: token})

		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)

		logger.Log.Debug("User login successfully", "user", body.Username)
	}
}
