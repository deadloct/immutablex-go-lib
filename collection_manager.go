package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

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

func (c *CollectionManager) GetShortcutByName(name string) *CollectionShortcut {
	v, ok := c.shortcuts[name]
	if !ok {
		return nil
	}

	return &v
}

func (c *CollectionManager) loadShortcuts() {
	if ShortcutLocation == "" {
		return
	}

	if _, err := os.Stat(ShortcutLocation); err != nil {
		log.Printf("skipping shortcuts load because shortcut location %s does not exist", ShortcutLocation)
	}

	content, err := ioutil.ReadFile(ShortcutLocation)
	if err != nil {
		log.Printf("could not load shortcuts file %s: %v", ShortcutLocation, err)
	}

	var data []CollectionShortcut
	if err = json.Unmarshal(content, &data); err != nil {
		log.Printf("could not parse shortcuts file %s: %v", ShortcutLocation, err)
	}

	c.shortcuts = make(map[string]CollectionShortcut, len(data))
	for _, shortcut := range data {
		c.shortcuts[shortcut.Shortcut] = shortcut
	}
}
