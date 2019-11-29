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
	"github.com/bedag/kusible/pkg/values"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type Inventory struct {
	Entries entries `mapstructure:"inventory"`
	Ejson   *values.EjsonSettings
}

type Targets struct {
	limits     []string
	filter     string
	valuesPath string
	ejson      *values.EjsonSettings
	targets    map[string]*Target
}

type Target struct {
	entry  *entry
	values values.Values
}

type entries map[string]entry

type entry struct {
	Name            string   `mapstructure:"name"`
	Groups          []string `mapstructure:"groups"`
	ConfigNamespace string   `mapstructure:"config_namespace"`
	Kubeconfig      *kubeconfig
}

type kubeconfigLoader interface {
	Load() ([]byte, error)
	Type() string
	ConfigYaml(unsafe bool) ([]byte, error)
	Config(unsafe bool) map[string]interface{}
}

type kubeconfig struct {
	Loader kubeconfigLoader
	config *clientcmdapi.Config
}

type kubeconfigS3Loader struct {
	AccessKey  string `mapstructure:"accesskey"`
	SecretKey  string `mapstructure:"secretkey"`
	Region     string `mapstructure:"region"`
	Server     string `mapstructure:"server"`
	DecryptKey string `mapstructure:"decrypt_key"`
	Bucket     string `mapstructure:"bucket"`
	Path       string `mapstructure:"path"`
	Downloader s3manageriface.DownloaderAPI
}

type kubeconfigFileLoader struct {
	Path       string
	DecryptKey string
}
