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

import "github.com/mitchellh/mapstructure"

// Config holds the list of entries that can serve as targets
// for application deployments. It is the root of the actual inventory
type Config struct {
	Inventory []*Entry `yaml:"inventory"`
}

// Entry is a single inventory entry representing a possible deploy
// target.
type Entry struct {
	// Name to uniquely identify the entry
	Name string `yaml:"name"`
	// Groups is a "least specific to most specific" ordered list of
	// groups associated with this entry, used to compile the "values"
	// for this entry and as selector to target entries in specific groups.
	// Each entry is always part of the "all" group and a group
	// with the name of the entry.
	Groups []string `yaml:"groups"`
	// ConfigNamespace refers to the namespace where the "cluster-inventory"
	// ConfigMap can be found
	ConfigNamespace string `yaml:"config_namespace"`
	// Kubeconfig holds the kubeconfig loader configuration
	Kubeconfig *Kubeconfig `yaml:"kubeconfig"`
}

// Kubeconfig holds information on how / where to retrieve / generate
// the Kubeconfig for an entry in the inventory
type Kubeconfig struct {
	Backend string  `yaml:"backend"`
	Params  *Params `yaml:"params"`
}

// Params holds the parameters used by a kubeconfig backend to
// retrieve / generate a kubeconfig. The exact fields depend
// on the kubeconfig loader.
type Params map[string]interface{}

// decode the given data with the default decoder settings
func decode(data *map[string]interface{}, result interface{}) error {
	decoderConfig := &mapstructure.DecoderConfig{
		ZeroFields:       true,
		ErrorUnused:      false,
		WeaklyTypedInput: true,
		TagName:          "yaml",
		Result:           &result,
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return err
	}
	err = decoder.Decode(data)
	return err
}

// NewConfigFromMap takes raw config data and parses it into an
// inventory config
func NewConfigFromMap(data *map[string]interface{}) (*Config, error) {
	var config Config
	err := decode(data, &config)
	return &config, err
}

// NewConfig returns an empty inventory config
func NewConfig() *Config {
	return &Config{
		Inventory: make([]*Entry, 0),
	}
}
