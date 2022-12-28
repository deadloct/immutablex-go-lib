package orders

import (
	"context"
	"log"

	"github.com/deadloct/immutablex-cli/lib"
	"github.com/immutable/imx-core-sdk-golang/imx/api"
)

type ListOrdersConfig struct {
	AuxiliaryFeePercentages string
	AuxiliaryFeeRecipients  string
	BuyAssetID              string
	BuyMaxQuantity          string
	BuyMetadata             string
	BuyMinQuantity          string
	BuyTokenAddress         string
	BuyTokenID              string
	BuyTokenName            string
	BuyTokenType            string
	Cursor                  string
	Direction               string
	IncludeFees             bool
	MaxTimestamp            string
	MinTimestamp            string
	OrderBy                 string
	Orders                  []api.Order
	PageSize                int
	SellAssetID             string
	SellMaxQuantity         string
	SellMetadata            string
	SellMinQuantity         string
	SellTokenAddress        string
	SellTokenID             string
	SellTokenName           string
	SellTokenType           string
	Status                  string
	UpdatedMaxTimestamp     string
	UpdatedMinTimestamp     string
	User                    string
}

type Client interface {
	Start() error
	Stop()
	GetOrder() (api.Order, error)
	ListOrders(ctx context.Context, cfg *ListOrdersConfig) ([]api.Order, error)
}

func NewClientConfig(alchemyKey string) interface{} {
	if alchemyKey == "" {
		return RESTClientConfig{URL: lib.DefaultImmutableAPIURL}
	}

	return AlchemyClientConfig{alchemyKey: alchemyKey}
}

func NewClient(cfg interface{}) Client {
	switch v := cfg.(type) {
	case RESTClientConfig:
		return NewRESTClient(v)
	case AlchemyClientConfig:
		return NewAlchemyClient(v)
	default:
		log.Panicf("invalid client config")
	}

	return nil
}
