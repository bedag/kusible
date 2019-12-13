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

package target

import (
	"testing"

	"github.com/bedag/kusible/pkg/inventory"
	invconf "github.com/bedag/kusible/pkg/inventory/config"
	"github.com/bedag/kusible/pkg/values"
	"gotest.tools/assert"
)

func TestTarget(t *testing.T) {
	config := &invconf.Entry{
		Name:   "cluster-01",
		Groups: []string{"group-01", "group-02"},
		Kubeconfig: &invconf.Kubeconfig{
			Params: new(invconf.Params),
		},
	}
	wantValues := map[string]interface{}{
		"key1": "file-02",
		"key2": "file-02",
		"key3": "file-01",
	}
	entry, err := inventory.NewEntryFromConfig(config)
	assert.NilError(t, err)
	target, err := New(entry, "testdata/group_vars", &values.EjsonSettings{})
	assert.NilError(t, err)
	gotValues := target.Values().Map()
	assert.DeepEqual(t, wantValues, gotValues)
}
