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

package helm

import (
	"helm.sh/helm/v3/pkg/releaseutil"
)

func SplitSortManifest(bigManifest string) ([]string, error) {
	input := map[string]string{
		"kusible": bigManifest,
	}

	_, manifests, err := releaseutil.SortManifests(input, nil, releaseutil.InstallOrder)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, manifest := range manifests {
		result = append(result, manifest.Content)
	}

	return result, nil
}
