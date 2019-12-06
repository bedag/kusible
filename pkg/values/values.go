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
	"os"

	groupsfilter "github.com/bedag/kusible/pkg/groups"
)

func New(path string, groups []string, skipEval bool, ejsonSettings EjsonSettings) (Values, error) {
	var result Values
	var err error

	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if stat.Mode().IsRegular() {
		// the path provided is a file, treat it as a single value
		// file, thus loading it with ejson and spruc operator support
		result, err = NewFile(path, false, ejsonSettings)
		if err != nil {
			return nil, err
		}
	} else {
		// The path provided is a directory, treat it as a values
		// directory. As the valuesDirectory type requires a list
		// of groups to determine which files to process, first
		// get a list of all groups in the given directory
		dirGroups := groups
		if len(dirGroups) <= 0 {
			dirGroups, err = groupsfilter.Groups(path, ".*", []string{})
			if err != nil {
				return nil, err
			}
		}
		result, err = NewDirectory(path, dirGroups, false, ejsonSettings)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
