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
	"gotest.tools/assert"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)

func TestS3BackendType(t *testing.T) {
	backend := &S3Backend{}
	assert.Equal(t, "s3", backend.Type())
}

func TestS3BackendCreate(t *testing.T) {
	accessKey := "aaaaa"
	secretKey := "bbbbb"
	region := "ccccc"
	server := "ddddd"
	decryptKey := "eeeee"
	bucket := "fffff"
	path := "ggggg"

	backend, err := NewS3Backend(accessKey, secretKey, region, server, decryptKey, bucket, path)
	assert.NilError(t, err)

	assert.Equal(t, accessKey, backend.config.AccessKey)
	assert.Equal(t, secretKey, backend.config.SecretKey)
	assert.Equal(t, region, backend.config.Region)
	assert.Equal(t, server, backend.config.Server)
	assert.Equal(t, decryptKey, backend.config.DecryptKey)
	assert.Equal(t, bucket, backend.config.Bucket)
	assert.Equal(t, path, backend.config.Path)
}

func TestS3BackendCreateParamsNoEnv(t *testing.T) {
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

	backend, err := NewS3BackendFromParams(params)
	assert.NilError(t, err)

	assert.Equal(t, accessKey, backend.config.AccessKey)
	assert.Equal(t, secretKey, backend.config.SecretKey)
	assert.Equal(t, region, backend.config.Region)
	assert.Equal(t, server, backend.config.Server)
	assert.Equal(t, decryptKey, backend.config.DecryptKey)
	assert.Equal(t, bucket, backend.config.Bucket)
	assert.Equal(t, path, backend.config.Path)
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

	backend, err := NewS3BackendFromParams(params)
	assert.NilError(t, err)

	assert.Equal(t, accessKey, backend.config.AccessKey)
	assert.Equal(t, secretKey, backend.config.SecretKey)
	assert.Equal(t, region, backend.config.Region)
	assert.Equal(t, server, backend.config.Server)
	assert.Equal(t, decryptKey, backend.config.DecryptKey)
	assert.Equal(t, "kubernetes", backend.config.Bucket)
	assert.Equal(t, path, backend.config.Path)
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

	backend, err := NewS3BackendFromParams(params)
	assert.NilError(t, err)

	assert.Equal(t, accessKey, backend.config.AccessKey)
	assert.Equal(t, secretKey, backend.config.SecretKey)
	assert.Equal(t, region, backend.config.Region)
	assert.Equal(t, server, backend.config.Server)
	assert.Equal(t, decryptKey, backend.config.DecryptKey)
	assert.Equal(t, bucket, backend.config.Bucket)
	assert.Equal(t, "kubeconfig/kubeconfig.enc.7z", backend.config.Path)
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

	config := &S3Config{
		AccessKey:  "foo",
		SecretKey:  "foo",
		Server:     "foo",
		DecryptKey: decryptKey,
		Bucket:     bucket,
		Path:       path,
	}
	backend := &S3Backend{
		config:     config,
		Downloader: mockedS3DownloadManager{},
	}

	resultConfigBytesIn, err := backend.Load()
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

	backend, err := NewS3BackendFromParams(params)
	assert.NilError(t, err)

	var expected S3Config
	var result S3Config

	decoderConfig := &mapstructure.DecoderConfig{
		Result:  &expected,
		TagName: "json",
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	assert.NilError(t, err)
	err = decoder.Decode(params)
	assert.NilError(t, err)

	// TODO: unsafe vs. safe test
	resultRaw, err := backend.Config().Yaml(true)
	assert.NilError(t, err)
	err = yaml.Unmarshal(resultRaw, &result)
	assert.NilError(t, err)

	assert.DeepEqual(t, expected, result)
}
