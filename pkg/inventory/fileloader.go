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
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/gabriel-vasile/mimetype"
)

func NewKubeconfigFileLoader(path string, decryptKey string) *kubeconfigFileLoader {
	return &kubeconfigFileLoader{path: path, decryptKey: decryptKey}
}

func (loader *kubeconfigFileLoader) Load() (string, error) {
	_, err := os.Stat(loader.path)
	if err != nil {
		return "", err
	}

	mime, _, err := mimetype.DetectFile(loader.path)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer

	switch mime {
	case "text/plain":
		file, err := os.Open(loader.path)
		if err != nil {
			return "", err
		}
		defer file.Close()
		if _, err := io.Copy(&buffer, file); err != nil {
			return "", err
		}
	case "application/x-7z-compressed":
		buffer, err = extractSingleTar7ZipFile(loader.path, loader.decryptKey)
	default:
		return "", errors.New("Unknown kubeconfig source file type: " + mime)
	}

	kubeconfig := buffer.String()
	return kubeconfig, nil
}
