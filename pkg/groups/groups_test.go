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

package groups

import (
	"sort"
	"testing"

	"gotest.tools/assert"
)

func TestGroups(t *testing.T) {
	tests := map[string]struct {
		filter   string
		limits   []string
		expected []string
	}{
		"all":               {filter: ".*", limits: []string{}, expected: []string{"group01", "group02", "group03", "group04", "group05"}},
		"filter":            {filter: ".*[23]", limits: []string{}, expected: []string{"group02", "group03"}},
		"limits":            {filter: ".*", limits: []string{".*[12]", ".*[45]"}, expected: []string{"group01", "group02", "group04", "group05"}},
		"empty":             {filter: "", limits: []string{}, expected: []string{}},
		"non-group(root)":   {filter: "test", limits: []string{}, expected: []string{}},
		"non-group(subdir)": {filter: "file", limits: []string{}, expected: []string{}},
		"unknown":           {filter: "unknown", limits: []string{}, expected: []string{}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotGroups, err := SortedGroups("testdata", tc.filter, tc.limits)
			assert.NilError(t, err)
			wantGroups := tc.expected
			sort.Strings(wantGroups)
			assert.DeepEqual(t, wantGroups, gotGroups)
		})
	}
}

func TestLimitGroups(t *testing.T) {
	tests := map[string]struct {
		groups   []string
		limits   []string
		expected []string
	}{
		"no-limit":                {groups: []string{"a", "b", "c"}, limits: []string{}, expected: []string{"a", "b", "c"}},
		"empty-limit":             {groups: []string{"a", "b", "c"}, limits: []string{""}, expected: []string{}},
		"match-all":               {groups: []string{"a", "b", "c"}, limits: []string{".*"}, expected: []string{"a", "b", "c"}},
		"match-explicit":          {groups: []string{"a", "aa", "aba", "bab"}, limits: []string{"a"}, expected: []string{"a"}},
		"match-pattern":           {groups: []string{"a", "aa", "aba", "bab"}, limits: []string{".a."}, expected: []string{"bab"}},
		"multi-match(all-match)":  {groups: []string{"a", "aa", "aba", "bab"}, limits: []string{"a", "aba"}, expected: []string{"a", "aba"}},
		"multi-match(some-match)": {groups: []string{"a", "aa", "aba", "bab"}, limits: []string{"a", ""}, expected: []string{"a"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotGroups, err := LimitGroups(tc.groups, tc.limits)
			assert.NilError(t, err)
			wantGroups := tc.expected
			sort.Strings(gotGroups)
			sort.Strings(wantGroups)
			assert.DeepEqual(t, wantGroups, gotGroups)
		})
	}

}
