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
	"archive/tar"
	"bytes"
	"errors"
	"io"

	"github.com/gabriel-vasile/mimetype"
	"github.com/kjk/lzmadec"
)

func extractSingleTar7ZipFile(path string, password string) (bytes.Buffer, error) {
	var archive *lzmadec.Archive
	result := bytes.Buffer{}

	mime, _, err := mimetype.DetectFile(path)
	if err != nil {
		return result, err
	}
	if mime != "application/x-7z-compressed" {
		return result, errors.New("Expected MIME type application/x-7z-compressed but got " + mime)
	}

	if password != "" {
		archive, err = lzmadec.NewEncryptedArchive(path, password)
		if err != nil {
			return result, errors.New("Failed to open archive: " + err.Error())
		}
	} else {
		archive, err = lzmadec.NewArchive(path)
		if err != nil {
			return result, errors.New("Failed to open archive: " + err.Error())
		}
	}

	file := archive.Entries[0].Path
	reader, err := archive.GetFileReader(file)
	if err != nil {
		return result, nil
	}
	defer reader.Close()

	tarReader := tar.NewReader(reader)

	header, err := tarReader.Next()
	if err != nil || header == nil {
		return result, errors.New("Failed to read tar inside the 7zip archive: " + err.Error())
	}

	if header.Typeflag != tar.TypeReg {
		return result, errors.New("The contents of the archive is not a single file")
	}

	if _, err := io.Copy(&result, tarReader); err != nil {
		return result, err
	}

	return result, nil
}
