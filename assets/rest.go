package assets

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/deadloct/immutablex-cli/lib/collections"
	"github.com/immutable/imx-core-sdk-golang/imx/api"
	log "github.com/sirupsen/logrus"
)

const (
	GetAssetEndpoint   = "/v1/assets"
	ListAssetsEndpoint = "/v1/assets"
)

type RESTClientConfig struct {
	url string
}

type RESTClient struct {
	client    *http.Client
	url       string
	shortcuts collections.Shortcuts
}

func NewRESTClient(cfg RESTClientConfig) *RESTClient {
	return &RESTClient{
		client:    &http.Client{},
		url:       cfg.url,
		shortcuts: collections.NewShortcuts(),
	}
}

func (c *RESTClient) Start() error {
	return nil
}

func (c *RESTClient) Stop() {}

func (c *RESTClient) GetAsset(ctx context.Context, tokenAddress, tokenID string, includeFees bool) (*api.Asset, error) {
	if s := c.shortcuts.GetShortcutByName(tokenAddress); s != nil {
		tokenAddress = s.Addr
	}

	log.Debugf("fetching asset id %s from collection %s (with fees:%b)", tokenAddress, tokenID, includeFees)
	url := strings.Join([]string{c.url + GetAssetEndpoint, tokenAddress, tokenID}, "/")
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result api.Asset
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Errorf("could not parse response from server: %#v", err)
		return nil, err
	}

	return &result, nil
}

func (c *RESTClient) ListAssets(ctx context.Context, cfg ListAssetsConfig) ([]api.AssetWithOrders, error) {
	url := c.getListAssetsURL(cfg)
	getResp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer getResp.Body.Close()

	var resp api.ListAssetsResponse
	if err := json.NewDecoder(getResp.Body).Decode(&resp); err != nil {
		log.Errorf("could not parse response from server: %#v", err)
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
		return c.ListAssets(ctx, cfg)
	}

	// Attempt to fetch earlier assets
	if len(resp.Result) > 0 {
		cfg.Before = last
		return c.ListAssets(ctx, cfg)
	}

	return cfg.Assets, nil
}

func (c *RESTClient) getListAssetsURL(cfg ListAssetsConfig) string {
	v := url.Values{}

	if cfg.BuyOrders {
		v.Set("buy_orders", "true")
	}

	if cfg.Before != "" {
		v.Set("before", cfg.Before)
	}

	collectionAddr := cfg.Collection
	if s := c.shortcuts.GetShortcutByName(collectionAddr); s != nil {
		collectionAddr = s.Addr
	}

	if collectionAddr != "" {
		v.Set("collection", collectionAddr)
	}

	if cfg.Cursor != "" {
		v.Set("cursor", cfg.Cursor)
	}

	if cfg.Direction != "" {
		v.Set("direction", cfg.Direction)
	}

	if cfg.IncludeFees {
		v.Set("include_fees", "true")
	}

	if len(cfg.Metadata) > 0 {
		v.Set("metadata", ParseMetadata(cfg.Metadata))
	}

	if cfg.Name != "" {
		v.Set("name", cfg.Name)
	}

	if cfg.OrderBy != "" {
		v.Set("order_by", cfg.OrderBy)
	}

	if cfg.SellOrders {
		v.Set("sell_orders", "true")
	}

	if cfg.Status != "" {
		v.Set("status", cfg.Status)
	}

	if cfg.UpdatedMaxTimestamp != "" {
		v.Set("updated_max_timestamp", cfg.UpdatedMaxTimestamp)

	}

	if cfg.UpdatedMinTimestamp != "" {
		v.Set("updated_min_timestamp", cfg.UpdatedMinTimestamp)
	}

	if cfg.User != "" {
		v.Set("user", cfg.User)
	}

	return c.url + ListAssetsEndpoint + "?" + v.Encode()
}
