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

package values

import (
	"sort"
	"testing"

	"gotest.tools/assert"
)

func TestUtilDataFiles(t *testing.T) {
	files, ok := DirectoryDataFiles("testdata/util", "test-*")
	assert.Assert(t, ok)
	assert.Equal(t, 8, len(files))
	expected := []string{
		"testdata/util/test-a.ejson",
		"testdata/util/test-a.json",
		"testdata/util/test-a.yaml",
		"testdata/util/test-a.yml",
		"testdata/util/test-b.ejson",
		"testdata/util/test-b.json",
		"testdata/util/test-b.yaml",
		"testdata/util/test-b.yml",
	}
	sort.Strings(expected)
	sort.Strings(files)
	assert.DeepEqual(t, expected, files)
}

func TestUtilSpruceEval(t *testing.T) {
	tests := map[string]struct {
		data     map[string]interface{}
		skip     bool
		prune    []string
		expected map[string]interface{}
	}{
		"simple-eval": {
			data:     map[string]interface{}{"key1": "test", "key2": "(( grab key1 ))"},
			skip:     false,
			prune:    []string{},
			expected: map[string]interface{}{"key1": "test", "key2": "test"},
		},
		"skip-eval": {
			data:     map[string]interface{}{"key1": "test", "key2": "(( grab key1 ))"},
			skip:     true,
			prune:    []string{},
			expected: map[string]interface{}{"key1": "test", "key2": "(( grab key1 ))"},
		},
		"prune-key": {
			data:     map[string]interface{}{"key1": "test", "key2": "(( grab key1 ))"},
			skip:     false,
			prune:    []string{"key1"},
			expected: map[string]interface{}{"key2": "test"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := SpruceEval(&tc.data, tc.skip, tc.prune)
			assert.NilError(t, err)
			assert.DeepEqual(t, tc.expected, tc.data)
		})
	}
}
