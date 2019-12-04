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
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
	"gotest.tools/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestS3LoaderType(t *testing.T) {
	loader := &kubeconfigS3Loader{}
	assert.Equal(t, "s3", loader.Type())
}

func TestS3LoaderCreate(t *testing.T) {
	accessKey := "aaaaa"
	secretKey := "bbbbb"
	region := "ccccc"
	server := "ddddd"
	decryptKey := "eeeee"
	bucket := "fffff"
	path := "ggggg"

	loader := NewKubeconfigS3Loader(accessKey, secretKey, region, server, decryptKey, bucket, path)
	if loader == nil {
		t.Errorf("failed to create s3 loader")
	}

	assert.Equal(t, accessKey, loader.AccessKey)
	assert.Equal(t, secretKey, loader.SecretKey)
	assert.Equal(t, region, loader.Region)
	assert.Equal(t, server, loader.Server)
	assert.Equal(t, decryptKey, loader.DecryptKey)
	assert.Equal(t, bucket, loader.Bucket)
	assert.Equal(t, path, loader.Path)
}

func TestS3LoaderCreateParamsNoEnv(t *testing.T) {
	accessKey := "aaaaa"
	secretKey := "bbbbb"
	region := "ccccc"
	server := "ddddd"
	decryptKey := "eeeee"
	bucket := "fffff"
	path := "ggggg"

	params := map[string]interface{}{
		"accesskey":   accessKey,
		"secretkey":   secretKey,
		"region":      region,
		"server":      server,
		"decrypt_key": decryptKey,
		"bucket":      bucket,
		"path":        path,
	}

	loader := NewKubeconfigS3LoaderFromParams(params)
	if loader == nil {
		t.Errorf("failed to create s3 loader")
	}

	assert.Equal(t, accessKey, loader.AccessKey)
	assert.Equal(t, secretKey, loader.SecretKey)
	assert.Equal(t, region, loader.Region)
	assert.Equal(t, server, loader.Server)
	assert.Equal(t, decryptKey, loader.DecryptKey)
	assert.Equal(t, bucket, loader.Bucket)
	assert.Equal(t, path, loader.Path)
}

func TestS3LoaderCreateParamsPartialEnv(t *testing.T) {
	accessKey := "aaaaa"
	secretKey := "bbbbb"
	region := "ccccc"
	server := "ddddd"
	decryptKey := "eeeee"
	path := "ggggg"

	params := map[string]interface{}{
		"region": region,
		"path":   path,
	}

	err := os.Setenv("S3_ACCESSKEY", accessKey)
	assert.NilError(t, err, "failed to set environment %s=%s", "S3_ACCESSKEY", accessKey)
	err = os.Setenv("S3_SECRETKEY", secretKey)
	assert.NilError(t, err, "failed to set environment %s=%s", "S3_SECRETKEY", secretKey)
	err = os.Setenv("S3_SERVER", server)
	assert.NilError(t, err, "failed to set environment %s=%s", "S3_SERVER", server)
	err = os.Setenv("EJSON_PRIVKEY", decryptKey)
	assert.NilError(t, err, "failed to set environment %s=%s", "EJSON_PRIVKEY", decryptKey)

	loader := NewKubeconfigS3LoaderFromParams(params)
	if loader == nil {
		t.Errorf("failed to create s3 loader")
	}

	assert.Equal(t, accessKey, loader.AccessKey)
	assert.Equal(t, secretKey, loader.SecretKey)
	assert.Equal(t, region, loader.Region)
	assert.Equal(t, server, loader.Server)
	assert.Equal(t, decryptKey, loader.DecryptKey)
	assert.Equal(t, "kubernetes", loader.Bucket)
	assert.Equal(t, path, loader.Path)
}

func TestS3LoaderCreateParamsFullEnv(t *testing.T) {
	accessKey := "aaaaa"
	secretKey := "bbbbb"
	region := "ccccc"
	server := "ddddd"
	decryptKey := "eeeee"
	bucket := "fffff"

	params := map[string]interface{}{}

	err := os.Setenv("S3_ACCESSKEY", accessKey)
	assert.NilError(t, err, "failed to set environment %s=%s", "S3_ACCESSKEY", accessKey)
	err = os.Setenv("S3_SECRETKEY", secretKey)
	assert.NilError(t, err, "failed to set environment %s=%s", "S3_SECRETKEY", secretKey)
	err = os.Setenv("S3_REGION", region)
	assert.NilError(t, err, "failed to set environment %s=%s", "S3_REGION", region)
	err = os.Setenv("S3_SERVER", server)
	assert.NilError(t, err, "failed to set environment %s=%s", "S3_SERVER", server)
	err = os.Setenv("S3_BUCKET", bucket)
	assert.NilError(t, err, "failed to set environment %s=%s", "S3_BUCKET", bucket)
	err = os.Setenv("EJSON_PRIVKEY", decryptKey)
	assert.NilError(t, err, "failed to set environment %s=%s", "EJSON_PRIVKEY", decryptKey)

	loader := NewKubeconfigS3LoaderFromParams(params)
	if loader == nil {
		t.Errorf("failed to create s3 loader")
	}

	assert.Equal(t, accessKey, loader.AccessKey)
	assert.Equal(t, secretKey, loader.SecretKey)
	assert.Equal(t, region, loader.Region)
	assert.Equal(t, server, loader.Server)
	assert.Equal(t, decryptKey, loader.DecryptKey)
	assert.Equal(t, bucket, loader.Bucket)
	assert.Equal(t, "kubeconfig/kubeconfig.enc.7z", loader.Path)
}

func TestS3LoaderEntryDefault(t *testing.T) {
	entry := "test"
	defaultPath := "kubeconfig/kubeconfig.enc.7z"
	params := map[string]interface{}{
		"_entry": entry,
	}

	expected := fmt.Sprintf("%s/%s", entry, defaultPath)
	loader := NewKubeconfigS3LoaderFromParams(params)
	if loader == nil {
		t.Errorf("failed to create s3 loader")
	}

	assert.Equal(t, expected, loader.Path)
}

type mockedS3DownloadManager struct {
	s3manageriface.DownloaderAPI
}

func (d mockedS3DownloadManager) Download(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error) {
	return d.DownloadWithContext(aws.BackgroundContext(), w, input, options...)
}

func (d mockedS3DownloadManager) DownloadWithContext(ctx aws.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error) {
	path := fmt.Sprintf("%s/%s", aws.StringValue(input.Bucket), aws.StringValue(input.Key))
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	count, err := w.WriteAt(content, 0)
	return int64(count), err
}

func TestS3LoaderLoad(t *testing.T) {
	decryptKey := "test123"
	bucket := "testdata"
	path := "kubeconfig.enc"

	loader := &kubeconfigS3Loader{
		AccessKey:  "foo",
		SecretKey:  "foo",
		Server:     "foo",
		DecryptKey: decryptKey,
		Bucket:     bucket,
		Path:       path,
		Downloader: mockedS3DownloadManager{},
	}

	resultConfigBytesIn, err := loader.Load()
	assert.NilError(t, err)
	resultConfig, err := clientcmd.Load(resultConfigBytesIn)
	assert.NilError(t, err)
	resultConfigBytes, err := clientcmd.Write(*resultConfig)
	assert.NilError(t, err)

	expectedConfigPath := fmt.Sprintf("%s/%s", bucket, "kubeconfig")
	assert.NilError(t, err)
	expectedConfigBytesIn, err := ioutil.ReadFile(expectedConfigPath)
	assert.NilError(t, err)
	expectedConfig, err := clientcmd.Load(expectedConfigBytesIn)
	assert.NilError(t, err)
	expectedConfigBytes, err := clientcmd.Write(*expectedConfig)
	assert.NilError(t, err)
	assert.Equal(t, string(expectedConfigBytes), string(resultConfigBytes))
}

func TestS3LoaderConfig(t *testing.T) {
	params := map[string]interface{}{
		"accesskey":   "aaaaa",
		"secretkey":   "bbbbb",
		"region":      "ccccc",
		"server":      "ddddd",
		"decrypt_key": "eeeee",
		"bucket":      "fffff",
		"path":        "ggggg",
	}

	loader := NewKubeconfigS3LoaderFromParams(params)
	if loader == nil {
		t.Errorf("failed to create s3 loader")
	}

	type config struct {
		AccessKey  string `yaml:"accesskey"`
		SecretKey  string `yaml:"secretkey"`
		Region     string `yaml:"region"`
		Server     string `yaml:"server"`
		DecryptKey string `yaml:"decrypt_key"`
		Bucket     string `yaml:"bucket"`
		Path       string `yaml:"path"`
	}

	var expected config
	var result config

	decoderConfig := &mapstructure.DecoderConfig{
		Result:  &expected,
		TagName: "yaml",
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	assert.NilError(t, err)
	err = decoder.Decode(params)
	assert.NilError(t, err)

	// TODO: unsafe vs. safe test
	resultRaw, err := loader.ConfigYaml(true)
	assert.NilError(t, err)
	err = yaml.Unmarshal(resultRaw, &result)
	assert.NilError(t, err)

	assert.DeepEqual(t, expected, result)
}
