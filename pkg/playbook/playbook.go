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

package config

import (
	"bufio"
	"os"

	config "github.com/bedag/kusible/pkg/playbook/config"
	"github.com/bedag/kusible/pkg/spruce"
	"github.com/bedag/kusible/pkg/target"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

/*
Each run-relevant inventory entry has its own "view" on the given playbook containting
only the relevant plays for its groups.

Given a list of groups, the playbook loader

* loads the playbook (without evaluation)
* filters the plays based on the given groups
* loads the values relevant for the given groups (without evaluation)
* merges the filtered playbook and values
* evaluates the result
* unmarshalls the merged/evaluated playbook/value map into a valid playbook config structure
*/

func New(path string, targets *target.Targets, skipEval bool) ([]*config.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return NewFromReader(bufio.NewReader(file), targets, skipEval)
}

func NewFromReader(reader *bufio.Reader, targets *target.Targets, skipEval bool) ([]*config.Config, error) {
	// Get the base config of the given playbook
	// The base config contains all playbook data but only the name and groups of
	// each play are required and parsed. We need the groups of the plays
	// to generate a unique set of plays applicable to each target
	baseConfig, err := config.NewBaseConfigFromReader(reader)
	if err != nil {
		return nil, err
	}

	for _, target := range targets.Targets() {
		// TODO: this should be the data of the cluster inventory ConfigMap
		result := map[string]interface{}{}
		// Based on the groups of the target and the groups of each play,
		// generate a new base config containing only the plays relevant
		// for the current target
		applicableBaseConfig := baseConfig.FilterApplicable(target.Entry().Groups())

		// convert the base config to a simple map to perpare the merge of the
		// remaining plays with the cluster inventory config and the target values
		var playbookMap map[string]interface{}
		data, err := yaml.Marshal(applicableBaseConfig)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(data, &playbookMap)
		if err != nil {
			return nil, err
		}

		values := target.Values().Map()
		err = mergo.Merge(&result, playbookMap, mergo.WithOverride)
		if err != nil {
			return nil, err
		}
		err = mergo.Merge(&result, values, mergo.WithOverride)
		if err != nil {
			return nil, err
		}
		err = spruce.Eval(&result, skipEval, []string{})

		// TODO: Marshal back to a complete playbook datastructure and append
		// to the list of playbook configs for this playbook
	}

	return nil, nil
}
