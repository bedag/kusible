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

package values

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Shopify/ejson"
	"github.com/geofffranks/simpleyaml"
	"github.com/geofffranks/spruce"
	log "github.com/sirupsen/logrus"
)

// Compile takes a directory and a list of groups as parameters and
// compiles a map of values based on the files in the given directory
// where the file names / directory names on the top level of the directory
// match the given group names.
//
// Each group in the list may match either
//  * directories
//  * *.yaml, *.yml, *.json
//  * *.ejson
// in the given directory. It is not required that a group has any matching
// files / directories.
//
// The contents of the files / directories will then be merged according
// to the order of the given group list. Values of groups at the end of the
// list will override values from the end of the list (least specific to most
// specific ordering).
//
// If an entry in the given directory is itself a directory, its contents
// (including all subdirectories) will be merged in alphabetical order.
//
// All files / directories belonging to the same group or having the same
// basename (foo/, foo.yaml, foo.json all have the same basename) will
// be merged with the following priority (least to most specific)
//  * directories
//  * *.yaml, *.yml, *.json
//  * *.ejson
// Note: *.yaml, *.yml, *.json have the same priority so no guarantees
// for the merge order are made.
//
// Files can make use of spruce operators (https://github.com/geofffranks/spruce/blob/master/doc/operators.md).
// *.ejson will be treated as ejson (https://github.com/Shopify/ejson) encrypted
// and decrypted before merging if a matching private key was provided.
func Compile(directory string, groups []string, ejsonKeyDir string, ejsonPrivKey string, skipEval bool, skipDecrypt bool) (map[interface{}]interface{}, error) {
	// List of keys that should be pruned when running the merged data
	// through the spruce evaluator
	//   * _public_key: only present in encrypted ejson files to identify the correct private key
	//                  not required in resulting document
	pruneKeys := []string{"_public_key"}
	root := make(map[interface{}]interface{})

	// get the list of files that should be merged
	files, err := GetOrderedDataFileList(directory, groups)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"files": strings.Join(files[:], " "),
	}).Debug("Ordered list of files to merge")

	// merge everything while decrypting any ejson files encountered
	merger := &spruce.Merger{AppendByDefault: false}
	for _, path := range files {
		var data []byte

		// Check if the current path is an ejson file and if so, try
		// to decrypt it. If it cannot be decrypted, continue as there
		// is no harm in using the encrypted values
		isEjson, err := filepath.Match("*.ejson", filepath.Base(path))
		if err == nil && isEjson && !skipDecrypt {
			file, err := os.Open(path)
			if err != nil {
				return nil, err
			}
			defer file.Close()
			var outBuffer bytes.Buffer

			err = ejson.Decrypt(file, &outBuffer, ejsonKeyDir, ejsonPrivKey)
			if err != nil {
				log.WithFields(log.Fields{
					"file":  path,
					"error": err.Error(),
				}).Warn("Failed to decrypt ejson file, continuing with encrypted data")

				data, err = ioutil.ReadFile(path)
				if err != nil {
					return nil, err
				}
			} else {
				data = outBuffer.Bytes()
			}
		} else {
			data, err = ioutil.ReadFile(path)
			if err != nil {
				return nil, err
			}
		}

		yamlData, err := simpleyaml.NewYaml(data)
		if err != nil {
			return nil, err
		}

		doc, err := yamlData.Map()
		if err != nil {
			return nil, err
		}

		merger.Merge(root, doc)
	}

	if merger.Error() != nil {
		return nil, merger.Error()
	}

	evaluator := &spruce.Evaluator{Tree: root, SkipEval: skipEval}
	err = evaluator.Run(pruneKeys, nil)
	return evaluator.Tree, err
}

// GetOrderedDataFileList traverses the given directory and returns a list of
// files according to the rules described for the Compile method
func GetOrderedDataFileList(directory string, groups []string) ([]string, error) {
	var orderedFileList []string

	for _, group := range groups {
		var orderedGroupFileList []string
		groupDirectory := filepath.Join(directory, group)

		if stat, err := os.Stat(groupDirectory); err == nil && stat.Mode().IsDir() {
			// TODO: this adds directories in revers alphabetical order (files are fine though)
			err := filepath.Walk(groupDirectory, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					log.WithFields(log.Fields{
						"path": path,
					}).Warn(err.Error())
					return nil
				}

				if info.IsDir() && path != groupDirectory {
					files, _ := DirectoryDataFiles(path, "*")
					orderedGroupFileList = append(files, orderedGroupFileList...)
					return nil
				}

				return nil
			})
			if err != nil {
				return nil, err
			}
		}
		// add all files contained in subdirectories of the group directory
		// e.g. <directory>/<group>/**/*.{yml,yaml,json,ejson}
		orderedFileList = append(orderedFileList, orderedGroupFileList...)

		// add all files contained in the group directory
		// e.g. <directory>/<group>/*.{yml,yaml,json,ejson}
		files, _ := DirectoryDataFiles(groupDirectory, "*")
		orderedFileList = append(orderedFileList, files...)

		// add all group files
		// e.g. <directory>/<group>.{yml,yaml,json,ejson}
		files, _ = DirectoryDataFiles(directory, group)
		orderedFileList = append(orderedFileList, files...)
	}

	return orderedFileList, nil
}

// DirectoryDataFiles returns all data files of a given directory matching
// the provided pattern. Only the filetypes give in the description of the
// Compile method are considered. The operation is non-recursive.
//
// The pattern syntax is the same as the one for fmt.Match.
func DirectoryDataFiles(directory string, pattern string) ([]string, bool) {
	dataFileExt := [...]string{".yaml", ".yml", ".json", ".ejson"}
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
