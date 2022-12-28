package assets

import (
	"encoding/json"
	"fmt"
	"strings"

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
		name := asset.GetName()
		status := asset.Status
		id := *asset.Id

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
			asset.TokenId,
		}, ",")
		fmt.Printf("%s (Status: %v): (%s)\n", name, status, url)
	}
}
