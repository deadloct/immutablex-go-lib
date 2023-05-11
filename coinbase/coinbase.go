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

type CryptoSymbol string

type FiatSymbol string

const (
	BaseURL  = "https://api.coinbase.com"
	CacheFor = 30 * time.Second

	CryptoETH  CryptoSymbol = "ETH"
	CryptoIMX  CryptoSymbol = "IMX"
	CryptoUSDC CryptoSymbol = "USDC"

	FiatUSD FiatSymbol = "USD"
	FiatEUR FiatSymbol = "EUR"
	FiatGBP FiatSymbol = "GBP"
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
	lastSpotPrices map[string]Price
}

func GetCoinbaseClientInstance() *CoinbaseClient {
	muCoinbase.Lock()
	defer muCoinbase.Unlock()

	if coinbaseClientInstance == nil {
		coinbaseClientInstance = &CoinbaseClient{
			client:         &http.Client{},
			lastSpotPrices: make(map[string]Price),
		}
	}

	return coinbaseClientInstance
}

func (c *CoinbaseClient) RetrieveSpotPrice(crypto CryptoSymbol, fiat FiatSymbol) float64 {
	if fiat == "" {
		fiat = FiatUSD
	}

	if crypto == "" {
		crypto = CryptoETH
	}

	spotKey := c.getSpotKey(fiat, crypto)

	last, ok := c.lastSpotPrices[spotKey]
	if ok && time.Since(last.LastRetrieved) <= CacheFor {
		return c.lastSpotPrices[spotKey].Price
	}

	resp, err := c.client.Get(fmt.Sprintf("%s/v2/prices/%s-%s/spot", BaseURL, crypto, fiat))
	if err != nil {
		log.Errorf("all %s-%s prices will be zero b/c error retrieving spot price: %v", crypto, fiat, err)
		return 0
	}
	defer resp.Body.Close()

	var result CoinbaseSpotPriceReponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Errorf("all %s-%s prices will be zero b/c error parsing spot price: %v", crypto, fiat, err)
		return 0
	}

	amount, err := strconv.ParseFloat(result.Data.Amount, 64)
	if err != nil {
		log.Errorf("all %s-%s prices will be zero b/c error parsing spot price: %v", crypto, fiat, err)
		return 0
	}

	c.lastSpotPrices[spotKey] = Price{Price: amount, LastRetrieved: time.Now()}
	return amount
}

func (c *CoinbaseClient) getSpotKey(f FiatSymbol, cr CryptoSymbol) string {
	return fmt.Sprintf("%s-%s", f, cr)
}
