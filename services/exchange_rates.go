package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	btcRatio = 0.7
	ethRatio = 0.3
)

type Rates struct {
	BTC string `json:"BTC"`
	ETH string `json:"ETH"`
}

type CryptoExchangeRates struct {
	Data struct {
		Currency string `json:"currency"`
		Rates    Rates  `json:"rates"`
	} `json:"data"`
}

type ExchangeRatesFetcher interface {
	GetExchangeRates(client *http.Client, apiEndpoint string) (CryptoExchangeRates, error)
}

type CryptoExchangeRatesFetcher struct{}

func (c CryptoExchangeRatesFetcher) GetExchangeRates(client *http.Client, apiEndpoint string) (CryptoExchangeRates, error) {
	res, err := client.Get(apiEndpoint)
	if err != nil {
		return CryptoExchangeRates{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return CryptoExchangeRates{}, fmt.Errorf("%v error occurred while attempting to fetch crypto exchange rates.", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return CryptoExchangeRates{}, err
	}

	var exchangeRates CryptoExchangeRates
	err = json.Unmarshal(body, &exchangeRates)
	return exchangeRates, err
}
