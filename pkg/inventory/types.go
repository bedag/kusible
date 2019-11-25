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
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type Inventory struct {
	Entries []inventoryEntry `mapstructure:"inventory"`
}

type kubeconfigLoader interface {
	Load() ([]byte, error)
	Type() string
	Config() []byte
}

type kubeconfig struct {
	Loader kubeconfigLoader
	Config *clientcmdapi.Config
}

type inventoryEntry struct {
	Name            string
	Groups          []string
	Tiller          tillerSettings
	ConfigNamespace string `mapstructure:"config_namespace"`
	Kubeconfig      kubeconfig
}

type tillerSettings struct {
	Namespace string
	TLS       bool
	CA        string
	Cert      string
	Key       string
}

type kubeconfigS3Loader struct {
	AccessKey  string
	SecretKey  string
	Region     string
	Server     string
	DecryptKey string
	Bucket     string
	Path       string
	Downloader s3manageriface.DownloaderAPI
}

type kubeconfigFileLoader struct {
	Path       string
	DecryptKey string
}
