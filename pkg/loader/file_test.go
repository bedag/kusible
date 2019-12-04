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
package loader

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mitchellh/mapstructure"
	"gotest.tools/assert"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)

func TestFileBackendType(t *testing.T) {
	backend := &FileBackend{}
	assert.Equal(t, "file", backend.Type())
}

func TestFileBackendCreate(t *testing.T) {
	decryptKey := "aaaaa"
	path := "bbbbb"

	backend := NewFileBackend(path, decryptKey)
	if backend == nil {
		t.Errorf("failed to create file backend")
	}

	assert.Equal(t, decryptKey, backend.config.DecryptKey)
	assert.Equal(t, path, backend.config.Path)
}

func TestFileBackendCreateParamsNoEnv(t *testing.T) {
	decryptKey := "aaaaa"
	path := "bbbbb"

	params := map[string]interface{}{
		"decrypt_key": decryptKey,
		"path":        path,
	}

	backend, err := NewFileBackendFromParams(params)
	assert.NilError(t, err)

	assert.Equal(t, decryptKey, backend.config.DecryptKey)
	assert.Equal(t, path, backend.config.Path)
}

func TestFileBackendCreateParamsPartialEnv(t *testing.T) {
	decryptKey := "aaaaa"
	path := "bbbbb"

	params := map[string]interface{}{
		"path": path,
	}

	err := os.Setenv("EJSON_PRIVKEY", decryptKey)
	assert.NilError(t, err, "failed to set environment %s=%s", "EJSON_PRIVKEY", decryptKey)

	backend, err := NewFileBackendFromParams(params)
	assert.NilError(t, err)

	assert.Equal(t, decryptKey, backend.config.DecryptKey)
	assert.Equal(t, path, backend.config.Path)
}

func TestFileBackendCreateParamsFullEnv(t *testing.T) {
	decryptKey := "aaaaa"

	params := map[string]interface{}{}

	err := os.Setenv("EJSON_PRIVKEY", decryptKey)
	assert.NilError(t, err, "failed to set environment %s=%s", "EJSON_PRIVKEY", decryptKey)

	backend, err := NewFileBackendFromParams(params)
	assert.NilError(t, err)

	assert.Equal(t, decryptKey, backend.config.DecryptKey)
	assert.Equal(t, "kubeconfig", backend.config.Path)
}

func TestFileBackendLoad(t *testing.T) {
	decryptKey := "test123"
	path := "testdata/kubeconfig.enc"

	backend := NewFileBackend(path, decryptKey)
	if backend == nil {
		t.Errorf("failed to create file backend")
	}

	resultConfigBytesIn, err := backend.Load()
	assert.NilError(t, err)
	resultConfig, err := clientcmd.Load(resultConfigBytesIn)
	assert.NilError(t, err)
	resultConfigBytes, err := clientcmd.Write(*resultConfig)
	assert.NilError(t, err)

	expectedConfigPath := "testdata/kubeconfig"
	assert.NilError(t, err)
	expectedConfigBytesIn, err := ioutil.ReadFile(expectedConfigPath)
	assert.NilError(t, err)
	expectedConfig, err := clientcmd.Load(expectedConfigBytesIn)
	assert.NilError(t, err)
	expectedConfigBytes, err := clientcmd.Write(*expectedConfig)
	assert.NilError(t, err)
	assert.Equal(t, string(expectedConfigBytes), string(resultConfigBytes))
}

func TestFileConfig(t *testing.T) {
	params := map[string]interface{}{
		"decrypt_key": "aaaaa",
		"path":        "bbbbb",
	}

	backend, err := NewFileBackendFromParams(params)
	assert.NilError(t, err)

	var expected FileConfig
	var result FileConfig

	decoderConfig := &mapstructure.DecoderConfig{
		Result:  &expected,
		TagName: "yaml",
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	assert.NilError(t, err)
	err = decoder.Decode(params)
	assert.NilError(t, err)

	// TODO: unsafe vs. safe test
	resultRaw, err := backend.Config().Yaml(true)
	assert.NilError(t, err)

	err = yaml.Unmarshal(resultRaw, &result)
	assert.NilError(t, err)

	assert.DeepEqual(t, expected, result)
}
