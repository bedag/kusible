/*
Copyright © 2019 Copyright © 2021 Bedag Informatik AG

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
	"github.com/bedag/kusible/pkg/inventory/config"
	"github.com/bedag/kusible/pkg/loader"
	"github.com/bedag/kusible/pkg/wrapper/ejson"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Inventory struct {
	entries map[string]*Entry
	ejson   *ejson.Settings
}

type Entry struct {
	name                   string
	groups                 []string
	clusterInventoryConfig *config.ClusterInventory
	kubeconfig             *Kubeconfig
}

type Kubeconfig struct {
	loader    loader.Loader
	config    clientcmd.ClientConfig
	client    kubernetes.Interface // *kubernetes.Clientset
	namespace string
}
