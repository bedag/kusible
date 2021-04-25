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

package values

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/bedag/kusible/pkg/wrapper/ejson"
	log "github.com/sirupsen/logrus"
	"gotest.tools/assert"
	"sigs.k8s.io/yaml"
)

func TestDirectory(t *testing.T) {
	// disable logs during this test as it would spam warnings because
	// it cannot decrypt the (for this test) unencrypted ejson files
	log.SetOutput(ioutil.Discard)

	loadYaml := func(path string) (map[string]interface{}, error) {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		var result map[string]interface{}
		err = yaml.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	tests := map[string]struct {
		groupVarsDir string
		expectedFile string
		groups       []string
	}{
		"multi-dir-multi-file-123":         {groupVarsDir: "multi-dir-multi-file", expectedFile: "multi-dir-multi-file-123.yml", groups: []string{"multi-file-01", "multi-file-02", "multi-file-03"}},
		"multi-dir-multi-file-213":         {groupVarsDir: "multi-dir-multi-file", expectedFile: "multi-dir-multi-file-213.yml", groups: []string{"multi-file-02", "multi-file-01", "multi-file-03"}},
		"multi-dir-multi-file-312":         {groupVarsDir: "multi-dir-multi-file", expectedFile: "multi-dir-multi-file-312.yml", groups: []string{"multi-file-03", "multi-file-01", "multi-file-02"}},
		"multi-dir-multi-file-132":         {groupVarsDir: "multi-dir-multi-file", expectedFile: "multi-dir-multi-file-132.yml", groups: []string{"multi-file-01", "multi-file-03", "multi-file-02"}},
		"multi-dir-multi-file-231":         {groupVarsDir: "multi-dir-multi-file", expectedFile: "multi-dir-multi-file-231.yml", groups: []string{"multi-file-02", "multi-file-03", "multi-file-01"}},
		"multi-dir-multi-file-321":         {groupVarsDir: "multi-dir-multi-file", expectedFile: "multi-dir-multi-file-321.yml", groups: []string{"multi-file-03", "multi-file-02", "multi-file-01"}},
		"multi-dir-single-file-123":        {groupVarsDir: "multi-dir-single-file", expectedFile: "multi-dir-single-file-123.yml", groups: []string{"single-file-01", "single-file-02", "single-file-03"}},
		"multi-dir-single-file-213":        {groupVarsDir: "multi-dir-single-file", expectedFile: "multi-dir-single-file-213.yml", groups: []string{"single-file-02", "single-file-01", "single-file-03"}},
		"multi-dir-single-file-312":        {groupVarsDir: "multi-dir-single-file", expectedFile: "multi-dir-single-file-312.yml", groups: []string{"single-file-03", "single-file-01", "single-file-02"}},
		"multi-dir-single-file-132":        {groupVarsDir: "multi-dir-single-file", expectedFile: "multi-dir-single-file-132.yml", groups: []string{"single-file-01", "single-file-03", "single-file-02"}},
		"multi-dir-single-file-231":        {groupVarsDir: "multi-dir-single-file", expectedFile: "multi-dir-single-file-231.yml", groups: []string{"single-file-02", "single-file-03", "single-file-01"}},
		"multi-dir-single-file-321":        {groupVarsDir: "multi-dir-single-file", expectedFile: "multi-dir-single-file-321.yml", groups: []string{"single-file-03", "single-file-02", "single-file-01"}},
		"multi-file-123":                   {groupVarsDir: "multi-file", expectedFile: "multi-file-123.yml", groups: []string{"file-01", "file-02", "file-03"}},
		"multi-file-213":                   {groupVarsDir: "multi-file", expectedFile: "multi-file-213.yml", groups: []string{"file-02", "file-01", "file-03"}},
		"multi-file-312":                   {groupVarsDir: "multi-file", expectedFile: "multi-file-312.yml", groups: []string{"file-03", "file-01", "file-02"}},
		"multi-file-132":                   {groupVarsDir: "multi-file", expectedFile: "multi-file-132.yml", groups: []string{"file-01", "file-03", "file-02"}},
		"multi-file-231":                   {groupVarsDir: "multi-file", expectedFile: "multi-file-231.yml", groups: []string{"file-02", "file-03", "file-01"}},
		"multi-file-321":                   {groupVarsDir: "multi-file", expectedFile: "multi-file-321.yml", groups: []string{"file-03", "file-02", "file-01"}},
		"multi-mixed-123":                  {groupVarsDir: "multi-mixed", expectedFile: "multi-mixed-123.yml", groups: []string{"file-01", "file-02", "file-03"}},
		"multi-mixed-213":                  {groupVarsDir: "multi-mixed", expectedFile: "multi-mixed-213.yml", groups: []string{"file-02", "file-01", "file-03"}},
		"multi-mixed-312":                  {groupVarsDir: "multi-mixed", expectedFile: "multi-mixed-312.yml", groups: []string{"file-03", "file-01", "file-02"}},
		"multi-mixed-132":                  {groupVarsDir: "multi-mixed", expectedFile: "multi-mixed-132.yml", groups: []string{"file-01", "file-03", "file-02"}},
		"multi-mixed-231":                  {groupVarsDir: "multi-mixed", expectedFile: "multi-mixed-231.yml", groups: []string{"file-02", "file-03", "file-01"}},
		"multi-mixed-321":                  {groupVarsDir: "multi-mixed", expectedFile: "multi-mixed-321.yml", groups: []string{"file-03", "file-02", "file-01"}},
		"multi-mixed-dirfile-123":          {groupVarsDir: "multi-mixed-dirfile", expectedFile: "multi-mixed-dirfile-123.yml", groups: []string{"single-mixed-01", "single-mixed-02", "single-mixed-03"}},
		"multi-mixed-dirfile-213":          {groupVarsDir: "multi-mixed-dirfile", expectedFile: "multi-mixed-dirfile-213.yml", groups: []string{"single-mixed-02", "single-mixed-01", "single-mixed-03"}},
		"multi-mixed-dirfile-312":          {groupVarsDir: "multi-mixed-dirfile", expectedFile: "multi-mixed-dirfile-312.yml", groups: []string{"single-mixed-03", "single-mixed-01", "single-mixed-02"}},
		"multi-mixed-dirfile-132":          {groupVarsDir: "multi-mixed-dirfile", expectedFile: "multi-mixed-dirfile-132.yml", groups: []string{"single-mixed-01", "single-mixed-03", "single-mixed-02"}},
		"multi-mixed-dirfile-231":          {groupVarsDir: "multi-mixed-dirfile", expectedFile: "multi-mixed-dirfile-231.yml", groups: []string{"single-mixed-02", "single-mixed-03", "single-mixed-01"}},
		"multi-mixed-dirfile-321":          {groupVarsDir: "multi-mixed-dirfile", expectedFile: "multi-mixed-dirfile-321.yml", groups: []string{"single-mixed-03", "single-mixed-02", "single-mixed-01"}},
		"single-dir-multi-dir-multi-file":  {groupVarsDir: "single-dir-multi-dir-multi-file", expectedFile: "single-dir-multi-dir-multi-file.yml", groups: []string{"multi-dir-multi-file"}},
		"single-dir-multi-dir-single-file": {groupVarsDir: "single-dir-multi-dir-single-file", expectedFile: "single-dir-multi-dir-single-file.yml", groups: []string{"multi-dir-single-file"}},
		"single-dir-multi-file":            {groupVarsDir: "single-dir-multi-file", expectedFile: "single-dir-multi-file.yml", groups: []string{"multi-file"}},
		"single-dir-multi-mixed-dirfile":   {groupVarsDir: "single-dir-multi-mixed-dirfile", expectedFile: "single-dir-multi-mixed-dirfile.yml", groups: []string{"multi-mixed-dirfile"}},
		"single-dir-multi-mixed":           {groupVarsDir: "single-dir-multi-mixed", expectedFile: "single-dir-multi-mixed.yml", groups: []string{"multi-mixed"}},
		"single-dir-single-file":           {groupVarsDir: "single-dir-single-file", expectedFile: "single-dir-single-file.yml", groups: []string{"single-file"}},
		"single-dir-single-mixed-dirfile":  {groupVarsDir: "single-dir-single-mixed-dirfile", expectedFile: "single-dir-single-mixed-dirfile.yml", groups: []string{"single-mixed-dirfile"}},
		"single-dir-single-mixed":          {groupVarsDir: "single-dir-single-mixed", expectedFile: "single-dir-single-mixed.yml", groups: []string{"single-mixed"}},
		"single-file":                      {groupVarsDir: "single-file", expectedFile: "single-file.yml", groups: []string{"file"}},
		"single-mixed-dirfile":             {groupVarsDir: "single-mixed-dirfile", expectedFile: "single-mixed-dirfile.yml", groups: []string{"single-mixed"}},
		"single-mixed":                     {groupVarsDir: "single-mixed", expectedFile: "single-mixed.yml", groups: []string{"file"}},
	}

	marshalMethods := []string{"JSON", "YAML"}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			d, err := NewDirectory("testdata/directory/"+tc.groupVarsDir, tc.groups, true, ejson.Settings{})
			assert.NilError(t, err)
			got := d.Map()
			assert.NilError(t, err)

			want, err := loadYaml("testdata/directory/expected/" + tc.expectedFile)
			assert.NilError(t, err)

			assert.DeepEqual(t, want, got)

			for _, method := range marshalMethods {
				r := reflect.ValueOf(d).MethodByName(method).Call([]reflect.Value{})
				resultBytes := r[0].Bytes()
				// cannot use assert.NilError here because we cannot
				// cast nil to error
				assert.Assert(t, r[1].Interface() == nil)

				err = yaml.Unmarshal(resultBytes, map[string]interface{}{})
				assert.NilError(t, err)
			}
		})
	}
}
