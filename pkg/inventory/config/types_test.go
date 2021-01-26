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
	"testing"

	"gotest.tools/assert"
	"sigs.k8s.io/yaml"
	//"k8s.io/apimachinery/pkg/util/yaml"
)

func TestEmptyConfig(t *testing.T) {
	config := NewConfig()
	assert.Assert(t, config != nil)
	assert.Assert(t, config.Inventory != nil)
	assert.Equal(t, 0, len(config.Inventory))
}

func TestConfig(t *testing.T) {
	data := []byte(`---
inventory:
  - name: "testentry"
    groups: ["g0","g1","g2"]
    cluster_inventory:
      namespace: "kube-system"
      configmap: "cluster-inventory"
    kubeconfig:
      backend: "s3"
      params:
        param1: "value1"
        param2: "value2"
        param3: "value3"
        param4: "value4"
        path: "some/path"
`)

	var expectedMap map[string]interface{}
	err := yaml.Unmarshal(data, &expectedMap)
	assert.NilError(t, err)

	// ensure that the config gets parsed correctly
	config, err := NewConfigFromMap(&expectedMap)
	assert.NilError(t, err)
	assert.Assert(t, config != nil)
	assert.Assert(t, config.Inventory != nil)
	assert.Equal(t, 1, len(config.Inventory))
	assert.Equal(t, "testentry", config.Inventory[0].Name)
	assert.Equal(t, "kube-system", config.Inventory[0].ClusterInventory.Namespace)
	assert.Equal(t, "cluster-inventory", config.Inventory[0].ClusterInventory.ConfigMap)
	assert.Equal(t, "s3", config.Inventory[0].Kubeconfig.Backend)
	assert.Assert(t, config.Inventory[0].Kubeconfig.Params != nil)
	assert.Equal(t, 5, len(config.Inventory[0].Kubeconfig.Params))
	assert.Equal(t, "some/path", config.Inventory[0].Kubeconfig.Params["path"])

	// ensure that converting the parsed config back to yaml
	// results in the same yaml that was used to create the config
	var resultMap map[string]interface{}
	configYaml, err := yaml.Marshal(config)
	assert.NilError(t, err)
	err = yaml.Unmarshal(configYaml, &resultMap)
	assert.NilError(t, err)
	assert.DeepEqual(t, expectedMap, resultMap)
}

func TestEmptyEntryt(t *testing.T) {
	data := []byte(`---
inventory:
  - name: "testentry"
`)

	var expectedMap map[string]interface{}
	err := yaml.Unmarshal(data, &expectedMap)
	assert.NilError(t, err)

	// ensure that the config gets parsed correctly
	config, err := NewConfigFromMap(&expectedMap)
	assert.NilError(t, err)
	assert.Assert(t, config != nil)
	assert.Assert(t, config.Inventory != nil)
	assert.Equal(t, 1, len(config.Inventory))
	assert.Equal(t, "testentry", config.Inventory[0].Name)
	assert.Equal(t, "s3", config.Inventory[0].Kubeconfig.Backend)
	assert.Assert(t, config.Inventory[0].Kubeconfig.Params != nil)
	assert.Equal(t, "testentry/kubeconfig/kubeconfig.enc.7z", config.Inventory[0].Kubeconfig.Params["path"])
}
