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

	"github.com/geofffranks/simpleyaml"
	"github.com/mitchellh/mapstructure"
	"gotest.tools/assert"
)

func TestFile(t *testing.T) {
	loadYaml := func(path string) (map[string]interface{}, error) {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		yamlData, err := simpleyaml.NewYaml(data)
		if err != nil {
			return nil, err
		}

		raw, err := yamlData.Map()
		if err != nil {
			return nil, err
		}

		var result map[string]interface{}
		err = mapstructure.Decode(raw, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	marshalMethods := []string{"JSON", "YAML"}
	tests := map[string]struct {
		input    string
		expected string
	}{
		"simple":                {input: "simple.yml", expected: "simple.expected.yml"},
		"simple-ejson":          {input: "simple.ejson", expected: "simple.expected.yml"},
		"spruce-eval":           {input: "spruce-eval.yml", expected: "spruce-eval.expected.yml"},
		"spruce-eval-ejson":     {input: "spruce-eval.ejson", expected: "spruce-eval.expected.yml"},
		"simple-ejson-wrongkey": {input: "simple-wrongkey.ejson", expected: "simple-wrongkey.ejson"},
		// TODO: Fix handling of empty yaml files
		//"fully-empty":           {input: "fully-empty.yml", expected: "empty.yml"},
		//"empty-yaml":            {input: "empty.yml", expected: "empty.yml"},
		//"empty-json":            {input: "empty.json", expected: "empty.yml"},
	}

	ejsonSettings := EjsonSettings{
		KeyDir: "testdata/keydir",
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f, err := NewFile("testdata/file/"+tc.input, false, ejsonSettings)
			assert.NilError(t, err)
			got, err := f.Map()
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

				_, err = simpleyaml.NewYaml(resultBytes)
				assert.NilError(t, err)
			}
		})
	}
}
