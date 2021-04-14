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
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func NewKubeconfigFromConfig(config *invconfig.Kubeconfig) (*Kubeconfig, error) {
	return NewKubeconfigFromParams(config.Backend, config.Params)
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

func (k *Kubeconfig) SetNamespace(n string) error {
	k.namespace = n
	return k.loadConfig()
}

func (k *Kubeconfig) loadConfig() error {
	configData, err := k.loader.Load()
	if err != nil {
		return err
	}

	config, err := clientcmd.Load(configData)
	if err != nil {
		return err
	}

	if config.AuthInfos == nil {
		config.AuthInfos = map[string]*clientcmdapi.AuthInfo{}
	}
	if config.Clusters == nil {
		config.Clusters = map[string]*clientcmdapi.Cluster{}
	}
	if config.Contexts == nil {
		config.Contexts = map[string]*clientcmdapi.Context{}
	}
	if len(config.Contexts) > 0 {
		// normalize context names
		// the resulting contexts only include contexts with unique
		// cluster/user/namespace settings
		contexts := make(map[string]*clientcmdapi.Context, len(config.Contexts))
		for _, context := range config.Contexts {
			name := fmt.Sprintf("%s-%s", context.Cluster, context.AuthInfo)
			if context.Namespace != "" {
				name = fmt.Sprintf("%s-%s", name, context.Namespace)
			}
			if k.namespace != "" {
				// set namespace play config
				context.Namespace = k.namespace
			}

			contexts[name] = context
		}
		config.Contexts = contexts

		// If the current context is "", set it to the first
		// context we can retrieve from the config. Because
		// there is no guaranteed order of map elements, this is
		// not necessaryly the first context
		if config.CurrentContext == "" {
			for name := range config.Contexts {
				config.CurrentContext = name
				break
			}
		}
	}

	clientConfig := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{})

	k.config = clientConfig
	return nil
}
