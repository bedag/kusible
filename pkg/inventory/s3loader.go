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
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gabriel-vasile/mimetype"
)

func NewKubeconfigS3LoaderFromParams(params map[string]string) *kubeconfigS3Loader {
	result := map[string]string{
		"accesskey":  os.Getenv("S3_ACCESSKEY"),
		"secretkey":  os.Getenv("S3_SECRETKEY"),
		"region":     os.Getenv("S3_REGION"),
		"server":     os.Getenv("S3_SERVER"),
		"decryptkey": os.Getenv("EJSON_PRIVKEY"),
		"bucket":     "kubernetes",
		"path":       "kubeconfig/kubeconfig.enc.7z",
	}

	for k, v := range params {
		result[strings.ToLower(k)] = v
	}

	return NewKubeconfigS3Loader(
		result["accesskey"],
		result["secretkey"],
		result["region"],
		result["server"],
		result["decryptkey"],
		result["bucket"],
		result["path"])
}

func NewKubeconfigS3Loader(accessKey string, secretKey string, region string, server string, decryptKey string, bucket string, path string) *kubeconfigS3Loader {
	return &kubeconfigS3Loader{
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		Region:     region,
		Server:     server,
		DecryptKey: decryptKey,
		Bucket:     bucket,
		Path:       path,
	}
}

func (loader *kubeconfigS3Loader) Load() ([]byte, error) {
	// TODO: session caching
	awsConfig := &aws.Config{
		Region:           aws.String(loader.Region),
		Endpoint:         aws.String(loader.Server),
		Credentials:      credentials.NewStaticCredentials(loader.AccessKey, loader.SecretKey, ""),
		S3ForcePathStyle: aws.Bool(true),
	}
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}

	downloader := s3manager.NewDownloader(sess)
	requestInput := s3.GetObjectInput{
		Bucket: aws.String(loader.Bucket),
		Key:    aws.String(loader.Path),
	}

	buf := aws.NewWriteAtBuffer([]byte{})
	_, err = downloader.Download(buf, &requestInput)
	if err != nil {
		return nil, err
	}
	data := buf.Bytes()

	mime, _, err := mimetype.DetectReader(bytes.NewReader(data))

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

func (loader *kubeconfigS3Loader) Config() []byte {
	// TODO s3 loader config dump
	var result []byte
	return result
}
