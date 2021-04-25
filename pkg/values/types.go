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

package values

import "github.com/bedag/kusible/pkg/wrapper/ejson"

type Values interface {
	YAML() ([]byte, error)
	JSON() ([]byte, error)
	Map() map[string]interface{}
}

type file struct {
	data     map[string]interface{}
	path     string
	ejson    ejson.Settings
	skipEval bool
}

type directory struct {
	data            map[string]interface{}
	path            string
	groups          []string
	ejson           ejson.Settings
	skipEval        bool
	files           []file
	orderedFileList []string
}
