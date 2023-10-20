package services

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

type MockExchangeRatesFetcher struct{}

func (m *MockExchangeRatesFetcher) GetExchangeRates(client *http.Client, apiEndpoint string) (CryptoExchangeRates, error) {
	return CryptoExchangeRates{
		Data: struct {
			Currency string `json:"currency"`
			Rates    Rates  `json:"rates"`
		}{
			Currency: "USD",
			Rates: Rates{
				BTC: "0.000025",
				ETH: "0.00075",
			},
		},
	}, nil
}

func TestGetCryptoAllocations_HappyPath(t *testing.T) {
	mockFetcher := &MockExchangeRatesFetcher{}
	client := &http.Client{}
	investmentAmount := 10000.00

	res, err := GetCryptoAllocations(mockFetcher, client, investmentAmount)
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	parsedBtcAllocation, err := strconv.ParseFloat(res.BTC, 64)
	if err != nil {
		t.Errorf("Expected no error while parsing BTC allocation but got %v", err)
	}
	expectedBtcAllocation := 0.7 * 0.000025 * investmentAmount
	if parsedBtcAllocation != expectedBtcAllocation {
		t.Errorf("Expected BTC allocation of %v but got %v", expectedBtcAllocation, parsedBtcAllocation)
	}

	parsedEthAllocation, err := strconv.ParseFloat(res.ETH, 64)
	if err != nil {
		t.Errorf("Expected no error while parsing ETH allocation but got %v", err)
	}

	expectedEthAllocation := 0.3 * 0.00075 * investmentAmount
	if parsedEthAllocation != expectedEthAllocation {
		t.Errorf("Expected ETH allocation of %v but got %v", expectedEthAllocation, parsedEthAllocation)
	}
}

type MockExchangeRatesFetcherWithError struct{}

func (m *MockExchangeRatesFetcherWithError) GetExchangeRates(client *http.Client, apiEndpoint string) (CryptoExchangeRates, error) {
	return CryptoExchangeRates{}, fmt.Errorf("Error")
}

func TestGetCryptoAllocations_ExchangeRateError(t *testing.T) {
	mockFetcher := &MockExchangeRatesFetcherWithError{}
	client := &http.Client{}
	investmentAmount := 10000.00

	allocations, err := GetCryptoAllocations(mockFetcher, client, investmentAmount)
	if err == nil {
		t.Fatalf("Expected an error but got none")
	}
	if allocations.BTC != "" {
		t.Errorf(`Expected empty string for BTC allocation but got "%v"`, allocations.BTC)
	}
	if allocations.ETH != "" {
		t.Errorf(`Expected empty string for ETH allocation but got "%v"`, allocations.ETH)
	}
}

type MockExchangeRatesFetcherWithInvalidData struct{}

func (m *MockExchangeRatesFetcherWithInvalidData) GetExchangeRates(client *http.Client, apiEndpoint string) (CryptoExchangeRates, error) {
	return CryptoExchangeRates{
		Data: struct {
			Currency string `json:"currency"`
			Rates    Rates  `json:"rates"`
		}{
			Currency: "USD",
			Rates: Rates{
				BTC: "Invalid data",
				ETH: "0.1234",
			},
		},
	}, nil
}

func TestGetCryptoAllocations_InvalidData(t *testing.T) {
	mockFetcher := &MockExchangeRatesFetcherWithInvalidData{}
	client := &http.Client{}
	investmentAmount := 10000.00

	expectedError := `strconv.ParseFloat: parsing "Invalid data": invalid syntax`
	_, err := GetCryptoAllocations(mockFetcher, client, investmentAmount)
	if err.Error() != expectedError {
		t.Errorf(`Expected "%v" but got "%v"`, expectedError, err.Error())
	}
}
