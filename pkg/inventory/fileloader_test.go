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
	"io/ioutil"
	"os"
	"testing"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestFileLoaderType(t *testing.T) {
	loader := &kubeconfigFileLoader{}
	assert.Equal(t, "file", loader.Type())
}

func TestFileLoaderCreate(t *testing.T) {
	decryptKey := "aaaaa"
	path := "bbbbb"

	loader := NewKubeconfigFileLoader(path, decryptKey)
	if loader == nil {
		t.Errorf("failed to create file loader")
	}

	assert.Equal(t, decryptKey, loader.DecryptKey)
	assert.Equal(t, path, loader.Path)
}

func TestFileLoaderCreateParamsNoEnv(t *testing.T) {
	decryptKey := "aaaaa"
	path := "bbbbb"

	params := map[string]interface{}{
		"decrypt_key": decryptKey,
		"path":        path,
	}

	loader := NewKubeconfigFileLoaderFromParams(params)
	if loader == nil {
		t.Errorf("failed to create file loader")
	}

	assert.Equal(t, decryptKey, loader.DecryptKey)
	assert.Equal(t, path, loader.Path)
}

func TestFileLoaderCreateParamsPartialEnv(t *testing.T) {
	decryptKey := "aaaaa"
	path := "bbbbb"

	params := map[string]interface{}{
		"path": path,
	}

	err := os.Setenv("EJSON_PRIVKEY", decryptKey)
	assert.NilError(t, err, "failed to set environment %s=%s", "EJSON_PRIVKEY", decryptKey)

	loader := NewKubeconfigFileLoaderFromParams(params)
	if loader == nil {
		t.Errorf("failed to create file loader")
	}

	assert.Equal(t, decryptKey, loader.DecryptKey)
	assert.Equal(t, path, loader.Path)
}

func TestFileLoaderCreateParamsFullEnv(t *testing.T) {
	decryptKey := "aaaaa"

	params := map[string]interface{}{}

	err := os.Setenv("EJSON_PRIVKEY", decryptKey)
	assert.NilError(t, err, "failed to set environment %s=%s", "EJSON_PRIVKEY", decryptKey)

	loader := NewKubeconfigFileLoaderFromParams(params)
	if loader == nil {
		t.Errorf("failed to create file loader")
	}

	assert.Equal(t, decryptKey, loader.DecryptKey)
	assert.Equal(t, "kubeconfig", loader.Path)
}

func TestFileLoaderLoad(t *testing.T) {
	decryptKey := "test123"
	path := "testdata/kubeconfig.enc"

	loader := NewKubeconfigFileLoader(path, decryptKey)
	if loader == nil {
		t.Errorf("failed to create file loader")
	}

	resultConfigBytesIn, err := loader.Load()
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

	loader := NewKubeconfigFileLoaderFromParams(params)
	if loader == nil {
		t.Errorf("failed to create file loader")
	}

	type config struct {
		DecryptKey string `yaml:"decrypt_key"`
		Path       string `yaml:"path"`
	}

	var expected config
	var result config

	decoderConfig := &mapstructure.DecoderConfig{
		Result:  &expected,
		TagName: "yaml",
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	assert.NilError(t, err)
	err = decoder.Decode(params)
	assert.NilError(t, err)

	resultRaw, err := loader.Config()
	assert.NilError(t, err)

	err = yaml.Unmarshal(resultRaw, &result)
	assert.NilError(t, err)

	assert.DeepEqual(t, expected, result)
}
