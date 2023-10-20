package services

import (
	"fmt"
	"net/http"
	"strconv"
)

const apiEndpoint = "https://api.coinbase.com/v2/exchange-rates?currency=USD"

type CryptoAllocations struct {
	BTC string `json:"BTC"`
	ETH string `json:"ETH"`
}

func GetCryptoAllocations(fetcher ExchangeRatesFetcher, client *http.Client, investmentAmount float64) (CryptoAllocations, error) {
	exchangeRates, err := fetcher.GetExchangeRates(client, apiEndpoint)
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
