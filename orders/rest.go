package orders

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/immutable/imx-core-sdk-golang/imx/api"
	log "github.com/sirupsen/logrus"
)

const (
	ListOrdersEndpoint = "/v1/orders"
)

type RESTClientConfig struct {
	URL string
}

type RESTClient struct {
	client http.Client
	url    string
}

func NewRESTClient(cfg RESTClientConfig) *RESTClient {
	return &RESTClient{url: cfg.URL}
}

func (c *RESTClient) Start() error { return nil }

func (c *RESTClient) Stop() {}

func (c *RESTClient) GetOrder() (api.Order, error) { return api.Order{}, nil }

func (c *RESTClient) ListOrders(ctx context.Context, cfg *ListOrdersConfig) ([]api.Order, error) {
	url := c.getListOrdersURL(cfg)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var parsed api.ListOrdersResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		log.Errorf("could not parse response from server: %#v", err)
		return nil, err
	}

	if len(parsed.Result) == 0 {
		return cfg.Orders, nil
	}

	cfg.Orders = append(cfg.Orders, parsed.Result...)
	cfg.Cursor = parsed.Cursor

	first := *parsed.Result[0].UpdatedTimestamp.Get()
	last := *parsed.Result[len(parsed.Result)-1].UpdatedTimestamp.Get()
	log.Debugf("fetched %v orders from %v to %v", len(parsed.Result), first, last)

	if parsed.Remaining > 0 {
		return c.ListOrders(ctx, cfg)
	}

	return cfg.Orders, nil
}

func (c *RESTClient) getListOrdersURL(cfg *ListOrdersConfig) string {
	v := url.Values{}

	if cfg.AuxiliaryFeePercentages != "" {
		v.Set("auxiliary_fee_percentages", cfg.AuxiliaryFeePercentages)
	}

	if cfg.AuxiliaryFeeRecipients != "" {
		v.Set("auxiliary_fee_recipients", cfg.AuxiliaryFeeRecipients)
	}

	if cfg.BuyAssetID != "" {
		v.Set("buy_asset_id", cfg.BuyAssetID)
	}

	if cfg.BuyMaxQuantity != "" {
		v.Set("buy_max_quantity", cfg.BuyMaxQuantity)
	}

	if cfg.BuyMetadata != "" {
		v.Set("buy_metadata", cfg.BuyMetadata)
	}

	if cfg.BuyMinQuantity != "" {
		v.Set("buy_min_quantity", cfg.BuyMinQuantity)
	}

	if cfg.BuyTokenAddress != "" {
		v.Set("buy_token_address", cfg.BuyTokenAddress)
	}

	if cfg.BuyTokenID != "" {
		v.Set("buy_token_id", cfg.BuyTokenID)
	}

	if cfg.BuyTokenName != "" {
		v.Set("buy_token_name", cfg.BuyTokenName)
	}

	if cfg.BuyTokenType != "" {
		v.Set("buy_token_type", cfg.BuyTokenType)
	}

	if cfg.Direction != "" {
		v.Set("direction", cfg.Direction)
	}

	if cfg.IncludeFees {
		v.Set("include_fees", "true")
	}

	if cfg.MaxTimestamp != "" {
		v.Set("max_timestamp", cfg.MaxTimestamp)
	}

	if cfg.MinTimestamp != "" {
		v.Set("min_timestamp", cfg.MinTimestamp)
	}

	if cfg.OrderBy != "" {
		v.Set("order_by", cfg.OrderBy)
	}

	if cfg.SellAssetID != "" {
		v.Set("sell_asset_id", cfg.SellAssetID)
	}

	if cfg.SellMaxQuantity != "" {
		v.Set("sell_max_quantity", cfg.SellMaxQuantity)
	}

	if cfg.SellMetadata != "" {
		v.Set("sell_metadata", cfg.SellMetadata)
	}

	if cfg.SellMinQuantity != "" {
		v.Set("sell_min_quantity", cfg.SellMinQuantity)
	}

	if cfg.SellTokenAddress != "" {
		v.Set("sell_token_address", cfg.SellTokenAddress)
	}

	if cfg.SellTokenID != "" {
		v.Set("sell_token_id", cfg.SellTokenID)
	}

	if cfg.SellTokenName != "" {
		v.Set("sell_token_name", cfg.SellTokenName)
	}

	if cfg.SellTokenType != "" {
		v.Set("sell_token_type", cfg.SellTokenType)
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

	return c.url + ListOrdersEndpoint + "?" + v.Encode()
}
