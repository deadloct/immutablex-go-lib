package coinbase

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	BaseURL  = "https://api.coinbase.com"
	CacheFor = 30 * time.Second
)

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

type CoinbaseClient struct {
	client        *http.Client
	lastSpotPrice float64
	lastRetrieved time.Time
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
	if time.Since(c.lastRetrieved) <= CacheFor {
		return c.lastSpotPrice
	}

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

	c.lastSpotPrice, err = strconv.ParseFloat(result.Data.Amount, 64)
	if err != nil {
		log.Errorf("all ETH-USD prices will be zero b/c error parsing spot price: %v", err)
		return 0
	}

	c.lastRetrieved = time.Now()
	return c.lastSpotPrice
}
