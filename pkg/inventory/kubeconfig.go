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
	"fmt"

	invconfig "github.com/bedag/kusible/pkg/inventory/config"
	"github.com/bedag/kusible/pkg/loader"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

func NewKubeconfigFromConfig(config *invconfig.Kubeconfig) (*Kubeconfig, error) {
	return NewKubeconfigFromParams(config.Backend, *config.Params)
}

func NewKubeconfigFromParams(backend string, params map[string]interface{}) (*Kubeconfig, error) {
	ldr, err := loader.New(backend, params)
	if err != nil {
		return nil, err
	}

	kubeconfig, err := NewKubeconfigFromLoader(ldr)
	if err != nil {
		return nil, err
	}
	return kubeconfig, nil
}

func NewKubeconfigFromLoader(ldr loader.Loader) (*Kubeconfig, error) {
	if ldr == nil {
		return nil, fmt.Errorf("no kubeconfig loader defined")
	}

	kubeconfig := &Kubeconfig{
		loader: ldr,
	}

	return kubeconfig, nil
}

func (k *Kubeconfig) Yaml() ([]byte, error) {
	clientConfig, err := k.Config()
	if err != nil {
		return nil, err
	}
	rawConfig, err := clientConfig.RawConfig()
	if err != nil {
		return nil, err
	}
	data, err := clientcmd.Write(rawConfig)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (k *Kubeconfig) Loader() loader.Loader {
	return k.loader
}

func (k *Kubeconfig) Config() (clientcmd.ClientConfig, error) {
	if k.config == nil {
		err := k.loadConfig()
		if err != nil {
			return nil, err
		}
	}
	return k.config, nil
}

func (k *Kubeconfig) SetClient(clientset kubernetes.Interface) {
	k.client = clientset
}

// Client returns a clientset for the current kubeconfig. If no client
// currently exists, a new one will be created
func (k *Kubeconfig) Client() (kubernetes.Interface, error) {
	if k.client != nil {
		return k.client, nil
	}

	config, err := k.Config()
	if err != nil {
		return nil, err
	}

	clientConfig, err := config.ClientConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	k.SetClient(clientset)

	return clientset, nil
}

func (k *Kubeconfig) loadConfig() error {
	configData, err := k.loader.Load()
	if err != nil {
		return err
	}

	rawConfig, err := clientcmd.Load(configData)
	if err != nil {
		return err
	}
	if rawConfig.CurrentContext == "" {
		for name := range rawConfig.Contexts {
			rawConfig.CurrentContext = name
			break
		}
	}

	clientConfig := clientcmd.NewDefaultClientConfig(*rawConfig, nil)

	k.config = clientConfig
	return nil
}
