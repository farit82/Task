package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	updateInterval = 10 * time.Minute
)

type CoinData struct {
	ID     string  `json:"id"`
	Symbol string  `json:"symbol"`
	Name   string  `json:"name"`
	Price  float64 `json:"current_price"`
}

var (
	coinDataCache []CoinData
	cacheTime     time.Time
)

func getData() ([]CoinData, error) {

	url := "http://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1"
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Parse JSON response
	var coinData []CoinData
	err = json.Unmarshal(body, &coinData)
	if err != nil {
		return nil, err
	}
	coinDataCache = coinData
	cacheTime = time.Now()

	return coinData, nil
}

func getPrice(coinSymbol string) (float64, error) {
	coinData, err := getData()
	if err != nil {
		return 0, err
	}
	for _, coin := range coinData {
		if coin.Symbol == coinSymbol {
			return coin.Price, nil
		}
	}

	return 0, fmt.Errorf("Coin not found: %s", coinSymbol)
}

func main() {
	for {
		coinSymbol := "btc"
		price, err := getPrice(coinSymbol)
		if err != nil {
			fmt.Printf("Failed to get price for %s: %v\n", coinSymbol, err)
		} else {
			fmt.Printf("Price of %s: $%.2f\n", coinSymbol, price)
		}
		time.Sleep(updateInterval)
	}
}
