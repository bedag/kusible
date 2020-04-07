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
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/bedag/kusible/internal/wrapper/spruce"
	"github.com/bedag/kusible/pkg/wrapper/ejson"
	"github.com/imdario/mergo"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"
)

func NewDirectory(path string, groups []string, skipEval bool, ejsonSettings ejson.Settings) (*directory, error) {
	result := &directory{
		path:            path,
		ejson:           ejsonSettings,
		skipEval:        skipEval,
		groups:          groups,
		orderedFileList: []string{},
		data:            map[string]interface{}{},
	}
	err := result.load()
	return result, err
}

/*
LoadMap takes a directory and a list of groups as parameters and
compiles a map of values based on the files in the given directory
where the file names / directory names on the top level of the directory
match the given group names.

Each group in the list may match either

 * directories
 * *.yml
 * *.yaml
 * *.json
 * *.ejson

in the given directory. It is not required that a group has any matching
files / directories.

The contents of the files / directories will then be merged according
to the order of the given group list. Values of groups at the end of the
list will override values from the end of the list (least specific to most
specific ordering).

If an entry in the given directory is itself a directory, its contents
(including all subdirectories) will be merged in alphabetical order.

All files / directories belonging to the same group or having the same
basename (foo/, foo.yaml, foo.json all have the same basename) will
be merged with the following priority (least to most specific)

 * directories
 * *.yml
 * *.yaml
 * *.json
 * *.ejson

Files can make use of spruce operators (https://github.com/geofffranks/spruce/blob/master/doc/operators.md).
*.ejson will be treated as ejson (https://github.com/Shopify/ejson) encrypted
and decrypted before merging if a matching private key was provided.
*/
func (d *directory) load() error {
	if len(d.data) > 0 {
		return nil
	}
	// List of keys that should be pruned when running the merged data
	// through the spruce evaluator
	//   * _public_key: only present in encrypted ejson files to identify the correct private key
	//                  not required in resulting document
	pruneKeys := []string{"_public_key"}

	var err error
	// get the list of files that should be merged
	err = d.createOrderedDataFileList()
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"files": strings.Join(d.orderedFileList[:], " "),
	}).Debug("Ordered list of files to merge")

	// merge everything while decrypting any ejson files encountered
	for _, path := range d.orderedFileList {
		file, err := NewFile(path, true, d.ejson)
		if err != nil {
			return err
		}
		doc := file.Map()
		err = mergo.Merge(&d.data, doc, mergo.WithOverride)
		if err != nil {
			return err
		}
	}

	err = spruce.Eval(&d.data, d.skipEval, pruneKeys)
	return err
}

func (d *directory) createOrderedDataFileList() error {
	for _, group := range d.groups {
		var orderedGroupFileList []string
		groupDirectory := filepath.Join(d.path, group)

		if stat, err := os.Stat(groupDirectory); err == nil && stat.Mode().IsDir() {
			err := filepath.Walk(groupDirectory, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					log.WithFields(log.Fields{
						"path": path,
					}).Warn(err.Error())
					return nil
				}

				if info.IsDir() && path != groupDirectory {
					files, _ := DirectoryDataFiles(path, "*")
					orderedGroupFileList = append(orderedGroupFileList, files...)
					return nil
				}

				return nil
			})
			if err != nil {
				return err
			}
		}
		// add all files contained in subdirectories of the group directory
		// e.g. <directory>/<group>/**/*.{yml,yaml,json,ejson}
		d.orderedFileList = append(d.orderedFileList, orderedGroupFileList...)

		// add all files contained in the group directory
		// e.g. <directory>/<group>/*.{yml,yaml,json,ejson}
		files, _ := DirectoryDataFiles(groupDirectory, "*")
		d.orderedFileList = append(d.orderedFileList, files...)

		// add all group files
		// e.g. <directory>/<group>.{yml,yaml,json,ejson}
		files, _ = DirectoryDataFiles(d.path, group)
		d.orderedFileList = append(d.orderedFileList, files...)
	}
	return nil
}

/*
OrderedDataFileList traverses the given directory and returns a list of
files according to the rules described for the Compile method
*/
func (d *directory) OrderedDataFileList() []string {
	return d.orderedFileList
}

func (d *directory) Map() map[string]interface{} {
	return d.data
}

func (d *directory) YAML() ([]byte, error) {
	return yaml.Marshal(d.data)
}

func (d *directory) JSON() ([]byte, error) {
	return json.Marshal(d.data)
}
