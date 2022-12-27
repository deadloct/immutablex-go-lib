package assets

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/deadloct/immutablex-cli/lib"
	"github.com/immutable/imx-core-sdk-golang/imx/api"
	log "github.com/sirupsen/logrus"
)

func PrintAsset(asset *api.Asset) {
	data, err := json.MarshalIndent(asset, "", "  ")
	if err != nil {
		log.Debugf("could not convert asset to json: %v\nasset: %#v\n", err, asset)
		return
	}

	fmt.Println(string(data))
}

func PrintAssets(collectionAddr string, assets []api.AssetWithOrders) {
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

		fmt.Printf("%s (Status: %v): (%s)\n", name, status, path.Join(lib.ImmutascanURL, collectionAddr, asset.TokenId))
	}
}
