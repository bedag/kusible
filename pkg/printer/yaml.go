/*
Copyright © 2020 Michael Gruener

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

package printer

import (
	"fmt"

	"sigs.k8s.io/yaml"
)

func NewYAML(data []map[string]interface{}, options Options) *YAMLPrinter {
	return &YAMLPrinter{
		data:               data,
		listWrapSingleItem: options.ListWrapSingleItem,
	}
}

func (p *YAMLPrinter) Print() {
	var items interface{}
	items = map[string]interface{}{
		"items": p.data,
	}
	if len(p.data) == 1 && !p.listWrapSingleItem {
		items = p.data[0]
	}
	result, err := yaml.Marshal(items)
	if err != nil {
		fmt.Printf("Failed to print data as json: %s", err)
	}
	fmt.Printf("%s\n", string(result))
}
