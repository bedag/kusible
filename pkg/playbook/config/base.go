/*
Copyright © 2021 Bedag Informatik AG

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

package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"sigs.k8s.io/yaml"
)

// TODO
// Fields of map[string]interface{} can be easyly reached by doing
// x := mymap["field"].(type)
// Duh!

// NewBaseConfigFromFile loads a playbook base config from the given
// file path. The file must contain yaml data.
func NewBaseConfigFromFile(path string) (*BaseConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return NewBaseConfigFromReader(file)
}

// NewBaseConfigFromReader loads a playbook base config from the given
// bufio.Reader. The reader must point to yaml data.
func NewBaseConfigFromReader(reader io.Reader) (*BaseConfig, error) {
	data := []byte{}
	var result BaseConfig
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Applicable returns a BaseConfig that contains only the plays where
// the groups value of the play matches the groups given as parameter
func (bc *BaseConfig) Applicable(groups []string) (*BaseConfig, error) {
	if len(groups) <= 0 {
		return bc, nil
	}

	result := []*BasePlay{}

	for _, play := range bc.Plays {
		v := Validator{}
		for _, expr := range play.Groups {
			pattern, err := NewPattern(expr, groups)
			if err != nil {
				return nil, fmt.Errorf("failed to parse pattern expression '%s' of play '%s': %s", expr, play.Name, err)
			}
			v.Add(pattern)
		}
		if v.Valid() {
			result = append(result, play)
		}

	}

	return &BaseConfig{Plays: result}, nil
}

// ApplicableMap returns a map of the BaseConfig that contains only the plays where
// the groups value of the play matches the groups given as parameter
func (bc *BaseConfig) ApplicableMap(groups []string) (*map[string]interface{}, error) {
	config, err := bc.Applicable(groups)
	if err != nil {
		return nil, err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
