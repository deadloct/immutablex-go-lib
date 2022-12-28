package assets

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/deadloct/immutablex-cli/lib"
	"github.com/immutable/imx-core-sdk-golang/imx/api"
	log "github.com/sirupsen/logrus"
)

func PrintAssetJSON[T any](asset T) {
	data, err := json.MarshalIndent(asset, "", "  ")
	if err != nil {
		log.Debugf("could not convert asset to json: %v\nasset: %#v\n", err, asset)
		return
	}

	fmt.Println(string(data))
}

func printAssetCommon(name, status, id, tokenID, collectionAddr string) {
	if name == "" {
		name = "[no name set]"
	}

	if status == "" {
		status = "[no status set]"
	}

	if id == "" {
		id = "[no id set]"
	}

	url := strings.Join([]string{
		lib.ImmutascanURL,
		"address",
		collectionAddr,
		tokenID,
	}, "/")
	fmt.Printf("%s (Status: %v): (%s)\n", name, status, url)
}

func PrintAssetWithOrdersStandard(collectionAddr string, asset *api.AssetWithOrders) {
	printAssetCommon(
		asset.GetName(),
		asset.Status,
		*asset.Id,
		asset.TokenId,
		collectionAddr,
	)
}

func PrintAssetStandard(collectionAddr string, asset *api.Asset) {
	printAssetCommon(
		asset.GetName(),
		asset.Status,
		*asset.Id,
		asset.TokenId,
		collectionAddr,
	)
}

func PrintAsset(collectionAddr string, asset *api.Asset, output string) {
	switch strings.ToLower(output) {
	case "json":
		PrintAssetJSON(asset)
	default:
		PrintAssetStandard(collectionAddr, asset)
	}
}

func PrintAssets(collectionAddr string, assets []api.AssetWithOrders, output string) {
	for _, asset := range assets {
		switch strings.ToLower(output) {
		case "json":
			PrintAssetJSON(asset)
		default:
			PrintAssetWithOrdersStandard(collectionAddr, &asset)
		}
	}
}
