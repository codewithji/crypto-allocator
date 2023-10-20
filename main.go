package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/codewithji/crypto-allocator/services"
)

func getBanner(fileName string) (string, error) {
	banner, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return string(banner), nil
}

func main() {
	banner, err := getBanner("banner.txt")
	if err == nil {
		fmt.Printf("%v\n\n", banner)
	}

	reader := services.StdinReader{}
	investment, err := services.GetUserInvestmentInput(&reader, 3)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	fetcher := services.CryptoExchangeRatesFetcher{}

	cryptoAllocations, err := services.GetCryptoAllocations(fetcher, client, investment)
	if err != nil {
		fmt.Println("Error occurred while calculating crypto allocations")
	}

	jsonData, err := json.Marshal(cryptoAllocations)
	if err != nil {
		fmt.Println("Error occurred while serializing crypto allocations data")
	}

	fmt.Println(string(jsonData))
}
