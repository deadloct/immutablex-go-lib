package collections

import (
	"encoding/json"
	"io/ioutil"
	"os"

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

type Shortcut struct {
	Name     string `json:"name"`
	Addr     string `json:"addr"`
	Shortcut string `json:"shortcut"`
}

type Shortcuts map[string]Shortcut

func NewShortcuts() Shortcuts {
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

	var data []Shortcut
	if err := json.Unmarshal(content, &data); err != nil {
		log.Debugf("could not parse shortcuts file %s: %v", ShortcutLocation, err)
	}

	s := make(map[string]Shortcut, len(data))
	for _, shortcut := range data {
		s[shortcut.Shortcut] = shortcut
	}

	return s
}

func (s Shortcuts) GetShortcutByName(name string) *Shortcut {
	v, ok := s[name]
	if !ok {
		return nil
	}

	return &v
}
