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
	"strings"
	"testing"

	"gotest.tools/assert"
)

func TestPattern(t *testing.T) {
	tests := map[string]struct {
		expr      string
		inGroups  []string
		outGroups []string
		err       bool
		match     bool
	}{
		"no-expr-no-groups": {
			expr:      "",
			inGroups:  []string{},
			outGroups: []string{},
			err:       true,
			match:     false,
		},
		"only-modifier-!": {
			expr:      "!",
			inGroups:  []string{},
			outGroups: []string{},
			err:       true,
			match:     false,
		},
		"broken-regex": {
			expr:      "[",
			inGroups:  []string{},
			outGroups: []string{},
			err:       true,
			match:     false,
		},
		"only-modifier-&": {
			expr:      "&",
			inGroups:  []string{},
			outGroups: []string{},
			err:       true,
			match:     false,
		},
		"single-letter-expression": {
			expr:      "x",
			inGroups:  []string{"x"},
			outGroups: []string{"x"},
			err:       false,
			match:     true,
		},
		"single-match-multi-group": {
			expr:      "x",
			inGroups:  []string{"x", "y"},
			outGroups: []string{"x"},
			err:       false,
			match:     true,
		},
		"regexp-modifier-&-multi-match": {
			expr:      "&prod.*",
			inGroups:  []string{"prod-a", "test-x", "prod-b"},
			outGroups: []string{"prod-a", "prod-b"},
			err:       false,
			match:     true,
		},
		"regexp-modifier-!-multi-match": {
			expr:      "!test.*",
			inGroups:  []string{"prod-a", "test-x", "prod-b"},
			outGroups: []string{"test-x"},
			err:       false,
			match:     false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pattern, err := NewPattern(tc.expr, tc.inGroups)
			assert.Equal(t, tc.err, err != nil)
			if pattern != nil {
				assert.Equal(t, tc.match, pattern.Matches())
				assert.DeepEqual(t, tc.outGroups, pattern.Groups())
			}
		})
	}
}

func TestValidator(t *testing.T) {
	type input struct {
		groups []string
		valid  bool
	}
	tests := map[string]struct {
		expressions []string
		inputs      []*input
	}{
		"empty": {
			expressions: []string{},
			inputs: []*input{
				&input{
					groups: []string{},
					valid:  false,
				},
			},
		},
		"single-expression": {
			expressions: []string{"a"},
			inputs: []*input{
				&input{
					groups: []string{"a"},
					valid:  true,
				},
				&input{
					groups: []string{"b"},
					valid:  false,
				},
				&input{
					groups: []string{"a", "b"},
					valid:  true,
				},
				&input{
					groups: []string{"x", "y"},
					valid:  false,
				},
				&input{
					groups: []string{},
					valid:  false,
				},
				&input{
					groups: []string{"", ""},
					valid:  false,
				},
			},
		},
		"multi-expression": {
			expressions: []string{"a", "b"},
			inputs: []*input{
				&input{
					groups: []string{"a"},
					valid:  true,
				},
				&input{
					groups: []string{"b"},
					valid:  true,
				},
				&input{
					groups: []string{"a", "b"},
					valid:  true,
				},
				&input{
					groups: []string{"b", "a"},
					valid:  true,
				},
				&input{
					groups: []string{"a", "b", "c"},
					valid:  true,
				},
				&input{
					groups: []string{"ab", "ba"},
					valid:  false,
				},
			},
		},
		"&-modifier": {
			expressions: []string{"a", "&b"},
			inputs: []*input{
				&input{
					groups: []string{"a"},
					valid:  false,
				},
				&input{
					groups: []string{"b"},
					valid:  false,
				},
				&input{
					groups: []string{"a", "b"},
					valid:  true,
				},
				&input{
					groups: []string{"b", "a"},
					valid:  true,
				},
				&input{
					groups: []string{"a", "b", "c"},
					valid:  true,
				},
			},
		},
		"!-modifier": {
			expressions: []string{"a", "!b"},
			inputs: []*input{
				&input{
					groups: []string{"a"},
					valid:  true,
				},
				&input{
					groups: []string{"b"},
					valid:  false,
				},
				&input{
					groups: []string{"a", "b"},
					valid:  false,
				},
				&input{
					groups: []string{"b", "a"},
					valid:  false,
				},
				&input{
					groups: []string{"a", "b", "c"},
					valid:  false,
				},
			},
		},
		"solo-&-modifier": {
			expressions: []string{"&b"},
			inputs: []*input{
				&input{
					groups: []string{"a"},
					valid:  false,
				},
				&input{
					groups: []string{"b"},
					valid:  true,
				},
				&input{
					groups: []string{"a", "b"},
					valid:  true,
				},
				&input{
					groups: []string{"a", "b", "c"},
					valid:  true,
				},
			},
		},
		"solo-!-modifier": {
			expressions: []string{"!b"},
			inputs: []*input{
				&input{
					groups: []string{"a"},
					valid:  true,
				},
				&input{
					groups: []string{"b"},
					valid:  false,
				},
				&input{
					groups: []string{"a", "b"},
					valid:  false,
				},
				&input{
					groups: []string{"a", "b", "c"},
					valid:  false,
				},
				&input{
					groups: []string{"a", "c", "d"},
					valid:  true,
				},
			},
		},
		"multi-&-modifier": {
			expressions: []string{"a", "&b", "&c"},
			inputs: []*input{
				&input{
					groups: []string{"a"},
					valid:  false,
				},
				&input{
					groups: []string{"b"},
					valid:  false,
				},
				&input{
					groups: []string{"c"},
					valid:  false,
				},
				&input{
					groups: []string{"a", "b"},
					valid:  false,
				},
				&input{
					groups: []string{"a", "c"},
					valid:  false,
				},
				&input{
					groups: []string{"b", "c"},
					valid:  false,
				},
				&input{
					groups: []string{"a", "b", "c"},
					valid:  true,
				},
			},
		},
		"multi-!-modifier": {
			expressions: []string{"a", "!b", "!c"},
			inputs: []*input{
				&input{
					groups: []string{"a"},
					valid:  true,
				},
				&input{
					groups: []string{"b"},
					valid:  false,
				},
				&input{
					groups: []string{"c"},
					valid:  false,
				},
				&input{
					groups: []string{"a", "b"},
					valid:  false,
				},
				&input{
					groups: []string{"a", "c"},
					valid:  false,
				},
				&input{
					groups: []string{"b", "c"},
					valid:  false,
				},
				&input{
					groups: []string{"a", "b", "c"},
					valid:  false,
				},
			},
		},
		"mixed-with-regex": {
			expressions: []string{"a", "b-.*", "&c.*", "!c-.*"},
			inputs: []*input{
				&input{
					groups: []string{"a", "c"},
					valid:  true,
				},
				&input{
					groups: []string{"b-b", "c"},
					valid:  true,
				},
				&input{
					groups: []string{"b", "c"},
					valid:  false,
				},
				&input{
					groups: []string{"a", "b-b", "c"},
					valid:  true,
				},
				&input{
					groups: []string{"a", "b-b", "c", "c-c"},
					valid:  false,
				},
				&input{
					groups: []string{"a", "b-b", "c-c"},
					valid:  false,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			for _, i := range tc.inputs {
				id := fmt.Sprintf("%s|%s", strings.Join(tc.expressions, ","), strings.Join(i.groups, ","))
				t.Run(id, func(t *testing.T) {
					v := Validator{}
					for _, expr := range tc.expressions {
						pattern, err := NewPattern(expr, i.groups)
						assert.NilError(t, err)
						v.Add(pattern)
					}
					assert.Equal(t, i.valid, v.Valid())
				})
			}
		})
	}
}
