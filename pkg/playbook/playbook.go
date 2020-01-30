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

package playbook

import (
	"bufio"
	"fmt"
	"os"

	config "github.com/bedag/kusible/pkg/playbook/config"
	"github.com/bedag/kusible/pkg/spruce"
	"github.com/bedag/kusible/pkg/target"
	"github.com/imdario/mergo"
	"sigs.k8s.io/yaml"
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

	result := []*config.Config{}

	for _, target := range targets.Targets() {
		// TODO: this should be the data of the cluster inventory ConfigMap
		var mergeResult map[string]interface{}
		// Based on the groups of the target and the groups of each play,
		// generate a new base config containing only the plays relevant
		// for the current target
		targetBaseConfig, err := baseConfig.GetApplicable(target.Entry().Groups())
		if err != nil {
			return nil, err
		}

		// convert the base config to a simple map to perpare the merge of the
		// remaining plays with the cluster inventory config and the target values
		var playbookMap map[string]interface{}
		data, err := yaml.Marshal(targetBaseConfig)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(data, &playbookMap)
		if err != nil {
			return nil, err
		}

		err = mergo.Merge(&mergeResult, playbookMap, mergo.WithOverride)
		if err != nil {
			return nil, fmt.Errorf("failed merge cluster-inventory and playbook for target '%s': %s", target.Entry().Name(), err)
		}

		values := target.Values().Map()
		err = mergo.Merge(&mergeResult, values, mergo.WithOverride)
		if err != nil {
			return nil, fmt.Errorf("failed merge values and playbook for target '%s': %s", target.Entry().Name(), err)
		}

		err = spruce.Eval(&mergeResult, false, []string{})
		if err != nil {
			// TODO: add optional way to dump the unevaluated yaml here
			//doc, _ := yaml.Marshal(mergeResult)
			//fmt.Printf("%s\n", string(doc))
			return nil, fmt.Errorf("failed evaluate playbook config for target '%s': %s", target.Entry().Name(), err)
		}

		targetConfig, err := config.NewConfigFromMap(&mergeResult)
		if err != nil {
			return nil, fmt.Errorf("failed to create playbook config for target '%s': %s", target.Entry().Name(), err)
		}
		result = append(result, targetConfig)
	}

	return result, nil
}
