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

package printer

import (
	"fmt"
	"reflect"
)

func NewSingle(data []map[string]interface{}, field string) *singlePrinter {
	result := []interface{}{}
	for _, entry := range data {
		result = append(result, entry[field])
	}
	return &singlePrinter{
		data: result,
	}
}

func (p *singlePrinter) Print() {
	for _, entry := range p.data {
		switch reflect.TypeOf(entry).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(entry)

			for i := 0; i < s.Len(); i++ {
				fmt.Printf("%+v\n", s.Index(i))
			}
		default:
			fmt.Printf("%+v\n", entry)
		}
	}
}
