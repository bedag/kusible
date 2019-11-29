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
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func NewKubeconfigFromConfig(backend string, params map[string]interface{}) (*kubeconfig, error) {
	var kubeconfig *kubeconfig

	var loader kubeconfigLoader

	switch strings.ToLower(backend) {
	case "s3", "":
		// the default, if no specific kubeconfig backend was provided in the
		// inventory entry, is to load the kubeconfig from s3
		loader = NewKubeconfigS3LoaderFromParams(params)
	case "file":
		loader = NewKubeconfigFileLoaderFromParams(params)
	default:
		return nil, fmt.Errorf("unknown kubeconfig backend: %s", backend)
	}

	kubeconfig, err := NewKubeconfigFromLoader(loader)
	if err != nil {
		return nil, err
	}
	return kubeconfig, nil
}

func NewKubeconfigFromLoader(loader kubeconfigLoader) (*kubeconfig, error) {
	if loader == nil {
		return nil, fmt.Errorf("no kubeconfig loader defined")
	}

	kubeconfig := &kubeconfig{
		Loader: loader,
	}

	return kubeconfig, nil
}

func (k *kubeconfig) Yaml() ([]byte, error) {
	config, err := k.Config()
	if err != nil {
		return nil, err
	}
	data, err := clientcmd.Write(*config)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (k *kubeconfig) Config() (*clientcmdapi.Config, error) {
	if k.config == nil {
		err := k.loadConfig()
		if err != nil {
			return nil, err
		}
	}
	return k.config, nil
}

func (k *kubeconfig) loadConfig() error {
	rawConfig, err := k.Loader.Load()
	if err != nil {
		return err
	}

	clientConfig, err := clientcmd.Load(rawConfig)
	if err != nil {
		return err
	}
	if clientConfig.CurrentContext == "" {
		for name := range clientConfig.Contexts {
			clientConfig.CurrentContext = name
			break
		}
	}

	k.config = clientConfig
	return nil
}

func kubeconfigDecoderHookFunc(entryName string) mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t.Name() == "kubeconfig" {
			var config struct {
				Backend string
				Params  map[string]interface{}
			}
			err := mapstructure.Decode(data, &config)

			// keys starting with _ are treated as metadata by the kubeconfig loaders
			// add the name of the entry currently being decoded as metadata
			// so a loader can use it to construct its default values
			config.Params["_entry"] = entryName

			kubeconfig, err := NewKubeconfigFromConfig(config.Backend, config.Params)
			return kubeconfig, err
		}
		return data, nil
	}
}
