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

/*
Package config implements the playbook config format
*/
package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"sigs.k8s.io/yaml"
)

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

func (bc *BaseConfig) GetApplicable(targetGroups []string) (*BaseConfig, error) {
	if len(targetGroups) <= 0 {
		return bc, nil
	}

	result := []*BasePlay{}

	for _, play := range bc.Plays {
		v := Validator{}
		for _, expr := range play.Groups {
			pattern, err := NewPattern(expr, targetGroups)
			if err != nil {
				return nil, fmt.Errorf("failed to parse pattern expression '%s': %s", expr, err)
			}
			v.Add(pattern)
		}
		if v.Valid() {
			result = append(result, play)
		}

	}

	return &BaseConfig{Plays: result}, nil
}
