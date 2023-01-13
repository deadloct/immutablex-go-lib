package coinbase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Currency string

const (
	BaseURL  = "https://api.coinbase.com"
	CacheFor = 30 * time.Second

	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyGBP Currency = "GBP"
)

var (
	coinbaseClientInstance *CoinbaseClient
	muCoinbase             sync.Mutex
)

type SupportedCurrencies string

type CoinbaseSpotPriceReponse struct {
	Data CoinbaseSpotPrice `json:"data"`
}

type CoinbaseSpotPrice struct {
	Base     string `json:"base"`
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type Price struct {
	Price         float64
	LastRetrieved time.Time
}

type CoinbaseClient struct {
	client         *http.Client
	lastSpotPrices map[Currency]Price
	lastRetrieved  time.Time
}

func GetCoinbaseClientInstance() *CoinbaseClient {
	muCoinbase.Lock()
	defer muCoinbase.Unlock()

	if coinbaseClientInstance == nil {
		coinbaseClientInstance = &CoinbaseClient{
			client:         &http.Client{},
			lastSpotPrices: make(map[Currency]Price),
		}
	}

	return coinbaseClientInstance
}

func (c *CoinbaseClient) RetrieveSpotPrice(currency Currency) float64 {
	if currency == "" {
		currency = CurrencyUSD
	}

	last, ok := c.lastSpotPrices[currency]
	if ok && time.Since(last.LastRetrieved) <= CacheFor {
		return c.lastSpotPrices[currency].Price
	}

	resp, err := c.client.Get(fmt.Sprintf("%s/v2/prices/ETH-%s/spot", BaseURL, currency))
	if err != nil {
		log.Errorf("all ETH-%s prices will be zero b/c error retrieving spot price: %v", currency, err)
		return 0
	}
	defer resp.Body.Close()

	var result CoinbaseSpotPriceReponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Errorf("all ETH-%s prices will be zero b/c error parsing spot price: %v", currency, err)
		return 0
	}

	amount, err := strconv.ParseFloat(result.Data.Amount, 64)
	if err != nil {
		log.Errorf("all ETH-USD prices will be zero b/c error parsing spot price: %v", err)
		return 0
	}

	c.lastSpotPrices[currency] = Price{Price: amount, LastRetrieved: time.Now()}
	return amount
}
