package assets

import (
	"context"

	"github.com/deadloct/immutablex-go-lib/collections"
	"github.com/deadloct/immutablex-go-lib/imx"
	"github.com/immutable/imx-core-sdk-golang/imx/api"
	log "github.com/sirupsen/logrus"
)

type AlchemyClientConfig struct {
	alchemyKey string
}

type AlchemyClient struct {
	client    imx.ClientWrapper
	shortcuts collections.Shortcuts
}

func NewAlchemyClient(cfg AlchemyClientConfig) *AlchemyClient {
	return &AlchemyClient{
		client:    imx.NewClient(cfg.alchemyKey),
		shortcuts: collections.NewShortcuts(),
	}
}

func (am *AlchemyClient) Start() error {
	return am.client.Start()
}

func (am *AlchemyClient) Stop() {
	am.client.Stop()
}

func (am *AlchemyClient) GetAsset(ctx context.Context, tokenAddress, tokenID string, includeFees bool) (*api.Asset, error) {
	if s := am.shortcuts.GetShortcutByName(tokenAddress); s != nil {
		tokenAddress = s.Addr
	}

	log.Debugf("fetching asset id %s from collection %s (with fees:%t)", tokenAddress, tokenID, includeFees)
	return am.client.GetClient().GetAsset(ctx, tokenAddress, tokenID, &includeFees)
}

func (am *AlchemyClient) ListAssets(
	ctx context.Context,
	cfg ListAssetsConfig,
) ([]api.AssetWithOrders, error) {

	req := am.getAPIListAssetsRequest(ctx, cfg)
	resp, err := am.client.GetClient().ListAssets(&req)
	if err != nil {
		return nil, err
	}

	if len(resp.Result) == 0 {
		return cfg.Assets, nil
	}

	cfg.Assets = append(cfg.Assets, resp.Result...)
	cfg.Cursor = resp.Cursor

	first := *resp.Result[0].UpdatedAt.Get()
	last := *resp.Result[len(resp.Result)-1].UpdatedAt.Get()
	log.Debugf("fetched %v assets from %v to %v", len(resp.Result), first, last)

	if resp.Remaining > 0 {
		return am.ListAssets(ctx, cfg)
	}

	return cfg.Assets, nil
}

func (am *AlchemyClient) getAPIListAssetsRequest(ctx context.Context, cfg ListAssetsConfig) api.ApiListAssetsRequest {
	collectionAddr := cfg.Collection
	if s := am.shortcuts.GetShortcutByName(collectionAddr); s != nil {
		collectionAddr = s.Addr
	}

	req := am.client.GetClient().NewListAssetsRequest(ctx).Collection(collectionAddr)

	if cfg.BuyOrders {
		req = req.BuyOrders(cfg.BuyOrders)
	}

	if cfg.Cursor != "" {
		req = req.Cursor(cfg.Cursor)
	}

	if cfg.Direction != "" {
		req = req.Direction(cfg.Direction)
	}

	if cfg.IncludeFees {
		req = req.IncludeFees(cfg.IncludeFees)
	}

	if cfg.Metadata != "" {
		req = req.Metadata(cfg.Metadata)
	}

	if cfg.Name != "" {
		req = req.IncludeFees(cfg.IncludeFees)
	}

	if cfg.OrderBy != "" {
		req.OrderBy(cfg.OrderBy)
	} else {
		req.OrderBy("updated_at")
	}

	if cfg.SellOrders {
		req = req.SellOrders(cfg.SellOrders)
	}

	if cfg.Status != "" {
		req = req.Status(cfg.Status)
	}

	if cfg.UpdatedMaxTimestamp != "" {
		req = req.UpdatedMaxTimestamp(cfg.UpdatedMaxTimestamp)
	}

	if cfg.UpdatedMinTimestamp != "" {
		req = req.UpdatedMinTimestamp(cfg.UpdatedMinTimestamp)
	}

	if cfg.User != "" {
		req = req.User(cfg.User)
	}

	return req
}
