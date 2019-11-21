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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gabriel-vasile/mimetype"
)

func NewKubeconfigS3Loader(accessKey string, secretKey string, region string, server string, decryptKey string, bucket string, path string) *kubeconfigS3Loader {
	return &kubeconfigS3Loader{
		accessKey:  accessKey,
		secretKey:  secretKey,
		region:     region,
		server:     server,
		decryptKey: decryptKey,
		bucket:     bucket,
		path:       path,
	}
}

func (loader *kubeconfigS3Loader) Load() ([]byte, error) {
	// TODO: session caching
	awsConfig := &aws.Config{
		Region:           aws.String(loader.region),
		Endpoint:         aws.String(loader.server),
		Credentials:      credentials.NewStaticCredentials(loader.accessKey, loader.secretKey, ""),
		S3ForcePathStyle: aws.Bool(true),
	}
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}

	downloader := s3manager.NewDownloader(sess)
	requestInput := s3.GetObjectInput{
		Bucket: aws.String(loader.bucket),
		Key:    aws.String(loader.path),
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
		rawKubeconfig, err = extractSingleTar7Zip(data, loader.decryptKey)
		if err != nil {
			return nil, err
		}
	case "application/octet-stream":
		rawKubeconfig, err = decryptOpensslSymmetric(data, loader.decryptKey)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("Unknown kubeconfig source file type: " + mime)
	}

	return rawKubeconfig, nil

}
