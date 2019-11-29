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
	"testing"

	"gotest.tools/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestNewKubeconfigFromConfig(t *testing.T) {
	backend := "file"
	params := map[string]interface{}{
		"decrypt_key": "test123",
		"path":        "testdata/kubeconfig.enc.7z",
	}

	kubeconfig, err := NewKubeconfigFromConfig(backend, params)
	assert.NilError(t, err)
	assert.Equal(t, "file", kubeconfig.Loader.Type())
	resultConfigBytes, err := kubeconfig.Yaml()
	assert.NilError(t, err)
	resultClientConfig, err := kubeconfig.Config()
	assert.NilError(t, err)
	resultCurrentContext := resultClientConfig.CurrentContext
	assert.Assert(t, resultCurrentContext != "")

	expectedConfigPath := "testdata/kubeconfig"
	assert.NilError(t, err)
	expectedConfigBytesIn, err := ioutil.ReadFile(expectedConfigPath)
	assert.NilError(t, err)
	expectedConfig, err := clientcmd.Load(expectedConfigBytesIn)
	assert.NilError(t, err)
	if expectedConfig.CurrentContext == "" {
		expectedConfig.CurrentContext = resultCurrentContext
	}
	expectedConfigBytes, err := clientcmd.Write(*expectedConfig)
	assert.NilError(t, err)
	assert.Equal(t, string(expectedConfigBytes), string(resultConfigBytes))
}

func TestNewKubeconfigFromLoader(t *testing.T) {
	params := map[string]interface{}{
		"decrypt_key": "test123",
		"path":        "testdata/kubeconfig.enc.7z",
	}

	loader := NewKubeconfigFileLoaderFromParams(params)
	if loader == nil {
		t.Errorf("failed to create file loader")
	}

	kubeconfig, err := NewKubeconfigFromLoader(loader)
	assert.NilError(t, err)
	assert.Equal(t, "file", kubeconfig.Loader.Type())
	resultConfigBytes, err := kubeconfig.Yaml()
	assert.NilError(t, err)
	resultClientConfig, err := kubeconfig.Config()
	assert.NilError(t, err)
	resultCurrentContext := resultClientConfig.CurrentContext
	assert.Assert(t, resultCurrentContext != "")

	expectedConfigPath := "testdata/kubeconfig"
	assert.NilError(t, err)
	expectedConfigBytesIn, err := ioutil.ReadFile(expectedConfigPath)
	assert.NilError(t, err)
	expectedConfig, err := clientcmd.Load(expectedConfigBytesIn)
	assert.NilError(t, err)
	if expectedConfig.CurrentContext == "" {
		expectedConfig.CurrentContext = resultCurrentContext
	}
	expectedConfigBytes, err := clientcmd.Write(*expectedConfig)
	assert.NilError(t, err)
	assert.Equal(t, string(expectedConfigBytes), string(resultConfigBytes))
}
