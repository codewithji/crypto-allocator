package services

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetExchangeRates_HappyPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"currency": "USD",
				"rates": {
					"BTC": "0.0000351465991008",
					"ETH": "0.0006380338349343"
				}	
			}
		}`))
	}))
	defer server.Close()

	client := &http.Client{}

	fetcher := CryptoExchangeRatesFetcher{}
	exchangeRates, err := fetcher.GetExchangeRates(client, server.URL)
	if err != nil {
		t.Fatalf("Expected no error but got %q", err)
	}

	if exchangeRates.Data.Currency != "USD" {
		t.Errorf("Expected currency to be USD but got %q", exchangeRates.Data.Currency)
	}

	expectedBTCRate := "0.0000351465991008"
	if exchangeRates.Data.Rates.BTC != expectedBTCRate {
		t.Errorf("Expected BTC rate %q but got %q", expectedBTCRate, exchangeRates.Data.Rates.BTC)
	}

	expectedETHRate := "0.0006380338349343"
	if exchangeRates.Data.Rates.ETH != expectedETHRate {
		t.Errorf("Expected ETH rate %q but got %q", expectedETHRate, exchangeRates.Data.Rates.ETH)
	}
}

func TestGetExchangeRates_StatusCodes(t *testing.T) {
	statusCodes := []int{
		http.StatusNotFound,
		http.StatusInternalServerError,
	}

	for _, statusCode := range statusCodes {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(statusCode)
		}))
		defer server.Close()

		client := &http.Client{}
		fetcher := CryptoExchangeRatesFetcher{}

		exchangeRates, err := fetcher.GetExchangeRates(client, server.URL)
		if err == nil {
			t.Fatalf("Expected error for status code %q but got none", statusCode)
		}

		expectedErr := fmt.Errorf("%v error occurred while attempting to fetch crypto exchange rates.", statusCode)

		if err.Error() != expectedErr.Error() {
			t.Errorf(`Expected %q but got %q`, expectedErr, err)
		}
		if exchangeRates.Data.Currency != "" {
			t.Errorf(`Expected empty string for Currency property but got %q`, exchangeRates.Data.Currency)
		}
		if exchangeRates.Data.Rates.BTC != "" {
			t.Errorf(`Expected empty string for BTC rate property but got %q`, exchangeRates.Data.Rates.BTC)
		}
		if exchangeRates.Data.Rates.ETH != "" {
			t.Errorf(`Expected empty string for ETH rate property but got %q`, exchangeRates.Data.Rates.ETH)
		}
	}
}

func TestGetExchangeRates_JsonDecodeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	client := &http.Client{}
	fetcher := CryptoExchangeRatesFetcher{}

	_, err := fetcher.GetExchangeRates(client, server.URL)
	expectedErr := "Failed to decode JSON: invalid character 'i' looking for beginning of value"

	if err == nil || err.Error() != expectedErr {
		t.Errorf("Expected error %q but got %q", expectedErr, err)
	}
}
