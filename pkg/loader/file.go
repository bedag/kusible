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

package loader

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gabriel-vasile/mimetype"
)

func NewFileBackend(path string, decryptKey string) *FileBackend {
	config := &FileConfig{
		DecryptKey: decryptKey,
		Path:       path,
	}
	return NewFileBackendFromConfig(config)
}

func NewFileBackendFromConfig(config *FileConfig) *FileBackend {
	return &FileBackend{
		config: config,
	}
}

func NewFileBackendFromParams(params map[string]interface{}) (*FileBackend, error) {
	config := FileConfig{
		DecryptKey: os.Getenv("EJSON_PRIVKEY"),
		Path:       "kubeconfig",
	}

	err := decode(params, &config)
	if err != nil {
		return nil, err
	}

	return NewFileBackendFromConfig(&config), nil
}

func (b *FileBackend) Load() ([]byte, error) {
	if b.config.Path == "" {
		return nil, fmt.Errorf("no path set for file backend")
	}
	_, err := os.Stat(b.config.Path)
	if err != nil {
		return nil, err
	}

	mime, err := mimetype.DetectFile(b.config.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to detect mimetype for file://%s", b.config.Path)
	}

	var raw []byte
	if mime.Is("text/plain") {
		raw, err = ioutil.ReadFile(b.config.Path)
		if err != nil {
			return nil, err
		}
	} else if mime.Is("application/x-7z-compressed") {
		raw, err = extractSingleTar7ZipFile(b.config.Path, b.config.DecryptKey)
		if err != nil {
			return nil, err
		}
	} else if mime.Is("application/octet-stream") {
		raw, err = decryptOpensslSymmetricFile(b.config.Path, b.config.DecryptKey)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Unknown source file type: " + mime.String())
	}

	return raw, nil
}

func (b *FileBackend) Type() string {
	return "file"
}

func (b *FileBackend) Config() BackendConfig {
	return b.config
}

func (c *FileConfig) Sanitize() BackendConfig {
	result := &FileConfig{
		DecryptKey: fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(c.DecryptKey))),
		Path:       c.Path,
	}
	return result
}

func (c *FileConfig) Yaml(unsafe bool) ([]byte, error) {
	return safeYaml(c, unsafe)
}
