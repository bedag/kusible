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
//
package inventory

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"k8s.io/client-go/tools/clientcmd"
)

func NewKubeconfigFromConfig(backend string, params map[interface{}]interface{}, skipLoading bool) (*kubeconfig, error) {
	var kubeconfig *kubeconfig

	var loaderParams map[string]string
	var loader kubeconfigLoader
	err := mapstructure.Decode(params, &loaderParams)
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(backend) {
	case "s3", "":
		// the default, if no specific kubeconfig backend was provided in the
		// inventory entry, is to load the kubeconfig from s3
		loader = NewKubeconfigS3LoaderFromParams(loaderParams)
	case "file":
		loader = NewKubeconfigFileLoaderFromParams(loaderParams)
	default:
		return nil, fmt.Errorf("unknown kubeconfig backend: %s", backend)
	}

	kubeconfig, err = NewKubeconfigFromLoader(loader, skipLoading)
	return kubeconfig, nil
}

func NewKubeconfigFromLoader(loader kubeconfigLoader, skipLoading bool) (*kubeconfig, error) {
	if loader == nil {
		return nil, fmt.Errorf("no kubeconfig loader defined")
	}

	var rawConfig []byte
	var err error
	if !skipLoading {
		// TODO split data loading and data decoding
		rawConfig, err = loader.Load()
		if err != nil {
			return nil, err
		}
	}

	// if rawConfig is empty because of skipLoading, clientcmd.Load()
	// returns an  empty config, see
	// https://github.com/kubernetes/client-go/blob/571c0ef67034a5e72b9e30e662044b770361641e/tools/clientcmd/loader.go#L408
	clientConfig, err := clientcmd.Load(rawConfig)
	if err != nil {
		return nil, err
	}

	kubeconfig := &kubeconfig{
		Loader: loader,
		Config: clientConfig,
	}

	return kubeconfig, nil
}

func (k *kubeconfig) Yaml() ([]byte, error) {
	config, err := clientcmd.Write(*k.Config)
	return config, err
}
