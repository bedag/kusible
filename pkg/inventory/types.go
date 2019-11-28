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
	"crypto/rsa"
	"crypto/x509"

	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type Inventory struct {
	Entries entries `mapstructure:"inventory"`
}

type entries map[string]entry

type entry struct {
	Name            string   `mapstructure:"name"`
	Groups          []string `mapstructure:"groups"`
	Tiller          *tiller  `mapstructure:"tiller"`
	ConfigNamespace string   `mapstructure:"config_namespace"`
	Kubeconfig      *kubeconfig
}

type kubeconfigLoader interface {
	Load() ([]byte, error)
	Type() string
	Config() ([]byte, error)
}

type kubeconfig struct {
	Loader kubeconfigLoader
	Config *clientcmdapi.Config
}

type tiller struct {
	Namespace string            `mapstructure:"namespace"`
	TLS       bool              `mapstructure:"tls"`
	CA        *x509.Certificate `mapstrucutre:"ca"`
	Cert      *x509.Certificate `mapstructure:"cert"`
	Key       *rsa.PrivateKey   `mapstrucutre:"key"`
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
