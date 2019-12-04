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

package loader

import (
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"gopkg.in/yaml.v2"
)

// Loader defines the loader interface
type Loader interface {
	Load() ([]byte, error) // loads the config from the source
	Type() string          // returns the loader backend type
	Config() BackendConfig // returns the backend config of the loader
}

type BackendConfig interface {
	Yaml(unsafe bool) ([]byte, error) // returns the sanitized loader config as yaml
	Sanitize() BackendConfig          // returns the sanitized loader config
}

type S3Config struct {
	BackendConfig
	AccessKey  string `json:"accesskey"`
	SecretKey  string `json:"secretkey"`
	Region     string `json:"region"`
	Server     string `json:"server"`
	DecryptKey string `json:"decrypt_key"`
	Bucket     string `json:"bucket"`
	Path       string `json:"path"`
}

type S3Backend struct {
	config     *S3Config
	Downloader s3manageriface.DownloaderAPI
}

type FileConfig struct {
	BackendConfig
	Path       string `json:"path"`
	DecryptKey string `json:"decrypt_key"`
}

type FileBackend struct {
	config *FileConfig
}

func safeYaml(c BackendConfig, unsafe bool) ([]byte, error) {
	config := c
	if !(unsafe) {
		config = config.Sanitize()
	}
	result, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}
	return result, nil
}
