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
package inventory

import (
	"testing"

	"github.com/bedag/kusible/pkg/inventory/config"
	"gotest.tools/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestEntryMatchLimits(t *testing.T) {
	limitsMatching := []string{"dev", "test"}
	limitsNotMatching := []string{"dev", "foo"}
	entry := &Entry{
		name:   "test",
		groups: []string{"dev", "test", "state", "prod"},
	}

	result, err := entry.MatchLimits(limitsMatching)
	assert.NilError(t, err)
	assert.Assert(t, result)

	result, err = entry.MatchLimits(limitsNotMatching)
	assert.NilError(t, err)
	assert.Assert(t, !result)
}

func TestEntryValidGroups(t *testing.T) {
	limits := []string{"dev", "test"}
	entry := &Entry{
		name:   "test",
		groups: []string{"dev", "test", "state", "prod"},
	}

	result, err := entry.ValidGroups(limits)
	assert.NilError(t, err)
	assert.DeepEqual(t, limits, result)
}

func TestClusterInventory(t *testing.T) {
	clusterInventoryConfig := &config.ClusterInventory{
		Namespace: "kube-system",
		ConfigMap: "cluster-inventory",
	}

	tests := map[string]struct {
		configmap   *v1.ConfigMap
		errExpected bool
	}{
		"working empty": {
			configmap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:        clusterInventoryConfig.ConfigMap,
					Namespace:   clusterInventoryConfig.Namespace,
					Annotations: map[string]string{},
				},
				Data: map[string]string{
					"inventory": "{}",
				},
			},
			errExpected: false,
		},
		"missing inventory": {
			configmap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:        clusterInventoryConfig.ConfigMap,
					Namespace:   clusterInventoryConfig.Namespace,
					Annotations: map[string]string{},
				},
				Data: map[string]string{},
			},
			errExpected: true,
		},
		"inventory unparsable": {
			configmap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:        clusterInventoryConfig.ConfigMap,
					Namespace:   clusterInventoryConfig.Namespace,
					Annotations: map[string]string{},
				},
				Data: map[string]string{
					"inventory": "{",
				},
			},
			errExpected: true,
		},
		"wrong namespace": {
			configmap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:        clusterInventoryConfig.ConfigMap,
					Namespace:   "wrong",
					Annotations: map[string]string{},
				},
				Data: map[string]string{
					"inventory": "{}",
				},
			},
			errExpected: true,
		},
		"wrong ConfigMap": {
			configmap: &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "some-inventory",
					Namespace:   clusterInventoryConfig.Namespace,
					Annotations: map[string]string{},
				},
				Data: map[string]string{
					"inventory": "{}",
				},
			},
			errExpected: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			clientset := fake.NewSimpleClientset(tc.configmap)
			assert.Assert(t, clientset != nil)

			k := &Kubeconfig{
				loader: nil,
				config: nil,
				client: clientset,
			}

			entry := &Entry{
				name:                   "test",
				groups:                 []string{"test"},
				clusterInventoryConfig: clusterInventoryConfig,
				kubeconfig:             k,
			}
			_, err := entry.ClusterInventory()
			assert.Equal(t, tc.errExpected, err != nil)
		})
	}
}
