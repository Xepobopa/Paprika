package utils

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		ENVIRONMENT string `env:"ENVIRONMENT" envDefault:"DEV"`
		Nats        NATS
		Mexc        MEXC
	}

	NATS struct {
		URL string `env:"NATS_URL" envDefault:"nats://localhost:4222"`
	}
	
	MEXC struct {
		// SPOT_WS       string `env:"MEXC_SPOT_WS,required"`
		// FUTURES_WS    string `env:"MEXC_FUTURES_WS,required"`
		// ACCESS_TOKEN  string `env:"MEXC_ACCESS_TOKEN,required"`
		// PRIVATE_TOKEN string `env:"MEXC_PRIVATE_TOKEN,required"`
		API           MEXC_API
	}
	
	MEXC_API struct {
		SPOT_TICKERS_24HR       string `env:"MEXC_API_TICKER_24HR,required"`
		CONTRACTS_TICKERS string `env:"MEXC_API_CONTRACTS_TICKERS,required"`
		// CONFIG_GETALL     string `env:"MEXC_API_CONFIG_GETALL,required"`
		// CONTRACTS_DETAIL  string `env:"MEXC_API_CONTRACTS_DETAIL,required"`
		// ORDER_BOOK_TICKER string `env:"MEXC_API_ORDER_BOOK_TICKER,required"`
		// EXCHANGE_INFO     string `env:"MEXC_API_EXCHANGE_INFO,required"`
	}
)

func NewConfig(filename string) (*Config, error) {
	cfg := &Config{}

	// Remove in prod
	err := godotenv.Load(filename)
	if err != nil {
		return nil, fmt.Errorf("loading .env error: %v", err)
	}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %v", err)
	}

	return cfg, nil
}
