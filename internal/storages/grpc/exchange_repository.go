package grpc

import (
	"context"
	"fmt"
	proto_exchange "github.com/flowXM/proto-exchange/exchange"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gw-currency-wallet/internal/config"
	"gw-currency-wallet/internal/storages"
	"gw-currency-wallet/pkg/logger"
	"time"
)

var (
	cache       = make(map[string]*proto_exchange.ExchangeRateResponse)
	lastUpdated = make(map[string]time.Time)
)

type exchangeRepository struct {
}

func NewExchangeRepository() storages.ExchangeRepository {
	return &exchangeRepository{}
}

func (e *exchangeRepository) GetExchangeRates() (*proto_exchange.ExchangeRatesResponse, error) {
	conn, err := grpc.NewClient(config.Cfg.ProtoExchangeUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Error("Error connection", "error", err)
		return nil, err
	}
	defer conn.Close()

	client := proto_exchange.NewExchangeServiceClient(conn)

	rates, err := client.GetExchangeRates(context.TODO(), &proto_exchange.Empty{})
	if err != nil {
		return nil, err
	}

	return rates, nil
}

func (e *exchangeRepository) GetExchangeRateForCurrency(from, to storages.Currency) (*proto_exchange.ExchangeRateResponse, error) {
	conn, err := grpc.NewClient(config.Cfg.ProtoExchangeUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Error("Error connection", "error", err)
		return nil, err
	}
	defer conn.Close()
	key := fmt.Sprintf("%s%s", string(from), string(to))

	t, ok := lastUpdated[key]
	if ok && time.Now().Sub(t) < time.Minute {
		return cache[key], nil
	}

	client := proto_exchange.NewExchangeServiceClient(conn)

	rates, err := client.GetExchangeRateForCurrency(context.TODO(), &proto_exchange.CurrencyRequest{FromCurrency: string(from), ToCurrency: string(to)})
	if err != nil {
		return nil, err
	}

	lastUpdated[key] = time.Now()
	cache[key] = rates

	return rates, nil
}
