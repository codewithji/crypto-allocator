package services

import (
	"errors"
	"fmt"
	"strconv"
)

type StringReader interface {
	ReadInput() (string, error)
}

type StdinReader struct{}

func (s *StdinReader) ReadInput() (string, error) {
	var input string
	_, err := fmt.Scanln(&input)
	return input, err
}

func GetUserInvestmentInput(reader StringReader, attempts int) (float64, error) {
	if attempts <= 0 {
		return 0, errors.New("Unauthorized")
	}

	fmt.Printf("Welcome to Crypto Allocator!\n\n")
	fmt.Printf("Tell us your investment amount in USD, and we'll let you know how much Bitcoin (BTC) and Ethereum (ETH) to buy. We allocate 70%% into BTC and 30%% into ETH based on current market rates from Coinbase.\n\n")

	var investmentAmount float64

	for attempts > 0 {
		attempts--

		fmt.Print("Enter amount: ")
		input, err := reader.ReadInput()
		fmt.Println()

		if err != nil {
			fmt.Println(err)
			continue
		}

		investmentAmount, err = processInput(input, attempts)
		if err != nil {
			if attempts == 0 {
				return 0, errors.New("Reached max attempts. Please try again later.")
			}
			fmt.Printf("%v\n\n", err)
			continue
		}

		// if input is successfully processed, break out of loop
		break
	}

	return investmentAmount, nil
}

func processInput(input string, attempts int) (float64, error) {
	investmentAmount, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid input. Please enter a number. (Attempts remaining: %v)", attempts)
	}

	if investmentAmount <= 0 {
		return 0, fmt.Errorf("Please enter a number greater than 0. (Attempts remaining: %v)", attempts)
	}

	return investmentAmount, nil
}
