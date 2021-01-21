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
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gabriel-vasile/mimetype"
)

func NewS3BackendFromConfig(config *S3Config) (*S3Backend, error) {
	awsConfig := &aws.Config{
		Region:           aws.String(config.Region),
		Endpoint:         aws.String(config.Server),
		Credentials:      credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, ""),
		S3ForcePathStyle: aws.Bool(true),
	}
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}

	downloader := s3manager.NewDownloader(sess)

	return &S3Backend{
		config:     config,
		Downloader: downloader,
	}, nil
}

func NewS3BackendFromParams(params map[string]interface{}) (*S3Backend, error) {
	config := &S3Config{
		AccessKey:  os.Getenv("S3_ACCESSKEY"),
		SecretKey:  os.Getenv("S3_SECRETKEY"),
		Region:     os.Getenv("S3_REGION"),
		Server:     os.Getenv("S3_SERVER"),
		DecryptKey: os.Getenv("EJSON_PRIVKEY"),
		Bucket:     os.Getenv("S3_BUCKET"),
		Path:       "kubeconfig/kubeconfig.enc.7z",
	}

	// for downward compatibility
	if config.Bucket == "" {
		config.Bucket = "kubernetes"
	}

	if config.Region == "" {
		// minio default region
		config.Region = "us-east-1"
	}

	err := decode(params, &config)
	if err != nil {
		return nil, err
	}

	return NewS3BackendFromConfig(config)
}

func NewS3Backend(accessKey string, secretKey string, region string, server string, decryptKey string, bucket string, path string) (*S3Backend, error) {
	config := &S3Config{
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		Region:     region,
		Server:     server,
		DecryptKey: decryptKey,
		Bucket:     bucket,
		Path:       path,
	}

	return NewS3BackendFromConfig(config)
}

func (b *S3Backend) Load() ([]byte, error) {
	if b.Downloader == nil {
		return nil, fmt.Errorf("no s3 client configured")
	}

	if b.config.Bucket == "" {
		return nil, fmt.Errorf("bucket for the S3 backend is empty")
	}

	if b.config.Path == "" {
		return nil, fmt.Errorf("path for the S3 backend is empty")
	}

	if b.config.AccessKey == "" {
		return nil, fmt.Errorf("AccessKey for the S3 backend is empty")
	}

	if b.config.SecretKey == "" {
		return nil, fmt.Errorf("SecretKey for the S3 backend is empty")
	}

	if b.config.Server == "" {
		return nil, fmt.Errorf("Server for the S3 backend is empty")
	}

	requestInput := s3.GetObjectInput{
		Bucket: aws.String(b.config.Bucket),
		Key:    aws.String(b.config.Path),
	}

	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := b.Downloader.Download(buf, &requestInput)
	if err != nil {
		return nil, err
	}
	data := buf.Bytes()

	mime, err := mimetype.DetectReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to detect mimetype for s3://%s/%s/%s", b.config.Server, b.config.Bucket, b.config.Path)
	}

	var rawKubeconfig []byte

	switch mime.String() {
	case "text/plain":
		rawKubeconfig = data
	case "application/x-7z-compressed":
		rawKubeconfig, err = extractSingleTar7Zip(data, b.config.DecryptKey)
		if err != nil {
			return nil, err
		}
	case "application/octet-stream":
		rawKubeconfig, err = decryptOpensslSymmetric(data, b.config.DecryptKey)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("Unknown kubeconfig source file type: " + mime.String())
	}

	return rawKubeconfig, nil

}

func (b *S3Backend) Type() string {
	return "s3"
}

func (b *S3Backend) Config() BackendConfig {
	return b.config
}

func (c *S3Config) Sanitize() BackendConfig {

	result := &S3Config{
		AccessKey:  c.AccessKey,
		SecretKey:  fmt.Sprintf("%x", sha256.Sum256([]byte(c.SecretKey))),
		Region:     c.Region,
		Server:     c.Server,
		DecryptKey: fmt.Sprintf("%x", sha256.Sum256([]byte(c.DecryptKey))),
		Bucket:     c.Bucket,
		Path:       c.Path,
	}

	return result
}

func (c *S3Config) Yaml(unsafe bool) ([]byte, error) {
	return safeYaml(c, unsafe)
}
