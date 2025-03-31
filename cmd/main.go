package main

import (
	"flag"
	"github.com/joho/godotenv"
	"gw-currency-wallet/internal/config"
	"gw-currency-wallet/internal/routes"
	"gw-currency-wallet/pkg/logger"
	"net/http"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "", "Config file location")
	flag.Parse()

	if configFile != "" {
		logger.Log.Debug("Loading env from file", "file", configFile)
		err := godotenv.Load(configFile)
		if err != nil {
			logger.Log.Error("Error loading .env file", "error", err)
			panic(err)
		}
	}

	config.Cfg = config.NewConfig()

	mux := http.NewServeMux()
	routes.Init(mux)
	err := http.ListenAndServe(":5002", mux)
	if err != nil {
		panic(err)
	}
}
