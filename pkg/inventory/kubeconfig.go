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

func NewKubeconfigFromBackend(backend string, params map[interface{}]interface{}) (*kubeconfig, error) {
	var kubeconfig *kubeconfig

	switch strings.ToLower(backend) {
	default:
		return nil, fmt.Errorf("unknown kubeconfig backend: %s", backend)
	case "s3":
		var loader kubeconfigS3Loader
		err := mapstructure.Decode(params, &loader)
		if err != nil {
			return nil, err
		}
		kubeconfig, err = NewKubeconfigFromLoader(&loader)
	case "file":
		var loader kubeconfigFileLoader
		err := mapstructure.Decode(params, &loader)
		if err != nil {
			return nil, err
		}
		kubeconfig, err = NewKubeconfigFromLoader(&loader)
	}
	return kubeconfig, nil
}

func NewKubeconfigFromLoader(loader kubeconfigLoader) (*kubeconfig, error) {
	if loader == nil {
		return nil, fmt.Errorf("no kubeconfig loader defined")
	}

	// TODO split data loading an data decoding
	rawConfig, err := loader.Load()
	if err != nil {
		return nil, err
	}

	clientConfig, err := clientcmd.Load(rawConfig)
	if err != nil {
		return nil, err
	}

	kubeconfig := &kubeconfig{
		config: clientConfig,
	}

	return kubeconfig, nil
}

func (k *kubeconfig) Config() ([]byte, error) {
	config, err := clientcmd.Write(*k.config)
	return config, err
}
