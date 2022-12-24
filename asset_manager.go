package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/immutable/imx-core-sdk-golang/imx/api"
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

func (am *AssetManager) GetAsset(ctx context.Context, collectionAddr, id string) (*api.Asset, error) {
	includeFees := false
	return am.client.GetClient().GetAsset(ctx, collectionAddr, id, &includeFees)
}

// Options from https://docs.x.immutable.com/reference/#/operations/listAssets
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

	req := am.GetAPIListAssetsRequest(ctx, cfg)
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
	log.Printf("fetched %v assets from %v to %v\n", len(resp.Result), first, last)

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

func (am *AssetManager) GetAPIListAssetsRequest(ctx context.Context, cfg *GetAssetsRequest) api.ApiListAssetsRequest {
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

func (am *AssetManager) PrintAsset(asset *api.Asset) {
	fmt.Println(FormatAssetInfo(asset))
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

		fmt.Printf("%s: %v (%s)\n", name, status, path.Join(ImmutascanURL, collectionAddr, asset.TokenId))
	}
}

func (am *AssetManager) PrintAssetCounts(name string, assets []api.AssetWithOrders) {
	counts := make(map[string]int, 4)
	for _, asset := range assets {
		rarity, ok := asset.Metadata["Rarity"].(string)
		if !ok {
			log.Printf("asset %s skipped because it doesn't have a rarity\n", asset.TokenId)
			continue
		}

		if !asset.Name.IsSet() {
			log.Printf("asset %s skipped since it has no name and must be messed up\n", asset.TokenId)
		}

		counts[rarity]++
		counts["Total"]++
	}

	fmt.Println(FormatAssetCounts(name, counts))
}

func (am *AssetManager) parseMetadata(metadata []string) string {
	metamap := make(map[string][]string, len(metadata))
	for _, item := range metadata {
		parts := strings.SplitN(item, "=", 2)
		if len(parts) != 2 {
			log.Printf("could not parse metadata item %s into a key=value pair", item)
			continue
		}

		metamap[parts[0]] = append(metamap[parts[0]], parts[1])
	}

	data, err := json.Marshal(metamap)
	if err != nil {
		log.Printf("skipping metamata completely because it could not be converted to json: %v", err)
	}

	return string(data)
}
