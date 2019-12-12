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
	"gotest.tools/assert"
)

func TestData(t *testing.T) {
	loadRaw := func(path string) (map[interface{}]interface{}, error) {
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
		return raw, nil
	}

	tests := map[string]struct {
		input  string
		method string
	}{
		"yaml": {input: "simple.yml", method: "YAML"},
		"json": {input: "simple.yml", method: "JSON"},
	}

	// try to parse output of the different data methods
	// as yaml to test if the results are valid
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var data data
			data, err := loadRaw("testdata/file/" + tc.input)
			assert.NilError(t, err)
			r := reflect.ValueOf(&data).MethodByName(tc.method).Call([]reflect.Value{})
			resultBytes := r[0].Bytes()
			// cannot use assert.NilError here because we cannot
			// cast nil to error
			assert.Assert(t, r[1].Interface() == nil)

			_, err = simpleyaml.NewYaml(resultBytes)
			assert.NilError(t, err)
		})
	}
}
