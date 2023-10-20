package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/codewithji/crypto-allocator/services"
	"github.com/joho/godotenv"
)

func getBanner(fileName string) (string, error) {
	banner, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return string(banner), nil
}

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func main() {
	banner, err := getBanner("banner.txt")
	if err == nil && len(banner) > 0 {
		fmt.Printf("%v\n\n", banner)
	}

	reader := services.StdinReader{}
	investment, err := services.GetUserInvestmentInput(&reader, 3)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	fetcher := services.CryptoExchangeRatesFetcher{}

	cryptoAllocations, err := services.GetCryptoAllocations(fetcher, client, investment)
	if err != nil {
		fmt.Println("Failed to calculate crypto allocations: %w", err)
		return
	}

	jsonData, err := json.Marshal(cryptoAllocations)
	if err != nil {
		fmt.Println("Failed to serialize crypto allocations: %w", err)
		return
	}

	fmt.Println(string(jsonData))
}
