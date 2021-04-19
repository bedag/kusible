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

package target

import (
	"sort"
	"testing"

	"github.com/bedag/kusible/pkg/inventory"
	invconf "github.com/bedag/kusible/pkg/inventory/config"
	"github.com/bedag/kusible/pkg/wrapper/ejson"
	"gotest.tools/assert"
)

func TestTargets(t *testing.T) {
	ejsonSettings := ejson.Settings{}

	inv, err := inventory.NewInventory("testdata/inventory.yml", ejsonSettings, true, invconf.ClusterInventory{})
	assert.NilError(t, err)

	type expected struct {
		entries []string
		values  map[string]interface{}
		error   bool
	}
	tests := map[string]struct {
		filter   string
		limits   []string
		skipEval bool
		expected expected
	}{
		"all": {
			filter:   ".*",
			limits:   []string{},
			skipEval: false,
			expected: expected{
				entries: []string{"cluster-01", "cluster-02"},
				values: map[string]interface{}{
					"cluster-01": map[string]interface{}{
						"key1": "file-01",
						"key2": "file-01",
						"key3": "file-01",
						"eval": "file-01",
					},
					"cluster-02": map[string]interface{}{
						"key1": "file-03",
						"key2": "file-01",
						"key3": "file-01",
						"eval": "file-03",
					},
				},
				error: false,
			},
		},
		"skip-eval": {
			filter:   ".*",
			limits:   []string{},
			skipEval: true,
			expected: expected{
				entries: []string{"cluster-01", "cluster-02"},
				values: map[string]interface{}{
					"cluster-01": map[string]interface{}{
						"key1": "file-01",
						"key2": "file-01",
						"key3": "file-01",
						"eval": "(( grab key1 ))",
					},
					"cluster-02": map[string]interface{}{
						"key1": "file-03",
						"key2": "file-01",
						"key3": "file-01",
						"eval": "(( grab key1 ))",
					},
				},
				error: false,
			},
		},
		"none(empty)": {
			filter:   "",
			limits:   []string{},
			skipEval: false,
			expected: expected{
				entries: []string{},
				values:  map[string]interface{}{},
				error:   false,
			},
		},
		"none(unknown)": {
			filter:   "unknown",
			limits:   []string{},
			skipEval: false,
			expected: expected{
				entries: []string{},
				values:  map[string]interface{}{},
				error:   false,
			},
		},
		"limited": {
			filter:   ".*",
			limits:   []string{"group-03"},
			skipEval: false,
			expected: expected{
				entries: []string{"cluster-02"},
				values: map[string]interface{}{
					"cluster-02": map[string]interface{}{
						"key1": "file-03",
						"key2": "file-01",
						"key3": "file-01",
						"eval": "file-03",
					},
				},
				error: false,
			},
		},
		"filtered": {
			filter:   ".*-01",
			limits:   []string{},
			skipEval: false,
			expected: expected{
				entries: []string{"cluster-01"},
				values: map[string]interface{}{
					"cluster-01": map[string]interface{}{
						"key1": "file-01",
						"key2": "file-01",
						"key3": "file-01",
						"eval": "file-01",
					},
				},
				error: false,
			},
		},
		"contradicting": {
			filter:   ".*-01",
			limits:   []string{"group-03"},
			skipEval: false,
			expected: expected{
				entries: []string{},
				values:  map[string]interface{}{},
				error:   false,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			targets, err := NewTargets(tc.filter, tc.limits, "testdata/group_vars", inv, tc.skipEval, &ejsonSettings)
			assert.Equal(t, tc.expected.error, err != nil)
			if !tc.expected.error {
				gotTargets := targets.Targets()
				assert.Equal(t, len(tc.expected.entries), len(gotTargets))
				wantNames := tc.expected.entries
				gotNames := targets.Names()
				sort.Strings(wantNames)
				sort.Strings(gotNames)
				assert.DeepEqual(t, wantNames, gotNames)
				for name, gotTarget := range gotTargets {
					wantValues := tc.expected.values[name]
					gotValues := gotTarget.Values().Map()
					assert.DeepEqual(t, wantValues, gotValues)
				}
			}
		})
	}
}
