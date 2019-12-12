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

package loader

import (
	"testing"

	"gotest.tools/assert"
)

func TestLoader(t *testing.T) {
	tests := map[string]struct {
		backend     string
		errExpected bool
	}{
		"file":    {backend: "file", errExpected: false},
		"s3":      {backend: "s3", errExpected: false},
		"unknown": {backend: "unknown", errExpected: true},
		"empty":   {backend: "", errExpected: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ldr, err := New(tc.backend, map[string]interface{}{})
			assert.Equal(t, tc.errExpected, err != nil)
			if !tc.errExpected {
				assert.Equal(t, tc.backend, ldr.Type())
			}
		})
	}
}
