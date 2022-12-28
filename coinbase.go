package lib

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
)

const BaseURL = "https://api.coinbase.com"

var (
	coinbaseClientInstance *CoinbaseClient
	muCoinbase             sync.Mutex
)

type CoinbaseSpotPriceReponse struct {
	Data CoinbaseSpotPrice `json:"data"`
}

type CoinbaseSpotPrice struct {
	Base     string `json:"base"`
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

func init() {
	GetCoinbaseClientInstance().RetrieveSpotPrice()
}

type CoinbaseClient struct {
	client        *http.Client
	LastSpotPrice float64
}

func GetCoinbaseClientInstance() *CoinbaseClient {
	muCoinbase.Lock()
	defer muCoinbase.Unlock()

	if coinbaseClientInstance == nil {
		coinbaseClientInstance = &CoinbaseClient{client: &http.Client{}}
	}

	return coinbaseClientInstance
}

func (c *CoinbaseClient) RetrieveSpotPrice() float64 {
	resp, err := c.client.Get(BaseURL + "/v2/prices/ETH-USD/spot")
	if err != nil {
		log.Errorf("all ETH-USD prices will be zero b/c error retrieving spot price: %v", err)
		return 0
	}
	defer resp.Body.Close()

	var result CoinbaseSpotPriceReponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Errorf("all ETH-USD prices will be zero b/c error parsing spot price: %v", err)
		return 0
	}

	c.LastSpotPrice, err = strconv.ParseFloat(result.Data.Amount, 64)
	if err != nil {
		log.Errorf("all ETH-USD prices will be zero b/c error parsing spot price: %v", err)
		return 0
	}

	return c.LastSpotPrice
}
