package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/immutable/imx-core-sdk-golang/imx/api"
	log "github.com/sirupsen/logrus"
)

var DefaultShortcutsContent = []byte(`
[
    {
        "name": "BitVerse Portals",
        "addr": "0xe4ac52f4b4a721d1d0ad8c9c689df401c2db7291",
        "shortcut": "portal"
    },
    {
        "name": "BitVerse Heroes",
        "addr": "0x6465ef3009f3c474774f4afb607a5d600ea71d95",
        "shortcut": "hero"
    }
]
`)

var ShortcutLocation string

func init() {
	ShortcutLocation = os.Getenv("IMX_SHORTCUT_LOCATION")
}

type CollectionShortcut struct {
	Name     string `json:"name"`
	Addr     string `json:"addr"`
	Shortcut string `json:"shortcut"`
}

type CollectionManager struct {
	client    IMXClientWrapper
	shortcuts map[string]CollectionShortcut
}

func NewCollectionManager() *CollectionManager {
	cm := &CollectionManager{client: NewClient()}
	cm.loadShortcuts()
	return cm
}

func (c *CollectionManager) Start() error {
	return c.client.Start()
}

func (c *CollectionManager) Stop() {
	c.client.Stop()
}

func (c *CollectionManager) GetCollection(ctx context.Context, collection string) (*api.Collection, error) {
	log.Debugf("fetching collection %s", collection)
	return c.client.GetClient().GetCollection(ctx, collection)
}

type ListCollectionsConfig struct {
	Blacklist string
	Direction string
	Keyword   string
	OrderBy   string
	Whitelist string

	// Used internally for recursion
	Collections []api.Collection
	Cursor      string
}

func (c *CollectionManager) ListCollections(ctx context.Context, cfg *ListCollectionsConfig) ([]api.Collection, error) {
	req := c.getAPIListCollectionsRequest(ctx, cfg)

	resp, err := c.client.GetClient().ListCollections(req)
	if err != nil {
		return nil, err
	}

	cfg.Collections = append(cfg.Collections, resp.Result...)
	cfg.Cursor = resp.Cursor

	first := *resp.Result[0].UpdatedAt.Get()
	last := *resp.Result[len(resp.Result)-1].UpdatedAt.Get()
	log.Debugf("fetched %v assets from %v to %v", len(resp.Result), first, last)

	if resp.Remaining > 0 {
		return c.ListCollections(ctx, cfg)
	}

	return cfg.Collections, nil
}

func (c *CollectionManager) GetShortcutByName(name string) *CollectionShortcut {
	v, ok := c.shortcuts[name]
	if !ok {
		return nil
	}

	return &v
}

func (c *CollectionManager) PrintCollection(collection *api.Collection) {
	data, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		log.Debugf("could not convert asset to json: %v\nasset: %#v\n", err, collection)
		return
	}

	fmt.Println(string(data))
}

func (c *CollectionManager) PrintCollections(collections []api.Collection, detailed bool) {
	for _, col := range collections {
		if detailed {
			c.PrintCollection(&col)
		} else {
			fmt.Printf("%s: %s\n", col.Name, ImmutascanURL+col.Address)
		}
	}
}

func (c *CollectionManager) getAPIListCollectionsRequest(ctx context.Context, cfg *ListCollectionsConfig) *api.ApiListCollectionsRequest {
	req := c.client.GetClient().NewListCollectionsRequest(ctx).PageSize(MaxAssetsPerReq)

	if cfg.Blacklist != "" {
		req = req.Blacklist(cfg.Blacklist)
	}

	if cfg.Direction != "" {
		req = req.Direction(cfg.Direction)
	}

	if cfg.Keyword != "" {
		req = req.Keyword(cfg.Keyword)
	}

	if cfg.OrderBy != "" {
		req = req.OrderBy(cfg.OrderBy)
	}

	if cfg.Whitelist != "" {
		req = req.Whitelist(cfg.Whitelist)
	}

	return &req
}

func (c *CollectionManager) loadShortcuts() {
	content := DefaultShortcutsContent
	if ShortcutLocation != "" {
		if _, err := os.Stat(ShortcutLocation); err == nil {
			content, err = ioutil.ReadFile(ShortcutLocation)
			if err != nil {
				log.Debugf("could not load shortcuts file %s: %v", ShortcutLocation, err)
				content = DefaultShortcutsContent
			}
		} else {
			log.Debugf("could not stat shortcuts file %s: %v", ShortcutLocation, err)
		}
	}

	var data []CollectionShortcut
	if err := json.Unmarshal(content, &data); err != nil {
		log.Debugf("could not parse shortcuts file %s: %v", ShortcutLocation, err)
	}

	c.shortcuts = make(map[string]CollectionShortcut, len(data))
	for _, shortcut := range data {
		c.shortcuts[shortcut.Shortcut] = shortcut
	}
}
