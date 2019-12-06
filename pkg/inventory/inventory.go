/*
Copyright © 2019 Michael Gruener

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package inventory

import (
	"fmt"
	"regexp"

	invconfig "github.com/bedag/kusible/pkg/inventory/config"
	"github.com/bedag/kusible/pkg/values"
	"github.com/mitchellh/mapstructure"
)

func NewInventory(path string, ejson values.EjsonSettings, skipKubeconfig bool) (*Inventory, error) {
	// load the raw inventory yaml data
	raw, err := values.New(path, []string{}, false, ejson)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = mapstructure.Decode(raw.Map(), &data)
	if err != nil {
		return nil, err
	}

	// parse the yaml data into the inventory config
	inventoryConfig, err := invconfig.NewConfigFromMap(&data)
	if err != nil {
		return nil, fmt.Errorf("failed load inventory config: %s", err)
	}

	// create the inventory based on the inventory config
	entries := make(map[string]*Entry, len(inventoryConfig.Inventory))
	for _, entryConf := range inventoryConfig.Inventory {
		entry, err := NewEntryFromConfig(entryConf)
		if err != nil {
			return nil, fmt.Errorf("failed to create entry '%s' from config: %s", entryConf.Name, err)
		}
		entries[entryConf.Name] = entry
	}

	return &Inventory{entries: entries, ejson: &ejson}, nil
}

func (i *Inventory) Entries() map[string]*Entry {
	return i.entries
}

func (i *Inventory) EntryNames(filter string, limits []string) ([]string, error) {
	var result []string

	regex, err := regexp.Compile("^" + filter + "$")
	if err != nil {
		return nil, fmt.Errorf("inventory entry filter '%s' is not a valid regex: %s", filter, err)
	}

	for _, entry := range i.entries {
		if regex.MatchString(entry.name) {
			valid, err := entry.MatchLimits(limits)
			if err != nil {
				return nil, err
			}
			if valid {
				result = append(result, entry.name)
			}
		}
	}
	return result, nil
}
