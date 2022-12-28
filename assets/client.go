package assets

import (
	"context"

	"github.com/deadloct/immutablex-cli/lib"
	"github.com/immutable/imx-core-sdk-golang/imx/api"
	log "github.com/sirupsen/logrus"
)

type ListAssetsConfig struct {
	BuyOrders           bool
	Collection          string
	Direction           string
	IncludeFees         bool
	Metadata            string
	Name                string
	OrderBy             string
	SellOrders          bool
	Status              string
	UpdatedMaxTimestamp string
	UpdatedMinTimestamp string
	User                string

	// Used internally for recursion
	Assets []api.AssetWithOrders
	Cursor string
}

type Client interface {
	Start() error
	Stop()
	GetAsset(ctx context.Context, tokenAddress, tokenID string, includeFees bool) (*api.Asset, error)
	ListAssets(ctx context.Context, cfg ListAssetsConfig) ([]api.AssetWithOrders, error)
}

func NewClientConfig(alchemyKey string) interface{} {
	if alchemyKey == "" {
		return RESTClientConfig{url: lib.DefaultImmutableAPIURL}
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
