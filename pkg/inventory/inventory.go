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
	"reflect"
	"regexp"

	"github.com/bedag/kusible/pkg/values"
	"github.com/mitchellh/mapstructure"
)

func NewInventory(path string, ejson values.EjsonSettings, skipKubeconfig bool) (*Inventory, error) {
	raw, err := values.NewValues(path, []string{}, false, ejson)
	if err != nil {
		return nil, err
	}
	data := raw.Map()

	var inventory Inventory

	hook := entriesDecoderHookFunc(skipKubeconfig)
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: hook,
		Result:     &inventory,
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return nil, err
	}
	err = decoder.Decode(data)
	inventory.Ejson = &ejson
	return &inventory, err
}

func entriesDecoderHookFunc(skipKubeconfig bool) mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t.Name() == "entries" {
			var entriesList []entry

			hook := entryDecoderHookFunc(skipKubeconfig)
			decoderConfig := &mapstructure.DecoderConfig{
				DecodeHook: hook,
				Result:     &entriesList,
			}
			decoder, err := mapstructure.NewDecoder(decoderConfig)
			if err != nil {
				return data, err
			}
			err = decoder.Decode(data)
			if err != nil {
				return data, err
			}
			entriesMap := make(entries, len(entriesList))
			for _, entry := range entriesList {
				name := entry.Name
				entriesMap[name] = entry
			}
			return entriesMap, nil
		}
		return data, nil
	}
}

func (i *Inventory) EntryNames(filter string, limits []string) ([]string, error) {
	var result []string

	regex, err := regexp.Compile("^" + filter + "$")
	if err != nil {
		return nil, fmt.Errorf("inventory entry filter '%s' is not a valid regex: %s", filter, err)
	}

	for _, entry := range i.Entries {
		if regex.MatchString(entry.Name) {
			valid, err := entry.MatchLimits(limits)
			if err != nil {
				return nil, err
			}
			if valid {
				result = append(result, entry.Name)
			}
		}
	}
	return result, nil
}
