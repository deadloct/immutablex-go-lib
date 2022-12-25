package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/immutable/imx-core-sdk-golang/imx/api"
	log "github.com/sirupsen/logrus"
)

type AssetManager struct {
	client IMXClientWrapper
}

func NewAssetManager() *AssetManager {
	return &AssetManager{
		client: NewClient(),
	}
}

func (am *AssetManager) Start() error {
	return am.client.Start()
}

func (am *AssetManager) Stop() {
	am.client.Stop()
}

func (am *AssetManager) GetAsset(ctx context.Context, tokenAddress, tokenID string, includeFees bool) (*api.Asset, error) {
	log.Debugf("fetching asset id %s from collection %s (with fees:%b)", tokenAddress, tokenID, includeFees)
	return am.client.GetClient().GetAsset(ctx, tokenAddress, tokenID, &includeFees)
}

type GetAssetsRequest struct {
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

func (am *AssetManager) GetAssets(
	ctx context.Context,
	cfg *GetAssetsRequest,
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
		return am.GetAssets(ctx, cfg)
	}

	// Attempt to fetch earlier assets
	if len(resp.Result) > 0 {
		cfg.Before = last
		return am.GetAssets(ctx, cfg)
	}

	return cfg.Assets, nil
}

func (am *AssetManager) PrintAsset(asset *api.Asset) {
	data, err := json.MarshalIndent(asset, "", "  ")
	if err != nil {
		log.Debugf("could not convert asset to json: %v\nasset: %#v\n", err, asset)
		return
	}

	fmt.Println(string(data))
}

func (am *AssetManager) PrintAssets(collectionAddr string, assets []api.AssetWithOrders) {
	for _, asset := range assets {
		name := "[no name set]"
		status := asset.Status
		id := *asset.Id

		if asset.Name.IsSet() && asset.Name.Get() != nil {
			name = *asset.Name.Get()
		}

		if status == "" {
			status = "[no status set]"
		}

		if id == "" {
			id = "[no id set]"
		}

		fmt.Printf("%s (Status: %v): (%s)\n", name, status, path.Join(ImmutascanURL, collectionAddr, asset.TokenId))
	}
}

func (am *AssetManager) getAPIListAssetsRequest(ctx context.Context, cfg *GetAssetsRequest) api.ApiListAssetsRequest {
	req := am.client.GetClient().NewListAssetsRequest(ctx).
		Collection(cfg.Collection).
		PageSize(MaxAssetsPerReq)

	if cfg.BuyOrders {
		req = req.BuyOrders(cfg.BuyOrders)
	}

	if cfg.Direction != "" {
		req = req.Direction(cfg.Direction)
	}

	if cfg.IncludeFees {
		req = req.IncludeFees(cfg.IncludeFees)
	}

	if cfg.Metadata != nil {
		if data := am.parseMetadata(cfg.Metadata); data != "" {
			req = req.Metadata(data)
		}
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

	// Recursion helpers
	if cfg.Before != "" {
		req = req.UpdatedMaxTimestamp(cfg.Before)
	}

	if cfg.Cursor != "" {
		req = req.Cursor(cfg.Cursor)
	}

	return req
}

func (am *AssetManager) parseMetadata(metadata []string) string {
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
