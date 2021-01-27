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
	"encoding/json"
	"fmt"
)

func NewJSON(data []map[string]interface{}) *JSONPrinter {
	return &JSONPrinter{
		data: data,
	}
}

func (p *JSONPrinter) Print() {
	items := map[string]interface{}{
		"items": p.data,
	}
	result, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		fmt.Printf("Failed to print data as json: %s", err)
	}
	fmt.Printf("%s\n", string(result))
}
