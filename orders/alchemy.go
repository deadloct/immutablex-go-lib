package orders

import (
	"context"
	"math"

	"github.com/deadloct/immutablex-cli/lib"
	"github.com/deadloct/immutablex-cli/lib/collections"
	"github.com/immutable/imx-core-sdk-golang/imx/api"
	log "github.com/sirupsen/logrus"
)

type AlchemyClientConfig struct {
	alchemyKey string
}

type AlchemyClient struct {
	client    lib.ClientWrapper
	shortcuts collections.Shortcuts
}

func NewAlchemyClient(cfg AlchemyClientConfig) *AlchemyClient {
	return &AlchemyClient{
		client:    lib.NewClient(cfg.alchemyKey),
		shortcuts: collections.NewShortcuts(),
	}
}

func (c *AlchemyClient) Start() error {
	return c.client.Start()
}

func (c *AlchemyClient) Stop() {
	c.client.Stop()
}

func (c *AlchemyClient) GetOrder() (api.Order, error) {
	return api.Order{}, nil
}

func (c *AlchemyClient) ListOrders(ctx context.Context, cfg *ListOrdersConfig) ([]api.Order, error) {
	req := c.getAPIListOrdersRequest(ctx, cfg)
	resp, err := c.client.GetClient().ListOrders(req)
	if err != nil {
		return nil, err
	}

	if len(resp.Result) == 0 {
		return cfg.Orders, nil
	}

	cfg.Orders = append(cfg.Orders, resp.Result...)
	cfg.Orders = cfg.Orders[:int(math.Min(float64(len(cfg.Orders)), float64(cfg.PageSize)))]
	cfg.Cursor = resp.Cursor

	first := *resp.Result[0].UpdatedTimestamp.Get()
	last := *resp.Result[len(resp.Result)-1].UpdatedTimestamp.Get()
	log.Debugf("fetched %v collections from %v to %v", len(resp.Result), first, last)

	getMore := len(cfg.Orders) < cfg.PageSize || cfg.PageSize == 0
	if resp.Remaining > 0 && getMore {
		return c.ListOrders(ctx, cfg)
	}

	return cfg.Orders, nil
}

func (c *AlchemyClient) getAPIListOrdersRequest(ctx context.Context, cfg *ListOrdersConfig) *api.ApiListOrdersRequest {
	req := c.client.GetClient().NewListOrdersRequest(ctx)

	if cfg.AuxiliaryFeePercentages != "" {
		req = req.AuxiliaryFeePercentages(cfg.AuxiliaryFeePercentages)
	}

	if cfg.AuxiliaryFeeRecipients != "" {
		req = req.AuxiliaryFeeRecipients(cfg.AuxiliaryFeeRecipients)
	}

	if cfg.BuyAssetID != "" {
		req = req.BuyAssetId(cfg.BuyAssetID)
	}

	if cfg.BuyMaxQuantity != "" {
		req = req.BuyMaxQuantity(cfg.BuyMaxQuantity)
	}

	if cfg.BuyMetadata != "" {
		req = req.BuyMetadata(cfg.BuyMetadata)
	}

	if cfg.BuyMinQuantity != "" {
		req = req.BuyMinQuantity(cfg.BuyMinQuantity)
	}

	if cfg.BuyTokenAddress != "" {
		req = req.BuyTokenAddress(cfg.BuyTokenAddress)
	}

	if cfg.BuyTokenID != "" {
		req = req.BuyTokenId(cfg.BuyTokenID)
	}

	if cfg.BuyTokenName != "" {
		req = req.BuyTokenName(cfg.BuyTokenName)
	}

	if cfg.BuyTokenType != "" {
		req = req.BuyTokenType(cfg.BuyTokenType)
	}

	if cfg.Direction != "" {
		req = req.Direction(cfg.Direction)
	}

	if cfg.IncludeFees {
		req = req.IncludeFees(cfg.IncludeFees)
	}

	if cfg.MaxTimestamp != "" {
		req = req.MaxTimestamp(cfg.MaxTimestamp)
	}

	if cfg.MinTimestamp != "" {
		req = req.MinTimestamp(cfg.MinTimestamp)
	}

	if cfg.OrderBy != "" {
		req = req.OrderBy(cfg.OrderBy)
	}

	if cfg.PageSize > 0 {
		req = req.PageSize(int32(cfg.PageSize))
	}

	if cfg.SellAssetID != "" {
		req = req.SellAssetId(cfg.SellAssetID)
	}

	if cfg.SellMaxQuantity != "" {
		req = req.SellMaxQuantity(cfg.SellMaxQuantity)
	}

	if cfg.SellMetadata != "" {
		req = req.SellMetadata(cfg.SellMetadata)
	}

	if cfg.SellMinQuantity != "" {
		req = req.SellMinQuantity(cfg.SellMinQuantity)
	}

	if cfg.SellTokenAddress != "" {
		req = req.SellTokenAddress(cfg.SellTokenAddress)
	}

	if cfg.SellTokenID != "" {
		req = req.SellTokenId(cfg.SellTokenID)
	}

	if cfg.SellTokenName != "" {
		req = req.SellTokenName(cfg.SellTokenName)
	}

	if cfg.SellTokenType != "" {
		req = req.SellTokenType(cfg.SellTokenType)
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

	return &req
}
