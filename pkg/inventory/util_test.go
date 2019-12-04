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
	"io/ioutil"
	"testing"

	"gotest.tools/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestExtractSingleTar7Zip(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/kubeconfig.enc.7z")
	assert.NilError(t, err)
	password := "test123"

	result, err := extractSingleTar7Zip(data, password)
	assert.NilError(t, err)
	_, err = clientcmd.Load(result)
	assert.NilError(t, err)
}

func TestExtractSingleTar7ZipFile(t *testing.T) {
	file := "testdata/kubeconfig.enc.7z"
	password := "test123"

	result, err := extractSingleTar7ZipFile(file, password)
	assert.NilError(t, err)
	_, err = clientcmd.Load(result)
	assert.NilError(t, err)
}

func TestDecryptOpensslSymmetric(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/kubeconfig.enc")
	assert.NilError(t, err)
	password := "test123"

	result, err := decryptOpensslSymmetric(data, password)
	assert.NilError(t, err)
	_, err = clientcmd.Load(result)
	assert.NilError(t, err)
}

func TestDecryptOpensslSymmetricFile(t *testing.T) {
	file := "testdata/kubeconfig.enc"
	password := "test123"

	result, err := decryptOpensslSymmetricFile(file, password)
	assert.NilError(t, err)
	_, err = clientcmd.Load(result)
	assert.NilError(t, err)
}
