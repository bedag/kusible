/*
Copyright © 2021 Bedag Informatik AG

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
	"archive/tar"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"

	openssl "github.com/Luzifer/go-openssl/v3"
	"github.com/gabriel-vasile/mimetype"
	"github.com/kjk/lzmadec"
	"github.com/mitchellh/mapstructure"
)

func extractSingleTar7Zip(data []byte, password string) ([]byte, error) {
	// extracting 7zip data only works with files stored in the filesystem
	tmpfile, err := ioutil.TempFile("", "s3loader")
	if err != nil {
		return nil, err
	}
	defer func(tmpfile *os.File) {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}(tmpfile)

	if _, err := tmpfile.Write(data); err != nil {
		return nil, err
	}
	result, err := extractSingleTar7ZipFile(tmpfile.Name(), password)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func extractSingleTar7ZipFile(path string, password string) ([]byte, error) {
	var archive *lzmadec.Archive

	mime, err := mimetype.DetectFile(path)
	if err != nil {
		return nil, err
	}
	if !mime.Is("application/x-7z-compressed") {
		return nil, errors.New("expected MIME type application/x-7z-compressed but got " + mime.String())
	}

	if password != "" {
		archive, err = lzmadec.NewEncryptedArchive(path, password)
		if err != nil {
			return nil, errors.New("failed to open archive: " + err.Error())
		}
	} else {
		archive, err = lzmadec.NewArchive(path)
		if err != nil {
			return nil, errors.New("failed to open archive: " + err.Error())
		}
	}

	file := archive.Entries[0].Path
	reader, err := archive.GetFileReader(file)
	if err != nil {
		return nil, nil
	}
	defer reader.Close()

	tarReader := tar.NewReader(reader)

	header, err := tarReader.Next()
	if err != nil || header == nil {
		return nil, errors.New("failed to read tar inside the 7zip archive: " + err.Error())
	}

	if header.Typeflag != tar.TypeReg {
		return nil, errors.New("the contents of the archive is not a single file")
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, tarReader); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decryptOpensslSymmetricFile(path string, password string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	result, err := decryptOpensslSymmetric(data, password)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func decryptOpensslSymmetric(data []byte, password string) ([]byte, error) {
	o := openssl.New()
	result, err := o.DecryptBinaryBytes(password, data, openssl.DigestSHA256Sum)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func decode(input interface{}, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		TagName: "json",
		Result:  output,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	if err := decoder.Decode(input); err != nil {
		return err
	}
	return nil
}
