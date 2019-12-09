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

package values

import (
	"bytes"

	"github.com/geofffranks/spruce"
	// TODO switch to "sigs.k8s.io/yaml"
	"github.com/geofffranks/yaml"
)

/*
YamlString returns the map data as yaml encoded string
*/
func (d *data) YAML() ([]byte, error) {
	yaml, err := yaml.Marshal(d)
	if err != nil {
		return nil, err
	}
	return yaml, nil
}

/*
JsonString returns the map data as yaml encoded string
*/
func (d *data) JSON() ([]byte, error) {
	// Although we want to create a json string, first convert
	// the data to yaml as there is no easy way to convert
	// a map that can have non-string keys to json. Then
	// convert the yaml data to json with the help of spruce
	yaml, err := yaml.Marshal(d)
	if err != nil {
		return nil, err
	}
	json, err := spruce.JSONifyIO(bytes.NewReader(yaml), false)
	if err != nil {
		// spruce errors can contain ansi colors
		return nil, StripAnsiError(err)
	}
	return []byte(json), nil
}
