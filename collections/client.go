package collections

import (
	"context"

	"github.com/deadloct/immutablex-cli/lib"
	"github.com/immutable/imx-core-sdk-golang/imx/api"
	log "github.com/sirupsen/logrus"
)

type ListCollectionsConfig struct {
	Blacklist string
	Direction string
	Keyword   string
	OrderBy   string
	Whitelist string

	// Used internally for recursion
	Collections []api.Collection
	Cursor      string
}

type Client interface {
	Start() error
	Stop()
	GetCollection(ctx context.Context, collection string) (*api.Collection, error)
	ListCollections(ctx context.Context, cfg *ListCollectionsConfig) ([]api.Collection, error)
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
