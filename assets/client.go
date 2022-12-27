package assets

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/deadloct/immutablex-cli/lib"
	"github.com/immutable/imx-core-sdk-golang/imx/api"
	log "github.com/sirupsen/logrus"
)

type ListAssetsConfig struct {
	BuyOrders           bool
	Collection          string
	Direction           string
	IncludeFees         bool
	Metadata            []string
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
	Before string
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

func ParseMetadata(metadata []string) string {
	metamap := make(map[string][]string, len(metadata))
	for _, item := range metadata {
		parts := strings.SplitN(item, "=", 2)
		if len(parts) != 2 {
			log.Debugf("could not parse metadata item %s into a key=value pair", item)
			continue
		}

		metamap[parts[0]] = append(metamap[parts[0]], parts[1])
	}

	data, err := json.Marshal(metamap)
	if err != nil {
		log.Debugf("skipping metamata completely because it could not be converted to json: %v", err)
	}

	return string(data)
}
