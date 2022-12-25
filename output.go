package lib

import (
	"encoding/json"
	"fmt"

	"github.com/immutable/imx-core-sdk-golang/imx/api"
	log "github.com/sirupsen/logrus"
)

const ImmutascanURL = "https://immutascan.io/address/"

func FormatAssetInfo(asset *api.Asset) string {
	// ownerURL := ImmutascanURL + asset.User

	data, err := json.MarshalIndent(asset, "", "  ")
	if err != nil {
		log.Debugf("could not stringify asset: %v", err)
		return fmt.Sprintf("%#v\n", asset)
	}

	return string(data)
}
