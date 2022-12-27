package collections

import (
	"encoding/json"
	"fmt"

	"github.com/deadloct/immutablex-cli/lib"
	"github.com/immutable/imx-core-sdk-golang/imx/api"
	log "github.com/sirupsen/logrus"
)

func PrintCollection(collection *api.Collection) {
	data, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		log.Debugf("could not convert asset to json: %v\nasset: %#v\n", err, collection)
		return
	}

	fmt.Println(string(data))
}

func PrintCollections(collections []api.Collection, detailed bool) {
	for _, col := range collections {
		if detailed {
			PrintCollection(&col)
		} else {
			fmt.Printf("%s: %s\n", col.Name, lib.ImmutascanURL+col.Address)
		}
	}
}
