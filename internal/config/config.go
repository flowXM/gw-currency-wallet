package config

import (
	"gw-currency-wallet/pkg/logger"
	"os"
	"strconv"
)

var Cfg Config

type PostgresConfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     uint16
}

type Config struct {
	ProtoExchangeUrl string
	JWTSecretKey     string
	SaltSize         uint8

	Postgres PostgresConfig
}

func NewConfig() Config {
	postgres := PostgresConfig{
		User:     GetEnv("POSTGRES_USER", DefaultDBUser),
		Password: GetEnv("POSTGRES_PASSWORD", DefaultDBPassword),
		Name:     GetEnv("POSTGRES_DB", DefaultDBName),
		Host:     GetEnv("POSTGRES_SERVER", DefaultDBHost),
		Port:     GetEnvUint16("POSTGRES_PORT", DefaultDBPort),
	}

	config := Config{
		ProtoExchangeUrl: GetEnv("PROTO_EXCHANGE_URL", DefaultProtoExchangeUrl),
		JWTSecretKey:     GetEnv("JWT_SECRET_KEY", DefaultJWTSecretKey),
		SaltSize:         GetEnvUint8("SALT_SIZE", DefaultSaltSize),
		Postgres:         postgres,
	}

	logger.Log.Debug("Loaded config", "config", config)

	return config
}

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

func GetEnvUint8(key string, fallback uint8) uint8 {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	res, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return fallback
	}

	return uint8(res)
}
