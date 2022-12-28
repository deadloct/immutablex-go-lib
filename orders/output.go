package orders

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/deadloct/immutablex-cli/lib"
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

func PrintOrderDetail(order api.Order) {
	data, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		log.Debugf("could not convert asset to json: %v\nasset: %#v\n", err, order)
		return
	}

	fmt.Println(string(data))
}

func PrintOrderSummary(order api.Order) {
	url := strings.Join([]string{lib.ImmutascanURL, "order", fmt.Sprint(order.OrderId)}, "/")
	ethPrice := getPrice(order)
	fiatPrice := ethPrice * lib.GetCoinbaseClientInstance().LastSpotPrice
	fmt.Printf(`Order:
- Status: %s
- Price With Fees: %f ETH / %.2f USD
- User: %s
- Date: %s
- Immutascan: %s%s`, order.Status, ethPrice, fiatPrice, order.User, order.GetUpdatedTimestamp(), url, "\n\n")
}

func PrintOrders(orders []api.Order, verbose bool) {
	for _, o := range orders {
		if verbose {
			PrintOrderDetail(o)
		} else {
			PrintOrderSummary(o)
		}
	}
}
