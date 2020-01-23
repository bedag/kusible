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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	// Use geofffranks yaml library instead of go-yaml
	// to ensure compatibility with spruce
	"github.com/Shopify/ejson"
	"github.com/bedag/kusible/pkg/spruce"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"
)

func NewFile(path string, skipEval bool, ejsonSettings EjsonSettings) (*file, error) {
	result := &file{
		path:     path,
		ejson:    ejsonSettings,
		skipEval: skipEval,
	}
	err := result.loadMap()
	return result, err
}

func (f *file) load() ([]byte, error) {
	var data []byte
	// Check if the current path is an ejson file and if so, try
	// to decrypt it. If it cannot be decrypted, continue as there
	// is no harm in using the encrypted values
	isEjson, err := filepath.Match("*.ejson", filepath.Base(f.path))
	if err == nil && isEjson && !f.ejson.SkipDecrypt {
		file, err := os.Open(f.path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		var outBuffer bytes.Buffer

		err = ejson.Decrypt(file, &outBuffer, f.ejson.KeyDir, f.ejson.PrivKey)
		if err != nil {
			log.WithFields(log.Fields{
				"file":  f.path,
				"error": err.Error(),
			}).Warn("Failed to decrypt ejson file, continuing with encrypted data")

			data, err = ioutil.ReadFile(f.path)
			if err != nil {
				return nil, err
			}
		} else {
			data = outBuffer.Bytes()
		}
	} else {
		data, err = ioutil.ReadFile(f.path)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (f *file) loadMap() error {
	data, err := f.load()
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &f.data)
	if err != nil {
		return err
	}

	if f.data == nil {
		f.data = make(map[string]interface{})
	}

	// if we want to skip the spruce evaluation, skip the evaluator
	// alltogether as an Evaluator with SkipEval: true only prunes / cherrypicks,
	// something we do not need here
	if !f.skipEval {
		err := spruce.Eval(&f.data, false, []string{})
		return err
	}
	return nil
}

func (f *file) Map() map[string]interface{} {
	return f.data
}

func (f *file) YAML() ([]byte, error) {
	return yaml.Marshal(f.data)
}

func (f *file) JSON() ([]byte, error) {
	return json.Marshal(f.data)
}
