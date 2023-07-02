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

	var coinData []CoinData
	err = json.Unmarshal(body, &coinData)
	if err != nil {
		return nil, err
	}
	coinDataCache = coinData
	cacheTime = time.Now()

	return coinData, nil
}

func getPrice(coinSymbol, ethereum string) (float64, float64, error) {
	coinData, err := getData()
	if err != nil {
		return 0, 0, err
	}

	var coinPrice, ethereumPrice float64

	for _, coin := range coinData {
		if coin.Symbol == coinSymbol {
			coinPrice = coin.Price
			fmt.Printf("Coin: %s, Price: %.2f USD\n", coin.Name, coin.Price)
		} else if coin.Symbol == ethereum {
			ethereumPrice = coin.Price
			fmt.Printf("Coin: %s, Price: %.2f USD\n", coin.Name, coin.Price)
		}
	}

	if coinPrice != 0 && ethereumPrice != 0 {
		return coinPrice, ethereumPrice, nil
	}

	return 0, 0, fmt.Errorf("Coins not found: %s and %s", coinSymbol, ethereum)
}

func main() {
	for {
		coinSymbol := "btc"
		ethereum := "eth"
		_, _, err := getPrice(coinSymbol, ethereum)
		if err != nil {
			fmt.Printf("Error in getting course: %v\n", err)
		}

		time.Sleep(updateInterval)
	}
}
