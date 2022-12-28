package collections

import (
	"context"

	"github.com/deadloct/immutablex-cli/lib"
	"github.com/immutable/imx-core-sdk-golang/imx/api"
	log "github.com/sirupsen/logrus"
)

type AlchemyClientConfig struct {
	alchemyKey string
}

type AlchemyClient struct {
	client    lib.ClientWrapper
	shortcuts Shortcuts
}

func NewAlchemyClient(cfg AlchemyClientConfig) *AlchemyClient {
	return &AlchemyClient{
		client:    lib.NewClient(cfg.alchemyKey),
		shortcuts: NewShortcuts(),
	}
}

func (c *AlchemyClient) Start() error {
	return c.client.Start()
}

func (c *AlchemyClient) Stop() {
	c.client.Stop()
}

func (c *AlchemyClient) GetCollection(ctx context.Context, collection string) (*api.Collection, error) {
	if v := c.shortcuts.GetShortcutByName(collection); v != nil {
		collection = v.Addr
	}

	log.Debugf("fetching collection %s", collection)
	return c.client.GetClient().GetCollection(ctx, collection)
}

func (c *AlchemyClient) ListCollections(ctx context.Context, cfg *ListCollectionsConfig) ([]api.Collection, error) {
	req := c.getAPIListCollectionsRequest(ctx, cfg)

	resp, err := c.client.GetClient().ListCollections(req)
	if err != nil {
		return nil, err
	}

	if len(resp.Result) == 0 {
		return cfg.Collections, nil
	}

	cfg.Collections = append(cfg.Collections, resp.Result...)
	cfg.Cursor = resp.Cursor

	first := *resp.Result[0].UpdatedAt.Get()
	last := *resp.Result[len(resp.Result)-1].UpdatedAt.Get()
	log.Debugf("fetched %v collections from %v to %v", len(resp.Result), first, last)

	if resp.Remaining > 0 {
		return c.ListCollections(ctx, cfg)
	}

	return cfg.Collections, nil
}

func (c *AlchemyClient) getAPIListCollectionsRequest(ctx context.Context, cfg *ListCollectionsConfig) *api.ApiListCollectionsRequest {
	req := c.client.GetClient().NewListCollectionsRequest(ctx)

	if cfg.Blacklist != "" {
		req = req.Blacklist(cfg.Blacklist)
	}

	if cfg.Cursor != "" {
		req = req.Cursor((cfg.Cursor))
	}

	if cfg.Direction != "" {
		req = req.Direction(cfg.Direction)
	}

	if cfg.Keyword != "" {
		req = req.Keyword(cfg.Keyword)
	}

	if cfg.OrderBy != "" {
		req = req.OrderBy(cfg.OrderBy)
	}

	if cfg.Whitelist != "" {
		req = req.Whitelist(cfg.Whitelist)
	}

	return &req
}
