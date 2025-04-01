package routes

import (
	"gw-currency-wallet/internal/handlers"
	"gw-currency-wallet/internal/storages/grpc"
	"gw-currency-wallet/internal/storages/postgres"
	"gw-currency-wallet/pkg/utils"
	"net/http"
	"reflect"
)

func Init(mux *http.ServeMux) {
	userRepository := postgres.NewUserRepository()
	walletRepository := postgres.NewWalletRepository()
	exchangeRepository := grpc.NewExchangeRepository()

	mux.Handle("POST /api/v1/register", utils.RateLimitedHandler(handlers.HandleRegister(userRepository)))
	mux.Handle("POST /api/v1/login", utils.RateLimitedHandler(handlers.HandleLogin(userRepository)))

	mux.Handle("GET /api/v1/balance", utils.RateLimitedHandler(utils.AuthHandler(handlers.HandleGetWallet(walletRepository))))

	mux.Handle("POST /api/v1/wallet/deposit", utils.RateLimitedHandler(utils.AuthHandler(handlers.HandleDepositWallet(walletRepository))))
	mux.Handle("POST /api/v1/wallet/withdraw", utils.RateLimitedHandler(utils.AuthHandler(handlers.HandleWithdrawWallet(walletRepository))))

	mux.Handle("GET /api/v1/exchange/rates", utils.RateLimitedHandler(utils.AuthHandler(handlers.HandleExchangeRates(exchangeRepository))))
	mux.Handle("POST /api/v1/exchange", utils.RateLimitedHandler(utils.AuthHandler(handlers.HandleExchangeWallet(walletRepository))))
}

func solve(mux *http.ServeMux) {
	_ = reflect.ValueOf(mux).Elem().FieldByName("tree").FieldByName("emptyChild")
	//todo если будет время
}
