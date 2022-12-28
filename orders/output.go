package orders

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/deadloct/immutablex-cli/lib"
	"github.com/immutable/imx-core-sdk-golang/imx/api"
	log "github.com/sirupsen/logrus"
)

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
	fmt.Printf(`Order:
- Status: %s
- User: %s
- Date: %s
- Immutascan: %s%s`, order.Status, order.User, order.GetUpdatedTimestamp(), url, "\n\n")
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
