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

package inventory

import (
	"testing"

	"gotest.tools/assert"
)

func TestEmptyConfig(t *testing.T) {
	config := NewConfig()
	assert.Assert(t, config != nil)
	assert.Assert(t, config.Inventory != nil)
	assert.Equal(t, 0, len(config.Inventory))
}

func TestConfig(t *testing.T) {
	data := map[interface{}]interface{}{
		"inventory": []map[interface{}]interface{}{
			{
				"name":             "testentry",
				"groups":           []interface{}{"g0", "g1", "g2"},
				"config_namespace": "kube-system",
				"kubeconfig": map[interface{}]interface{}{
					"backend": "s3",
					"params": map[interface{}]interface{}{
						"param1": "value1",
						"param2": "value2",
						"param3": "value3",
						"param4": "value4",
					},
				},
			},
		},
	}
	config, err := NewConfigFromMap(&data)
	assert.NilError(t, err)
	assert.Assert(t, config != nil)
	assert.Assert(t, config.Inventory != nil)
	assert.Equal(t, 1, len(config.Inventory))
	assert.Equal(t, "testentry", config.Inventory[0].Name)
	assert.Equal(t, "kube-system", config.Inventory[0].ConfigNamesace)
	assert.Assert(t, config.Inventory[0].Kubeconfig != nil)
	assert.Equal(t, "s3", config.Inventory[0].Kubeconfig.Backend)
	assert.Assert(t, config.Inventory[0].Kubeconfig.Params != nil)
	assert.Equal(t, 4, len(*config.Inventory[0].Kubeconfig.Params))
}
