package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"gw-currency-wallet/internal/config"
	"gw-currency-wallet/pkg/logger"
	"time"
)

func NewClient() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", config.Cfg.DBHost, config.Cfg.DBPort, config.Cfg.DBUser, config.Cfg.DBName, config.Cfg.DBPassword)

	db, err := DoWithRetries(
		func() (*sql.DB, error) {
			return sql.Open("postgres", dsn)
		},
		5,
	)

	if err != nil {
		return nil, err
	}

	logger.Log.Info("Successfully connected to DB")

	return db, nil
}

func DoWithRetries(fn func() (*sql.DB, error), attempts int) (*sql.DB, error) {
	var db *sql.DB
	var err error
	for attempt := 1; attempt <= attempts; attempt++ {
		logger.Log.Info("Trying connect to DB", "attempt", attempt)
		db, err = fn()
		if err != nil {
			time.Sleep(time.Millisecond * 5)
			continue
		}

		return db, nil
	}

	return nil, err
}
