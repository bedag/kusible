// Copyright © 2019 Michael Gruener
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package inventory

import (
	"testing"

	"gotest.tools/assert"
)

func TestEntryMatchLimits(t *testing.T) {
	limitsMatching := []string{"dev", "test"}
	limitsNotMatching := []string{"dev", "foo"}
	entry := &inventoryEntry{
		Name:   "test",
		Groups: []string{"dev", "test", "state", "prod"},
	}

	result, err := entry.MatchLimits(limitsMatching)
	assert.NilError(t, err)
	assert.Assert(t, result)

	result, err = entry.MatchLimits(limitsNotMatching)
	assert.NilError(t, err)
	assert.Assert(t, !result)
}

func TestEntryValidGroups(t *testing.T) {
	limits := []string{"dev", "test"}
	entry := &inventoryEntry{
		Name:   "test",
		Groups: []string{"dev", "test", "state", "prod"},
	}

	result, err := entry.ValidGroups(limits)
	assert.NilError(t, err)
	assert.DeepEqual(t, limits, result)
}
