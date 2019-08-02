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

type kubeconfigLoader interface {
	Load() string
}

type kubeconfig struct {
	kubeconfig string
	loader     kubeconfigLoader
}

type inventory struct {
	entries []inventoryEntry
}

type inventoryEntry struct {
	name            string
	groups          []string
	tiller          tillerSettings
	configNamespace string
	kubeconfig      kubeconfig
}

type tillerSettings struct {
	namespace string
	tls       bool
	ca        string
	cert      string
	key       string
}

type kubeconfigS3Loader struct {
	accessKey  string
	secretKey  string
	region     string
	server     string
	decryptKey string
	bucket     string
	path       string
}

type kubeconfigFileLoader struct {
	path       string
	decryptKey string
}
