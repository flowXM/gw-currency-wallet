package grpc

import (
	"context"
	proto_exchange "github.com/flowXM/proto-exchange/exchange"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gw-currency-wallet/internal/config"
	"gw-currency-wallet/internal/storages"
	"log"
)

type exchangeRepository struct {
}

func NewExchangeRepository() storages.ExchangeRepository {
	return &exchangeRepository{}
}

func (e *exchangeRepository) GetExchangeRates() (*proto_exchange.ExchangeRatesResponse, error) {
	conn, err := grpc.NewClient(config.Cfg.ProtoExchangeUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
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
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer conn.Close()

	client := proto_exchange.NewExchangeServiceClient(conn)

	rates, err := client.GetExchangeRateForCurrency(context.TODO(), &proto_exchange.CurrencyRequest{FromCurrency: string(from), ToCurrency: string(to)})
	if err != nil {
		return nil, err
	}

	return rates, nil
}
