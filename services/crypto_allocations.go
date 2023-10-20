package services

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type CryptoAllocations struct {
	BTC string `json:"BTC"`
	ETH string `json:"ETH"`
}

func GetCryptoAllocations(fetcher ExchangeRatesFetcher, client *http.Client, investmentAmount float64) (CryptoAllocations, error) {
	// The Coinbase API endpoint is public and does not require authentication
	// Using an env variable is simply for demonstrative purposes

	// Add the following to your local .env file if testing locally:
	// API_URL=https://api.coinbase.com/v2/exchange-rates?currency=USD
	apiUrl := os.Getenv("API_URL")

	exchangeRates, err := fetcher.GetExchangeRates(client, apiUrl)
	if err != nil {
		return CryptoAllocations{}, err
	}

	btcRate, err := strconv.ParseFloat(exchangeRates.Data.Rates.BTC, 64)
	if err != nil {
		return CryptoAllocations{}, err
	}
	ethRate, err := strconv.ParseFloat(exchangeRates.Data.Rates.ETH, 64)
	if err != nil {
		return CryptoAllocations{}, err
	}

	btcAllocation := btcRatio * investmentAmount * btcRate
	ethAllocation := ethRatio * investmentAmount * ethRate

	allocations := CryptoAllocations{
		BTC: fmt.Sprintf("%f", btcAllocation),
		ETH: fmt.Sprintf("%f", ethAllocation),
	}

	return allocations, nil
}
