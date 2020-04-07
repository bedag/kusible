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

	"github.com/bedag/kusible/pkg/wrapper/ejson"
	"github.com/bedag/kusible/pkg/inventory"
	"github.com/bedag/kusible/pkg/target"
	"gotest.tools/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestSet(t *testing.T) {
	tests := map[string]struct {
		inventory      string
		vars           string
		playbook       string
		skipEval       bool
		skipClusterInv bool
		hasConfig      bool
		hasClusterData bool
	}{
		"simple": {
			inventory:      "testdata/simple/inventory.yml",
			vars:           "testdata/simple/group_vars",
			playbook:       "testdata/simple/playbook.yml",
			skipEval:       false,
			skipClusterInv: false,
			hasConfig:      true,
			hasClusterData: true,
		},
		"skipEval": {
			inventory:      "testdata/simple/inventory.yml",
			vars:           "testdata/simple/group_vars",
			playbook:       "testdata/simple/playbook.yml",
			skipEval:       true,
			skipClusterInv: false,
			hasConfig:      false,
			hasClusterData: true,
		},
		"skipClusterData": {
			inventory:      "testdata/simple/inventory.yml",
			vars:           "testdata/simple/group_vars",
			playbook:       "testdata/simple/playbook.yml",
			skipEval:       true,
			skipClusterInv: true,
			hasConfig:      false,
			hasClusterData: false,
		},
		"complex": {
			inventory:      "testdata/complex/inventory.yml",
			vars:           "testdata/complex/group_vars",
			playbook:       "testdata/complex/playbook.yml",
			skipEval:       false,
			skipClusterInv: false,
			hasConfig:      true,
			hasClusterData: true,
		},
	}

	ejsonSettings := ejson.Settings{}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			inv, err := inventory.NewInventory(tc.inventory, ejsonSettings, true)
			assert.NilError(t, err)

			targets, err := target.NewTargets(".*", []string{}, tc.vars, inv, true, &ejsonSettings)
			assert.NilError(t, err)
			// create fake clients for each target so we can simulate
			// retrieving the cluster-inventory for each
			for _, tgt := range targets.Targets() {
				clientset := fake.NewSimpleClientset(&v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:        "cluster-inventory",
						Namespace:   tgt.Entry().ConfigNamespace(),
						Annotations: map[string]string{},
					},
					Data: map[string]string{
						"inventory": `{
							"foo": "bar"
						}`,
					},
				})
				tgt.Entry().Kubeconfig().SetClient(clientset)
			}

			playbookSet, err := NewSet(tc.playbook, targets, tc.skipEval, tc.skipClusterInv)
			assert.NilError(t, err)
			assert.Equal(t, len(targets.Targets()), len(playbookSet))
			for name, playbook := range playbookSet {
				t.Run(name, func(t *testing.T) {
					assert.Assert(t, playbook.Raw != nil)
					assert.Equal(t, tc.hasConfig, playbook.Config != nil)
					v := playbook.Raw["vars"]
					vars, ok := v.(map[string]interface{})
					assert.Assert(t, ok)
					_, ok = vars["foo"]
					assert.Equal(t, tc.hasClusterData, ok)
				})
			}
		})
	}
}
