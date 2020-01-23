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
	"bufio"
	"os"

	"github.com/standupdev/strset"
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
	return NewBaseConfigFromReader(bufio.NewReader(file))
}

// NewBaseConfigFromReader loads a playbook base config from the given
// bufio.Reader. The reader must point to yaml data.
func NewBaseConfigFromReader(reader *bufio.Reader) (*BaseConfig, error) {
	data := []byte{}
	var result BaseConfig
	_, err := reader.Read(data)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (bc *BaseConfig) FilterApplicable(targetGroups []string) *BaseConfig {
	targetSet := strset.Make(targetGroups...)

	applicablePlays := []*BasePlay{}

	for _, play := range bc.Plays {
		playSet := strset.Make(play.Groups...)
		if playSet.Len() > targetSet.Len() {
			continue
		}
		if playSet.SubsetOf(targetSet) {
			applicablePlays = append(applicablePlays, play)
		}
	}

	return &BaseConfig{Plays: applicablePlays}
}
