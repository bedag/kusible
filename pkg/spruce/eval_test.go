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

package spruce

import (
	"testing"

	"gotest.tools/assert"
)

func TestUtilSpruceEval(t *testing.T) {
	tests := map[string]struct {
		data     map[string]interface{}
		skip     bool
		prune    []string
		err      bool
		expected map[string]interface{}
	}{
		"simple-eval": {
			data:     map[string]interface{}{"key1": "test", "key2": "(( grab key1 ))"},
			skip:     false,
			prune:    []string{},
			err:      false,
			expected: map[string]interface{}{"key1": "test", "key2": "test"},
		},
		"skip-eval": {
			data:     map[string]interface{}{"key1": "test", "key2": "(( grab key1 ))"},
			skip:     true,
			prune:    []string{},
			err:      false,
			expected: map[string]interface{}{"key1": "test", "key2": "(( grab key1 ))"},
		},
		"prune-key": {
			data:     map[string]interface{}{"key1": "test", "key2": "(( grab key1 ))"},
			skip:     false,
			prune:    []string{"key1"},
			err:      false,
			expected: map[string]interface{}{"key2": "test"},
		},
		"error": {
			data:     map[string]interface{}{"key1": "test", "key2": "(( grab key3 ))"},
			skip:     false,
			prune:    []string{},
			err:      true,
			expected: map[string]interface{}{"key1": "test", "key2": "(( grab key3 ))"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := Eval(&tc.data, tc.skip, tc.prune)
			assert.Equal(t, tc.err, err != nil)
			assert.DeepEqual(t, tc.expected, tc.data)
		})
	}
}
