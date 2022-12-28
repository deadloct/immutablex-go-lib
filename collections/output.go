package collections

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/deadloct/immutablex-cli/lib"
	"github.com/immutable/imx-core-sdk-golang/imx/api"
	log "github.com/sirupsen/logrus"
)

func PrintCollectionJSON(collection *api.Collection) {
	data, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		log.Debugf("could not convert asset to json: %v\nasset: %#v\n", err, collection)
		return
	}

	fmt.Println(string(data))
}

func PrintCollectionStandard(collection *api.Collection) {
	url := strings.Join([]string{lib.ImmutascanURL, "address", collection.Address}, "/")
	fmt.Printf("%s: %s\n", collection.Name, url)
}

func PrintCollection(collection *api.Collection, output string) {
	switch strings.ToLower(output) {
	case "json":
		PrintCollectionJSON(collection)
	default:
		PrintCollectionStandard(collection)
	}
}

func PrintCollections(collections []api.Collection, output string) {
	for _, col := range collections {
		PrintCollection(&col, output)
	}
}
