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
	"errors"
	"io/ioutil"
	"os"

	"github.com/gabriel-vasile/mimetype"
)

func NewKubeconfigFileLoader(path string, decryptKey string) *kubeconfigFileLoader {
	return &kubeconfigFileLoader{path: path, decryptKey: decryptKey}
}

func (loader *kubeconfigFileLoader) Load() ([]byte, error) {
	_, err := os.Stat(loader.path)
	if err != nil {
		return nil, err
	}

	mime, _, err := mimetype.DetectFile(loader.path)
	if err != nil {
		return nil, err
	}

	var rawKubeconfig []byte
	switch mime {
	case "text/plain":
		rawKubeconfig, err = ioutil.ReadFile(loader.path)
		if err != nil {
			return nil, err
		}
	case "application/x-7z-compressed":
		rawKubeconfig, err = extractSingleTar7ZipFile(loader.path, loader.decryptKey)
		if err != nil {
			return nil, err
		}
	case "application/octet-stream":
		rawKubeconfig, err = decryptOpensslSymmetricFile(loader.path, loader.decryptKey)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("Unknown kubeconfig source file type: " + mime)
	}

	return rawKubeconfig, nil
}
