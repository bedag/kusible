/*
Copyright © 2019 Copyright © 2021 Bedag Informatik AG

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
	"fmt"
	"sort"
	"testing"

	"github.com/bedag/kusible/pkg/inventory/config"
	"github.com/bedag/kusible/pkg/loader"
	"github.com/bedag/kusible/pkg/wrapper/ejson"
	"github.com/go-test/deep"
	"gotest.tools/assert"
)

func basicInventoryTest(path string, filter string, limits []string, skip bool, clusterInvConfig config.ClusterInventory, expected []string) (*Inventory, error) {
	ejsonSettings := ejson.Settings{
		PrivKey:     "",
		KeyDir:      "",
		SkipDecrypt: false,
	}

	inventory, err := NewInventory(path, ejsonSettings, skip, clusterInvConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory: %s", err)
	}

	result, err := inventory.EntryNames(filter, limits)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve inventory entries: %s", err)
	}

	// we only want to compare the elements and not
	// the order of the elements
	sort.Strings(expected)
	sort.Strings(result)

	if diff := deep.Equal(result, expected); diff != nil {
		return nil, fmt.Errorf("unexpected list of inventory entries: %#v", diff)
	}
	return inventory, nil
}

func TestInventoryBare(t *testing.T) {
	inventoryPath := "testdata/clusters_bare.yaml"
	skipKubeconfig := true
	filter := ".*"
	limits := []string{}
	clusterInventory := config.ClusterInventory{}
	expected := []string{
		"test",
	}

	inventory, err := basicInventoryTest(inventoryPath, filter, limits, skipKubeconfig, clusterInventory, expected)
	assert.NilError(t, err)
	entry := inventory.entries["test"]
	assert.Assert(t, entry.kubeconfig != nil)
	assert.Assert(t, entry.kubeconfig.loader != nil)
	assert.Equal(t, "s3", entry.kubeconfig.loader.Type())
	assert.Equal(t, "all", entry.groups[0])
	assert.Equal(t, "test", entry.groups[1])
	assert.Equal(t, "kube-system", entry.ClusterInventoryConfig().Namespace)
	assert.Equal(t, "cluster-inventory", entry.ClusterInventoryConfig().ConfigMap)
	expectedPath := fmt.Sprintf("%s/%s", expected[0], "kubeconfig/kubeconfig.enc.7z")
	backendConfig := entry.kubeconfig.loader.Config().(*loader.S3Config)
	assert.Equal(t, expectedPath, backendConfig.Path)
	assert.Assert(t, backendConfig.Region != "")
}

func TestClusterInventoryConfigDefaults(t *testing.T) {
	inventoryPath := "testdata/clusters_bare.yaml"
	skipKubeconfig := true
	filter := ".*"
	limits := []string{}
	expected := []string{
		"test",
	}

	tests := map[string]struct {
		defaults      config.ClusterInventory
		wantNamespace string
		wantConfigMap string
	}{
		"no defaults": {
			defaults:      config.ClusterInventory{},
			wantNamespace: "kube-system",
			wantConfigMap: "cluster-inventory",
		},
		"namespace and configmap default": {
			defaults: config.ClusterInventory{
				Namespace: "some-namespace",
				ConfigMap: "some-configmap",
			},
			wantNamespace: "some-namespace",
			wantConfigMap: "some-configmap",
		},
		"namespace default": {
			defaults: config.ClusterInventory{
				Namespace: "some-namespace",
			},
			wantNamespace: "some-namespace",
			wantConfigMap: "cluster-inventory",
		},
		"configmap default": {
			defaults: config.ClusterInventory{
				ConfigMap: "some-configmap",
			},
			wantNamespace: "kube-system",
			wantConfigMap: "some-configmap",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			inventory, err := basicInventoryTest(inventoryPath, filter, limits, skipKubeconfig, tc.defaults, expected)
			assert.NilError(t, err)
			entry := inventory.entries["test"]
			assert.Equal(t, tc.wantNamespace, entry.ClusterInventoryConfig().Namespace)
			assert.Equal(t, tc.wantConfigMap, entry.ClusterInventoryConfig().ConfigMap)
		})
	}

}

func TestInventoryEntriesFull(t *testing.T) {
	inventoryPath := "testdata/clusters_default.yaml"
	skipKubeconfig := true
	filter := ".*"
	limits := []string{}
	clusterInventory := config.ClusterInventory{}
	expected := []string{
		"cluster-test-01-preflight",
		"cluster-dev-01",
		"cluster-test-01",
		"cluster-stage-01",
		"cluster-stage-02",
		"cluster-stage-03",
		"cluster-prod-01",
		"cluster-prod-02",
		"cluster-prod-03",
		"cluster-prod-04",
	}

	_, err := basicInventoryTest(inventoryPath, filter, limits, skipKubeconfig, clusterInventory, expected)
	assert.NilError(t, err)
}

func TestInventoryEntriesSingle(t *testing.T) {
	inventoryPath := "testdata/clusters_default.yaml"
	skipKubeconfig := true
	expected := []string{
		"cluster-dev-01",
	}
	limits := []string{}
	clusterInventory := config.ClusterInventory{}
	filter := expected[0]

	_, err := basicInventoryTest(inventoryPath, filter, limits, skipKubeconfig, clusterInventory, expected)
	assert.NilError(t, err)
}

func TestInventoryEntriesLimits(t *testing.T) {
	inventoryPath := "testdata/clusters_default.yaml"
	skipKubeconfig := true
	expected := []string{
		"cluster-stage-01",
		"cluster-stage-02",
		"cluster-stage-03",
	}
	limits := []string{
		"stage",
	}
	clusterInventory := config.ClusterInventory{}
	filter := ".*"

	_, err := basicInventoryTest(inventoryPath, filter, limits, skipKubeconfig, clusterInventory, expected)
	assert.NilError(t, err)
}

func TestInventoryLoader(t *testing.T) {
	inventoryPath := "testdata/clusters_file.yaml"
	skipKubeconfig := false
	filter := ".*"
	limits := []string{}
	clusterInventory := config.ClusterInventory{}
	expected := []string{
		"cluster-test-01",
		"cluster-test-02",
		"cluster-test-03",
	}
	inventory, err := basicInventoryTest(inventoryPath, filter, limits, skipKubeconfig, clusterInventory, expected)
	assert.NilError(t, err)
	for _, entry := range inventory.entries {
		ldr := entry.kubeconfig.loader
		assert.Assert(t, ldr != nil)
		assert.Equal(t, "file", ldr.Type())
	}
}

func TestInventoryEntryGroups(t *testing.T) {
	inventoryPath := "testdata/clusters_file.yaml"
	skipKubeconfig := false
	filter := ".*"
	limits := []string{}
	clusterInventory := config.ClusterInventory{}
	expected := []string{
		"cluster-test-01",
		"cluster-test-02",
		"cluster-test-03",
	}
	inventory, err := basicInventoryTest(inventoryPath, filter, limits, skipKubeconfig, clusterInventory, expected)
	assert.NilError(t, err)
	for _, entry := range inventory.entries {
		name := entry.name
		groups := entry.groups
		assert.Equal(t, "all", groups[0])
		assert.Equal(t, name, groups[len(groups)-1])
	}
}
