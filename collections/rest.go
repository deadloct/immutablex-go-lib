package collections

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/immutable/imx-core-sdk-golang/imx/api"
	log "github.com/sirupsen/logrus"
)

const (
	GetCollectionEndpoint   = "/v1/collections"
	ListCollectionsEndpoint = "/v1/collections"
)

type RESTClientConfig struct {
	url string
}

type RESTClient struct {
	client    *http.Client
	url       string
	shortcuts Shortcuts
}

func NewRESTClient(cfg RESTClientConfig) *RESTClient {
	return &RESTClient{
		url:       cfg.url,
		client:    &http.Client{},
		shortcuts: NewShortcuts(),
	}
}

func (c *RESTClient) Start() error { return nil }

func (c *RESTClient) Stop() {}

func (c *RESTClient) GetCollection(ctx context.Context, collection string) (*api.Collection, error) {
	if v := c.shortcuts.GetShortcutByName(collection); v != nil {
		collection = v.Addr
	}

	log.Debugf("fetching collection %s", collection)
	url := c.url + GetCollectionEndpoint + "/" + collection
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result api.Collection
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Errorf("could not parse response from server: %#v", err)
		return nil, err
	}

	return &result, nil
}

func (c *RESTClient) ListCollections(ctx context.Context, cfg *ListCollectionsConfig) ([]api.Collection, error) {
	url := c.getListCollectionsURL(cfg)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var parsed api.ListCollectionsResponse

	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		log.Errorf("could not parse response from server: %#v", err)
		return nil, err
	}

	if len(parsed.Result) == 0 {
		return cfg.Collections, nil
	}

	cfg.Collections = append(cfg.Collections, parsed.Result...)
	cfg.Cursor = parsed.Cursor

	first := *parsed.Result[0].UpdatedAt.Get()
	last := *parsed.Result[len(parsed.Result)-1].UpdatedAt.Get()
	log.Debugf("fetched %v collections from %v to %v", len(parsed.Result), first, last)

	if parsed.Remaining > 0 {
		return c.ListCollections(ctx, cfg)
	}

	return cfg.Collections, nil
}

func (c *RESTClient) getListCollectionsURL(cfg *ListCollectionsConfig) string {
	v := url.Values{}

	if cfg.Blacklist != "" {
		v.Set("blacklist", cfg.Blacklist)
	}

	if cfg.Cursor != "" {
		v.Set("cursor", cfg.Cursor)
	}

	if cfg.Direction != "" {
		v.Set("direction", cfg.Direction)
	}

	if cfg.Keyword != "" {
		v.Set("keyword", cfg.Keyword)
	}

	if cfg.OrderBy != "" {
		v.Set("order_by", cfg.OrderBy)
	}

	if cfg.Whitelist != "" {
		v.Set("whitelist", cfg.Whitelist)
	}

	return c.url + ListCollectionsEndpoint + "?" + v.Encode()
}
