package dry_config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

type DryConfig struct {
	config *ini.File
}

func loadConfig() (*ini.File, error) {
	var path string
	var home = os.Getenv("HOME") + "/.dry/conf"

	if _, err := os.Stat(".dryrc"); err == nil {
		path = ".dryrc"
	} else if _, err := os.Stat(home); err == nil {
		path = home
	} else {
		return nil, fmt.Errorf("no dryrc found")
	}

	return ini.Load(path)
}

func LoadConfig() (*DryConfig, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	return &DryConfig{
		config: config,
	}, nil
}

func (k *DryConfig) Resolve(service, name string) (map[string]string, error) {
	key := fmt.Sprintf("%s %s", service, name)
	if !k.config.HasSection(key) {
		return nil, fmt.Errorf("no such section with key '%s'", key)
	}

	section := k.config.Section(key)

	// if any of the string start with _
	if hasFunctionalKey(section.KeyStrings()) {
		return nil, nil
	} else {
		kv := section.KeysHash()
		return kv, nil
	}
}

func hasFunctionalKey(keys []string) bool {
	for _, key := range keys {
		if strings.HasPrefix(key, "_") {
			return true
		}
	}
	return false
}
