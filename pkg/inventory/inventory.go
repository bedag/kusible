// Copyright © 2019 Michael Gruener
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package inventory

import (
	"fmt"
	"os"
	"reflect"
	"regexp"

	"github.com/bedag/kusible/pkg/groups"
	"github.com/bedag/kusible/pkg/values"
	"github.com/mitchellh/mapstructure"
)

func NewInventory(path string, ejson values.EjsonSettings, skipKubeconfig bool) (*Inventory, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	var data map[interface{}]interface{}
	if stat.Mode().IsRegular() {
		// the path provided is a file, treat it as a single value
		// file, thus loading it with ejson and spruc operator support
		inventoryFile := values.NewValueFile(path, false, ejson)
		data, err = inventoryFile.LoadMap()
		if err != nil {
			return nil, err
		}
	} else {
		// The path provided is a directory, treat it as a values
		// directory. As the valuesDirectory type requires a list
		// of groups to determine which files to process, first
		// get a list of all groups in the given directory
		groups, err := groups.Groups(path, ".*", []string{})
		if err != nil {
			return nil, err
		}
		inventoryDir := values.NewValuesDirectory(path, groups, false, ejson)
		data, err = inventoryDir.LoadMap()
		if err != nil {
			return nil, err
		}
	}

	var inventory Inventory

	hook := loaderDecoderHookFunc(skipKubeconfig)
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: hook,
		Result:     &inventory,
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return nil, err
	}
	err = decoder.Decode(data)
	return &inventory, err
}

func loaderDecoderHookFunc(skipKubeconfig bool) mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t.Name() == "kubeconfig" {
			dataMap := data.(map[interface{}]interface{})

			var backend string
			var params map[interface{}]interface{}
			for k, v := range dataMap {
				key := k.(string)
				if key == "backend" {
					backend = v.(string)
				}
				if key == "params" {
					params = v.(map[interface{}]interface{})
				}
			}

			kubeconfig, err := NewKubeconfigFromConfig(backend, params, skipKubeconfig)
			return kubeconfig, err
		}
		return data, nil
	}
}

func (i *Inventory) EntryNames(filter string) ([]string, error) {
	var result []string

	regex, err := regexp.Compile("^" + filter + "$")
	if err != nil {
		return nil, fmt.Errorf("inventory entry filter '%s' is not a valid regex: %s", filter, err)
	}

	for _, entry := range i.Entries {
		if regex.MatchString(entry.Name) {
			result = append(result, entry.Name)
		}
	}
	return result, nil
}
