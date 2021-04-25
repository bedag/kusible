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

package playbook

import (
	"bufio"
	"fmt"
	"os"

	"github.com/bedag/kusible/pkg/playbook/config"
	"github.com/bedag/kusible/pkg/target"
)

/*
Each run-relevant inventory entry (target) has its own "view" on the given playbook containting
only the relevant plays for its groups.

Given a list of targets, the playbook loader

* loads the playbook (without evaluation)
* for each target
	* filters the plays based on the given groups
	* retrieves the cluster-inventory of the target (optional)
	* loads the values relevant for the given groups (without evaluation)
	* merges the cluster-inventory data, the target values and the filtered playbook
	* evaluates the result
	* unmarshalls the merged/evaluated playbook/value map into a valid playbook config structure
*/

func NewSet(path string, targets *target.Targets, skipEval bool, skipClusterInv bool) (Set, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return NewSetFromReader(bufio.NewReader(file), targets, skipEval, skipClusterInv)
}

func NewSetFromReader(reader *bufio.Reader, targets *target.Targets, skipEval bool, skipClusterInv bool) (Set, error) {
	// Get the base config of the given playbook
	// The base config contains all playbook data but only the name and groups of
	// each play are required and parsed. We need the groups of the plays
	// to generate a unique set of plays applicable to each target
	baseConfig, err := config.NewBaseConfigFromReader(reader)
	if err != nil {
		return nil, err
	}

	result := make(Set)

	for _, target := range targets.Targets() {

		playbook, err := New(baseConfig, target, skipEval, skipClusterInv)
		if err != nil {
			return nil, fmt.Errorf("Failed to create playbook for target '%s': '%s'", target.Entry().Name(), err)
		}

		result[target.Entry().Name()] = playbook
	}

	return result, nil
}
