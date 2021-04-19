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

package values

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

/*
DirectoryDataFiles returns all data files of a given directory matching
the provided pattern. Only the filetypes give in the description of the
LoadMap method are considered. The operation is non-recursive.

The pattern syntax is the same as the one for fmt.Match.
*/
func DirectoryDataFiles(directory string, pattern string) ([]string, bool) {
	dataFileExt := [...]string{".yml", ".yaml", ".json", ".ejson"}
	var dataFileGlobs []string

	for _, ext := range dataFileExt {
		dataFileGlobs = append(dataFileGlobs, pattern+ext)
	}

	var fileList []string
	ok := true

	for _, glob := range dataFileGlobs {
		files, err := filepath.Glob(filepath.Join(directory, glob))
		if err != nil {
			log.WithFields(log.Fields{
				"pattern": glob,
			}).Warn(err.Error())
			ok = false
		} else {
			fileList = append(fileList, files...)
		}
	}

	return fileList, ok
}
