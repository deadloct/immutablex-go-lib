package lib

import (
	"fmt"

	"github.com/immutable/imx-core-sdk-golang/imx/api"
)

const ImmutascanURL = "https://immutascan.io/address/"

func FormatAssetCounts(name string, counts map[string]int) string {
	var str = `
%s:
- Common: %d
- Rare: %d
- Epic: %d
- Legendary: %d
- Mythic: %d
- Total: %d
`
	return fmt.Sprintf(
		str,
		name,
		counts["Common"],
		counts["Rare"],
		counts["Epic"],
		counts["Legendary"],
		counts["Mythic"],
		counts["Total"],
	)
}

func FormatAssetInfo(asset *api.Asset) string {
	ownerURL := ImmutascanURL + asset.User
	str := `
%v:
- Background: %v
- Eyes: %v
- Frame: %v
- Gender: %v
- Generation: %v
- Hair: %v
- Hat: %v
- Outfit: %v
- Rarity: %v
- Skin: %v
- Description: %v
- Image URL: %v
- Game Meta JSON: %v
- Owner: %v
`

	return fmt.Sprintf(str,
		asset.Metadata["name"],
		asset.Metadata["Background"],
		asset.Metadata["Eye"],
		asset.Metadata["Frame"],
		asset.Metadata["Gender"],
		asset.Metadata["Generation"],
		asset.Metadata["Hair"],
		asset.Metadata["Hat"],
		asset.Metadata["Outfit"],
		asset.Metadata["Rarity"],
		asset.Metadata["Skin"],
		asset.Metadata["description"],
		asset.Metadata["image"],
		ownerURL,
		asset.Metadata["game_meta"],
	)
}
