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

package inventory

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/bedag/kusible/pkg/groups"
	"github.com/mitchellh/mapstructure"
)

func NewEntryFromParams(params map[string]interface{}) (*entry, error) {
	var entry entry

	if _, ok := params["name"].(string); !ok {
		return nil, fmt.Errorf("inventory entry has no name")
	}

	name := params["name"].(string)

	hook := kubeconfigDecoderHookFunc(name)
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: hook,
		Result:     &entry,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return nil, err
	}
	err = decoder.Decode(params)
	if err != nil {
		return nil, err
	}

	entry.Groups = append(entry.Groups, name)
	entry.Groups = append([]string{"all"}, entry.Groups...)

	if entry.ConfigNamespace == "" {
		entry.ConfigNamespace = "kube-config"
	}

	if entry.Kubeconfig == nil {
		config := map[string]interface{}{
			"_entry": name,
		}
		loader := NewKubeconfigS3LoaderFromParams(config)
		if err != nil {
			return &entry, fmt.Errorf("failed to create default loader for entry %s: %s", name, err)
		}
		entry.Kubeconfig, err = NewKubeconfigFromLoader(loader)
		if err != nil {
			return nil, fmt.Errorf("failed to load kubeconfig for entry '%s': %s", name, err)
		}
	}

	return &entry, nil
}

// MatchLimits returns true if the groups of the inventory entry satisfy all given
// limits, which are treated as ^$ enclosed regex
func (e *entry) MatchLimits(limits []string) (bool, error) {
	// no limits -> all groups are valid
	if len(limits) <= 0 {
		return true, nil
	}

	// no groups -> no limit matches
	if len(e.Groups) <= 0 {
		return false, nil
	}

	for _, limit := range limits {
		regex, err := regexp.Compile("^" + limit + "$")
		if err != nil {
			return false, err
		}

		matched := false
		for _, group := range e.Groups {
			if regex.MatchString(group) {
				matched = true
				break
			}
		}
		if !matched {
			return false, nil
		}
	}
	return true, nil
}

// ValidGroups returns all groups of the inventory entry that satisfy at
// least one limit
func (e *entry) ValidGroups(limits []string) ([]string, error) {
	return groups.LimitGroups(e.Groups, limits)
}

func entryDecoderHookFunc(skipKubeconfig bool) mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t.Name() == "entry" {
			var params map[string]interface{}
			err := mapstructure.Decode(data, &params)
			if err != nil {
				return data, err
			}

			entry, err := NewEntryFromParams(params)
			if err != nil {
				return data, err
			}
			return entry, nil
		}
		return data, nil
	}
}
