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

type EjsonSettings struct {
	KeyDir      string
	PrivKey     string
	SkipDecrypt bool
}

type Values interface {
	YAML() ([]byte, error)
	JSON() ([]byte, error)
	Map() (map[string]interface{}, error)
	Raw() *data
}

type data map[interface{}]interface{}

type file struct {
	data     data
	path     string
	ejson    EjsonSettings
	skipEval bool
}

type directory struct {
	data            data
	path            string
	groups          []string
	ejson           EjsonSettings
	skipEval        bool
	files           []file
	orderedFileList []string
}
