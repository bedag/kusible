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
	"regexp"

	"github.com/bedag/kusible/pkg/groups"
)

// MatchLimits returns true if the groups of the inventory entry satisfy all given
// limits, which are treated as ^$ enclosed regex
func (e *entry) MatchLimits(limits []string) (bool, error) {
	// no limits -> all groups are valid
	if len(limits) <= 0 {
		return true, nil
	}

	// no groups -> no limit matches
	if len(e.Groups) <= 0 {
		return false, nil
	}

	for _, limit := range limits {
		regex, err := regexp.Compile("^" + limit + "$")
		if err != nil {
			return false, err
		}

		matched := false
		for _, group := range e.Groups {
			if regex.MatchString(group) {
				matched = true
				break
			}
		}
		if !matched {
			return false, nil
		}
	}
	return true, nil
}

// ValidGroups returns all groups of the inventory entry that satisfy at
// least one limit
func (e *entry) ValidGroups(limits []string) ([]string, error) {
	return groups.LimitGroups(e.Groups, limits)
}
