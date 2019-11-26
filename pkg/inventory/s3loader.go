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
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gabriel-vasile/mimetype"
	"gopkg.in/yaml.v2"
)

func NewKubeconfigS3LoaderFromParams(params map[string]string) *kubeconfigS3Loader {
	result := map[string]string{
		"accesskey":   os.Getenv("S3_ACCESSKEY"),
		"secretkey":   os.Getenv("S3_SECRETKEY"),
		"region":      os.Getenv("S3_REGION"),
		"server":      os.Getenv("S3_SERVER"),
		"decrypt_key": os.Getenv("EJSON_PRIVKEY"),
		"bucket":      os.Getenv("S3_BUCKET"),
		"path":        "kubeconfig/kubeconfig.enc.7z",
	}

	// for downward compatibility
	if result["bucket"] == "" {
		result["bucket"] = "kubernetes"
	}

	for k, v := range params {
		result[strings.ToLower(k)] = v
	}

	return NewKubeconfigS3Loader(
		result["accesskey"],
		result["secretkey"],
		result["region"],
		result["server"],
		result["decrypt_key"],
		result["bucket"],
		result["path"])
}

func NewKubeconfigS3Loader(accessKey string, secretKey string, region string, server string, decryptKey string, bucket string, path string) *kubeconfigS3Loader {
	awsConfig := &aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(server),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		S3ForcePathStyle: aws.Bool(true),
	}
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil
	}

	downloader := s3manager.NewDownloader(sess)

	return &kubeconfigS3Loader{
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		Region:     region,
		Server:     server,
		DecryptKey: decryptKey,
		Bucket:     bucket,
		Path:       path,
		Downloader: downloader,
	}
}

func (loader *kubeconfigS3Loader) Load() ([]byte, error) {
	if loader.Downloader == nil {
		return nil, fmt.Errorf("no s3 client configured")
	}

	requestInput := s3.GetObjectInput{
		Bucket: aws.String(loader.Bucket),
		Key:    aws.String(loader.Path),
	}

	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := loader.Downloader.Download(buf, &requestInput)
	if err != nil {
		return nil, err
	}
	data := buf.Bytes()

	mime, _, err := mimetype.DetectReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to detect mimetype for s3://%s/%s/%s", loader.Server, loader.Bucket, loader.Path)
	}

	var rawKubeconfig []byte

	switch mime {
	case "text/plain":
		rawKubeconfig = data
	case "application/x-7z-compressed":
		rawKubeconfig, err = extractSingleTar7Zip(data, loader.DecryptKey)
		if err != nil {
			return nil, err
		}
	case "application/octet-stream":
		rawKubeconfig, err = decryptOpensslSymmetric(data, loader.DecryptKey)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("Unknown kubeconfig source file type: " + mime)
	}

	return rawKubeconfig, nil

}

func (loader *kubeconfigS3Loader) Type() string {
	return "s3"
}

func (loader *kubeconfigS3Loader) Config() ([]byte, error) {
	config := map[string]string{
		"accesskey":   loader.AccessKey,
		"secretkey":   loader.SecretKey,
		"region":      loader.Region,
		"server":      loader.Server,
		"decrypt_key": loader.DecryptKey,
		"bucket":      loader.Bucket,
		"path":        loader.Path,
	}
	result, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}
	return result, nil
}
