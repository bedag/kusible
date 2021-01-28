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
	"encoding/json"
	"fmt"

	"github.com/bedag/kusible/internal/third_party/deepcopy"
	"github.com/bedag/kusible/internal/wrapper/spruce"
	"github.com/bedag/kusible/pkg/playbook/config"
	"github.com/bedag/kusible/pkg/target"
	"github.com/imdario/mergo"
	"sigs.k8s.io/yaml"
)

// New creates a Playbook for one specific target. For a given BaseConfig, each target
// as an individual list of plays, based on the groups of the target and the plays.
func New(baseConfig *config.BaseConfig, target *target.Target, skipEval bool, skipClusterInv bool) (*Playbook, error) {
	// Based on the groups of the target and the groups of each play,
	// generate a new base config containing only the plays relevant
	// for the current target. As we have to merge the result with
	// data structures in the next step, retrieve a map instead of the
	// base config itself
	playbookMap, err := baseConfig.ApplicableMap(target.Entry().Groups())
	if err != nil {
		return nil, fmt.Errorf("failed to get plays: %s", err)
	}

	var mergeResult map[string]interface{}

	if !skipClusterInv {
		clusterInventory, err := target.Entry().ClusterInventory()
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve cluster-inventory: %s", err)
		}

		mergeResult, err = deepcopy.Map(*clusterInventory)
		if err != nil {
			return nil, fmt.Errorf("failed to copy cluster-inventory: %s", err)
		}
	}

	err = mergo.Merge(&mergeResult, playbookMap, mergo.WithOverride)
	if err != nil {
		return nil, fmt.Errorf("failed merge cluster-inventory and playbook: %s", err)
	}

	values := target.Values().Map()
	err = mergo.Merge(&mergeResult, values, mergo.WithOverride)
	if err != nil {
		return nil, fmt.Errorf("failed merge values and playbook: %s", err)
	}

	result := &Playbook{
		Raw: mergeResult,
	}
	if !skipEval {
		err = spruce.Eval(&mergeResult, false, []string{})
		if err != nil {
			// TODO: add optional way to dump the unevaluated yaml here
			//doc, _ := yaml.Marshal(mergeResult)
			//fmt.Printf("%s\n", string(doc))
			return nil, fmt.Errorf("failed evaluate playbook config: %s", err)
		}

		// TODO: if playbookMap does not contain any plays, but
		//       the targetConfig here is not empty, emit at least a warning
		//       because that means we have plays that were not part of
		//       the playbook but part of the values or cluster config map
		//       and therefore if was not tested if the play should actually be
		//       executed for the given inventory entry
		targetConfig, err := config.NewConfigFromMap(&mergeResult)
		if err != nil {
			return nil, fmt.Errorf("failed to create playbook config: %s", err)
		}
		result.Config = targetConfig
	}

	return result, nil
}

func (p *Playbook) YAML(raw bool) ([]byte, error) {
	// we want the raw, unevaluated config
	if raw {
		return yaml.Marshal(p.Raw)
	}

	// we want the evaluated config, but only if it
	// contains at least one play
	if p.Config != nil && (len(p.Config.Plays) > 0) {
		return p.Config.YAML()
	}

	// we neither want the unevaluated config nor do
	// we have an evaluated config that contains at
	// least one play
	return []byte{}, nil
}

func (p *Playbook) JSON(raw bool) ([]byte, error) {
	// we want the raw, unevaluated config
	if raw {
		return json.Marshal(p.Raw)
	}

	// we want the evaluated config, but only if it
	// contains at least one play
	if p.Config != nil && (len(p.Config.Plays) > 0) {
		return p.Config.JSON()
	}

	// we neither want the unevaluated config nor do
	// we have an evaluated config that contains at
	// least one play
	return []byte{}, nil
}

func (p *Playbook) Map(raw bool) (map[string]interface{}, error) {
	// we want the raw, unevaluated config
	if raw {
		return p.Raw, nil
	}

	// we want the evaluated config, but only if it
	// contains at least one play
	result := map[string]interface{}{}
	if p.Config != nil && (len(p.Config.Plays) > 0) {
		mergo.Map(&result, p.Config)
	}

	// we neither want the unevaluated config nor do
	// we have an evaluated config that contains at
	// least one play
	return result, nil
}
