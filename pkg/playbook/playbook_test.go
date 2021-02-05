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

	"github.com/bedag/kusible/pkg/inventory"
	invconfig "github.com/bedag/kusible/pkg/inventory/config"
	"github.com/bedag/kusible/pkg/target"
	"github.com/bedag/kusible/pkg/wrapper/ejson"
	"gotest.tools/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestConvert(t *testing.T) {
	invPath := "testdata/simple/inventory.yml"
	varsPath := "testdata/simple/group_vars"
	playbookPath := "testdata/simple/playbook.yml"

	tests := map[string]struct {
		skipEval       bool
		skipClusterInv bool
		hasConfig      bool
		hasClusterData bool
	}{
		"full": {
			skipEval:       false,
			skipClusterInv: false,
			hasConfig:      true,
			hasClusterData: true,
		},
		"skipEval": {
			skipEval:       true,
			skipClusterInv: false,
			hasConfig:      false,
			hasClusterData: true,
		},
		"skipClusterData": {
			skipEval:       true,
			skipClusterInv: true,
			hasConfig:      false,
			hasClusterData: false,
		},
	}

	ejsonSettings := ejson.Settings{}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			inv, err := inventory.NewInventory(invPath, ejsonSettings, true, invconfig.ClusterInventory{})
			assert.NilError(t, err)

			targets, err := target.NewTargets(".*", []string{}, varsPath, inv, true, &ejsonSettings)
			assert.NilError(t, err)
			// create fake clients for each target so we can simulate
			// retrieving the cluster-inventory for each
			for _, tgt := range targets.Targets() {
				clientset := fake.NewSimpleClientset(&v1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:        tgt.Entry().ClusterInventoryConfig().ConfigMap,
						Namespace:   tgt.Entry().ClusterInventoryConfig().Namespace,
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

			playbookSet, err := NewSet(playbookPath, targets, tc.skipEval, tc.skipClusterInv)
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

					// json/yaml with raw == true should always return data.
					// With raw == false we expect data if a config is present,
					// and if not, we expect an empty result (zero length byte
					// array)
					t.Run("json-raw", func(t *testing.T) {
						result, err := playbook.JSON(true)
						assert.NilError(t, err)
						assert.Assert(t, len(result) > 0)
					})
					t.Run("json", func(t *testing.T) {
						result, err := playbook.JSON(false)
						assert.NilError(t, err)
						assert.Equal(t, tc.hasConfig, len(result) > 0)
					})
					t.Run("yaml-raw", func(t *testing.T) {
						result, err := playbook.YAML(true)
						assert.NilError(t, err)
						assert.Assert(t, len(result) > 0)
					})
					t.Run("yaml", func(t *testing.T) {
						result, err := playbook.YAML(false)
						assert.NilError(t, err)
						assert.Equal(t, tc.hasConfig, len(result) > 0)
					})
					t.Run("map-raw", func(t *testing.T) {
						result, err := playbook.Map(true)
						assert.NilError(t, err)
						assert.Assert(t, len(result) > 0)
					})
					t.Run("map", func(t *testing.T) {
						result, err := playbook.Map(false)
						assert.NilError(t, err)
						assert.Equal(t, tc.hasConfig, len(result) > 0)
					})
				})
			}
		})
	}
}
