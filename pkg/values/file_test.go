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
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/bedag/kusible/pkg/wrapper/ejson"
	"gotest.tools/assert"
	"sigs.k8s.io/yaml"
)

func TestFile(t *testing.T) {
	loadYaml := func(path string) (map[string]interface{}, error) {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		var result map[string]interface{}
		err = yaml.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	marshalMethods := []string{"JSON", "YAML"}
	tests := map[string]struct {
		input    string
		skipEval bool
		expected string
	}{
		"simple":                {input: "simple.yml", skipEval: false, expected: "simple.expected.yml"},
		"simple-ejson":          {input: "simple.ejson", skipEval: false, expected: "simple.expected.yml"},
		"spruce-eval":           {input: "spruce-eval.yml", skipEval: false, expected: "spruce-eval.expected.yml"},
		"spruce-skip-eval":      {input: "spruce-eval.yml", skipEval: true, expected: "spruce-eval.yml"},
		"spruce-eval-ejson":     {input: "spruce-eval.ejson", skipEval: false, expected: "spruce-eval.expected.yml"},
		"simple-ejson-wrongkey": {input: "simple-wrongkey.ejson", skipEval: false, expected: "simple-wrongkey.ejson"},
		"fully-empty":           {input: "fully-empty.yml", skipEval: false, expected: "empty.yml"},
		"empty-yaml":            {input: "empty.yml", skipEval: false, expected: "empty.yml"},
		"empty-json":            {input: "empty.json", skipEval: false, expected: "empty.yml"},
	}

	ejsonSettings := ejson.Settings{
		KeyDir: "testdata/keydir",
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f, err := NewFile("testdata/file/"+tc.input, tc.skipEval, ejsonSettings)
			assert.NilError(t, err)
			got := f.Map()
			assert.NilError(t, err)
			delete(got, "_public_key")

			want, err := loadYaml("testdata/file/" + tc.expected)
			assert.NilError(t, err)
			delete(want, "_public_key")

			assert.DeepEqual(t, want, got)

			for _, method := range marshalMethods {
				r := reflect.ValueOf(f).MethodByName(method).Call([]reflect.Value{})
				resultBytes := r[0].Bytes()
				// cannot use assert.NilError here because we cannot
				// cast nil to error
				assert.Assert(t, r[1].Interface() == nil)

				err = yaml.Unmarshal(resultBytes, map[string]interface{}{})
				assert.NilError(t, err)
			}
		})
	}
}
