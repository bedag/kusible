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
	"testing"

	"github.com/bedag/kusible/internal/wrapper/ejson"
	"github.com/bedag/kusible/pkg/inventory"
	"github.com/bedag/kusible/pkg/target"
	"gotest.tools/assert"
)

func TestPlaybook(t *testing.T) {
	tests := map[string]struct {
		inventory string
		vars      string
		playbook  string
	}{
		"simple": {
			inventory: "testdata/simple/inventory.yml",
			vars:      "testdata/simple/group_vars",
			playbook:  "testdata/simple/playbook.yml",
		},
		"complex": {
			inventory: "testdata/complex/inventory.yml",
			vars:      "testdata/complex/group_vars",
			playbook:  "testdata/complex/playbook.yml",
		},
	}

	// skip spruce eval for target values, as this happens later
	// during the playbook creation
	skipEval := true
	ejsonSettings := ejson.Settings{}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			inv, err := inventory.NewInventory(tc.inventory, ejsonSettings, true)
			assert.NilError(t, err)

			targets, err := target.NewTargets(".*", []string{}, tc.vars, inv, skipEval, &ejsonSettings)
			assert.NilError(t, err)

			playbooks, err := New(tc.playbook, targets, false)
			assert.NilError(t, err)
			assert.Equal(t, len(targets.Targets()), len(playbooks))
		})
	}
}
