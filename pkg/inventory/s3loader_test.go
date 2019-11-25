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
	"io"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"gotest.tools/assert"
)

type mockedS3DownloadManager struct {
	s3manageriface.DownloaderAPI
	PartSize       int64
	Concurrency    int
	S3             s3iface.S3API
	RequestOptions []request.Option
}

func (s mockedS3DownloadManager) Download(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error) {
	return 0, nil
}

func (s mockedS3DownloadManager) DownloadWithContext(ctx aws.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error) {
	return 0, nil
}

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

	params := map[string]string{
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

	params := map[string]string{
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

	params := map[string]string{}

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

// TODO: S3Loader Load() test
//func TestS3LoaderLoad(t *testing.T) {
//	decryptKey := "test123"
//	bucket := "testdata"
//	path := "kubeconfig.enc"
//
//	loader := &kubeconfigS3Loader{
//		DecryptKey: decryptKey,
//		Bucket:     bucket,
//		Path:       path,
//		Downloader: mockedS3DownloadManager{},
//	}
//}
