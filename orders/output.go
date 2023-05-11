package orders

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/deadloct/immutablex-go-lib/coinbase"
	"github.com/deadloct/immutablex-go-lib/utils"
	"github.com/immutable/imx-core-sdk-golang/imx/api"
	log "github.com/sirupsen/logrus"
)

func getPrice(order api.Order) float64 {
	amount, err := strconv.Atoi(order.GetBuy().Data.QuantityWithFees)
	if err != nil {
		return 0
	}

	decimals := int(*order.GetBuy().Data.Decimals)
	return float64(amount) * math.Pow10(-1*decimals)
}

func PrintOrderJSON(order api.Order) {
	data, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		log.Debugf("could not convert asset to json: %v\nasset: %#v\n", err, order)
		return
	}

	fmt.Println(string(data))
}

func PrintOrderNormal(order api.Order) {
	url := strings.Join([]string{utils.ImmutascanURL, "order", fmt.Sprint(order.OrderId)}, "/")
	price := getPrice(order)
	symbol := coinbase.CryptoSymbol(order.GetBuy().Type)
	fiatPrice := price * coinbase.GetCoinbaseClientInstance().RetrieveSpotPrice(symbol, coinbase.FiatUSD)
	fmt.Printf(`Order:
- Status: %s
- Price With Fees: %f %s / %.2f %s
- User: %s
- Date: %s
- Immutascan: %s%s`, order.Status, price, symbol, fiatPrice, coinbase.FiatUSD, order.User, order.GetUpdatedTimestamp(), url, "\n\n")
}

func PrintOrders(orders []api.Order, output string) {
	for _, o := range orders {
		switch strings.ToLower(output) {
		case "json":
			PrintOrderJSON(o)
		default:
			PrintOrderNormal(o)
		}
	}
}
