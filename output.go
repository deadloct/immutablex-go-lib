package lib

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/immutable/imx-core-sdk-golang/imx/api"
)

const ImmutascanURL = "https://immutascan.io/address/"

func FormatAssetInfo(asset *api.Asset) string {
	// ownerURL := ImmutascanURL + asset.User

	data, err := json.MarshalIndent(asset, "", "  ")
	if err != nil {
		log.Printf("could not stringify asset: %v", err)
		return fmt.Sprintf("%#v\n", asset)
	}

	return string(data)
}
