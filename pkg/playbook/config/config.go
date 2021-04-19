/*
Copyright © 2019 Copyright © 2021 Bedag Informatik AG

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

/*
Package config implements the playbook config format
*/
package config

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
	"sigs.k8s.io/yaml"
)

// decode the given data with the default decoder settings
func decode(data *map[string]interface{}, result interface{}) error {
	// TODO: check https://github.com/mitchellh/mapstructure/issues/187 to
	// support mitchellh/mapstructure > 1.3.1
	decoderConfig := &mapstructure.DecoderConfig{
		ZeroFields:       true,
		ErrorUnused:      false,
		WeaklyTypedInput: true,
		TagName:          "json",
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

// NewConfig returns an empty playbook config
func NewConfig() *Config {
	return &Config{
		Plays: make([]*Play, 0),
	}
}

func (c *Config) YAML() ([]byte, error) {
	return yaml.Marshal(c)
}

func (c *Config) JSON() ([]byte, error) {
	return json.Marshal(c)
}
