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

package config

import (
	"fmt"
	"sort"
	"testing"

	"gotest.tools/assert"
)

func TestBaseConfig(t *testing.T) {
	config, err := NewBaseConfigFromFile("testdata/playbook.yml")
	assert.NilError(t, err)

	tests := map[int]struct {
		name   string
		groups []string
	}{
		0: {
			name:   "test01",
			groups: []string{"test01", "enabled"},
		},
		1: {
			name:   "test02",
			groups: []string{"test02"},
		},
		2: {
			name:   "test03",
			groups: []string{"test03", "enabled"},
		},
		3: {
			name:   "regexp-test",
			groups: []string{"prod-.*"},
		},
	}

	assert.Equal(t, len(tests), len(config.Plays))

	for id, tc := range tests {
		t.Run(fmt.Sprintf("play-%d", id), func(t *testing.T) {
			assert.Equal(t, tc.name, config.Plays[id].Name)
			assert.DeepEqual(t, tc.groups, config.Plays[id].Groups)
		})
	}
}

func TestBaseConfigFilter(t *testing.T) {
	config, err := NewBaseConfigFromFile("testdata/playbook.yml")
	assert.NilError(t, err)

	tests := map[string]struct {
		groups []string // input groups to filter for;
		// only plays having all of these should be returned
		plays int      // number of expected resulting plays
		names []string // names of expected resulting plays
	}{
		"empty": {
			groups: []string{},
			plays:  4,
			names:  []string{"test01", "test02", "test03", "regexp-test"},
		},
		"only-enabled": {
			groups: []string{"enabled"},
			plays:  2,
			names:  []string{"test01", "test03"},
		},
		"inclusive": {
			groups: []string{"test02", "enabled"},
			plays:  3,
			names:  []string{"test01", "test02", "test03"},
		},
		"notfound": {
			groups: []string{"notexisting"},
			plays:  0,
			names:  []string{},
		},
		"regex": {
			groups: []string{"prod-xx"},
			plays:  1,
			names:  []string{"regexp-test"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := config.Applicable(tc.groups)
			assert.NilError(t, err)
			assert.Equal(t, tc.plays, len(result.Plays))
			playNames := []string{}
			for _, play := range result.Plays {
				playNames = append(playNames, play.Name)
			}
			sort.Strings(playNames)
			sort.Strings(tc.names)
			assert.DeepEqual(t, tc.names, playNames)
		})
	}
}
