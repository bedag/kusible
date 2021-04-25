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

package ejson

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/Shopify/ejson"
	log "github.com/sirupsen/logrus"
)

func ReadFile(path string, settings Settings) ([]byte, error) {
	data := []byte{}

	// No decryption requested, just read the file
	if settings.SkipDecrypt {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	// try to decrypt the file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var outBuffer bytes.Buffer

	err = ejson.Decrypt(file, &outBuffer, settings.KeyDir, settings.PrivKey)
	if err == nil {
		data = outBuffer.Bytes()
		return data, nil
	}

	log.WithFields(log.Fields{
		"file":  path,
		"error": err.Error(),
	}).Warn("Failed to decrypt ejson file, continuing with encrypted data")

	data, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}
