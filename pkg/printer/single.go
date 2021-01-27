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

import "fmt"

func NewSingle(data []map[string]interface{}, field string) *SinglePrinter {
	result := []interface{}{}
	for _, entry := range data {
		result = append(result, entry[field])
	}
	return &SinglePrinter{
		data: result,
	}
}

func (p *SinglePrinter) Print() {
	for _, entry := range p.data {
		fmt.Printf("%+v\n", entry)
	}
}
