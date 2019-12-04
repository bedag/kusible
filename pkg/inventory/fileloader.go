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

package inventory

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"gopkg.in/yaml.v2"
)

func NewKubeconfigFileLoaderFromParams(params map[string]interface{}) *kubeconfigFileLoader {
	result := map[string]string{
		"decrypt_key": os.Getenv("EJSON_PRIVKEY"),
		"path":        "kubeconfig",
	}

	for k, v := range params {
		if !strings.HasPrefix(k, "_") {
			result[strings.ToLower(k)] = v.(string)
		}
	}

	return NewKubeconfigFileLoader(
		result["path"],
		result["decrypt_key"])
}

func NewKubeconfigFileLoader(path string, decryptKey string) *kubeconfigFileLoader {
	return &kubeconfigFileLoader{Path: path, DecryptKey: decryptKey}
}

func (loader *kubeconfigFileLoader) Load() ([]byte, error) {
	if loader.Path == "" {
		return nil, fmt.Errorf("no path set for kubeconfig file loader")
	}
	_, err := os.Stat(loader.Path)
	if err != nil {
		return nil, err
	}

	mime, _, err := mimetype.DetectFile(loader.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to detect mimetype for file://%s", loader.Path)
	}

	var rawKubeconfig []byte
	switch mime {
	case "text/plain":
		rawKubeconfig, err = ioutil.ReadFile(loader.Path)
		if err != nil {
			return nil, err
		}
	case "application/x-7z-compressed":
		rawKubeconfig, err = extractSingleTar7ZipFile(loader.Path, loader.DecryptKey)
		if err != nil {
			return nil, err
		}
	case "application/octet-stream":
		rawKubeconfig, err = decryptOpensslSymmetricFile(loader.Path, loader.DecryptKey)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("Unknown kubeconfig source file type: " + mime)
	}

	return rawKubeconfig, nil
}

func (loader *kubeconfigFileLoader) Type() string {
	return "file"
}

func (loader *kubeconfigFileLoader) Config(unsafe bool) map[string]interface{} {
	decryptKey := loader.DecryptKey
	if !unsafe {
		decryptKey = fmt.Sprintf("%x", sha256.Sum256([]byte(decryptKey)))
	}
	result := map[string]interface{}{
		"decrypt_key": decryptKey,
		"path":        loader.Path,
	}
	return result
}

func (loader *kubeconfigFileLoader) ConfigYaml(unsafe bool) ([]byte, error) {
	config := loader.Config(unsafe)
	result, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}
	return result, nil
}
