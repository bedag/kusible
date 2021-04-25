/*
Copyright © 2021 Bedag Informatik AG

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
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

/*
Groups returns an unsorted list of available group in the given directory
limited by the provided filter. Each element of the list given
in filter will be treated as a regex and only groups matching
any (!) of the filters will be returned. The regexp will be treated
as matching the whole group name, e.g. they are implicitely wrapped
in ^$
*/
func Groups(directory string, filter string, limits []string) ([]string, error) {
	groupFileExt := []string{"*.yaml", "*.yml", "*.json", "*.ejson"}

	stat, err := os.Stat(directory)
	if err != nil {
		return nil, err
	}

	if !stat.Mode().IsDir() {
		return nil, errors.New("'" + directory + "' is not a directory")
	}

	// get all elements in the given directory
	elements, err := filepath.Glob(filepath.Join(directory, "*"))
	if err != nil {
		return nil, err
	}

	groupSet := make(map[string]bool)

	for _, element := range elements {
		isGroupFile, err := filenameMultiMatch(groupFileExt, element)
		if err != nil {
			return nil, err
		}

		stat, err := os.Stat(element)
		if err != nil {
			return nil, err
		}

		if isGroupFile || stat.Mode().IsDir() {
			basename := filepath.Base(element)
			extension := filepath.Ext(basename)
			groupName := basename[0 : len(basename)-len(extension)]

			if !groupSet[groupName] {
				valid, err := groupRegexMatch([]string{filter}, groupName)
				if err != nil {
					return nil, err
				}

				if valid {
					groupSet[groupName] = true
				}
			}
		}
	}

	groups := make([]string, 0, len(groupSet))
	for group := range groupSet {
		groups = append(groups, group)
	}

	var result []string
	if len(limits) > 0 {
		result, err = LimitGroups(groups, limits)
		if err != nil {
			return nil, err
		}
	} else {
		result = groups
	}
	return result, nil
}

/*
SortedGroups is the same as Groups() but the resulting list is sorted alphabetically
*/
func SortedGroups(directory string, filter string, limits []string) ([]string, error) {
	g, err := Groups(directory, filter, limits)
	if err != nil {
		return nil, err
	}
	sort.Strings(g)
	return g, nil
}

// LimitGroups applies a list of limits (each treated as regex) to a list
// of groups and returns only groups matching at least one of the given limits
func LimitGroups(groups []string, limits []string) ([]string, error) {
	result := []string{}

	for _, group := range groups {
		valid, err := groupRegexMatch(limits, group)
		if err != nil {
			return nil, err
		}

		if valid {
			result = append(result, group)
		}
	}

	return result, nil
}

// filenameMultiMatch matches a given file name against a list of patterns
// The patterns use the same syntax as filepath.Match. For any given
// path, only the last element will be matched agains the pattern
func filenameMultiMatch(patterns []string, path string) (bool, error) {
	if len(patterns) <= 0 {
		return true, nil
	}

	basename := filepath.Base(path)

	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, basename)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}

	return false, nil
}

// groupRegexMatch takes a list of patterns that will be treated as regexp
// and match the provided group name agains these regexp. Each pattern will
// be wrapped in ^$ to match the whole string
//
// For the specific regexp syntax used, see https://github.com/google/re2/wiki/Syntax
func groupRegexMatch(patterns []string, group string) (bool, error) {
	if len(patterns) <= 0 {
		return true, nil
	}

	for _, pattern := range patterns {
		regex, err := regexp.Compile("^" + pattern + "$")
		if err != nil {
			return false, err
		}

		if regex.MatchString(group) {
			return true, nil
		}

	}

	return false, nil
}
