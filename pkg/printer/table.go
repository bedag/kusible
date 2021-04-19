/*
Copyright © 2020 Copyright © 2021 Bedag Informatik AG

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
	"os"
	"reflect"

	"github.com/olekukonko/tablewriter"
)

func NewTable(data []map[string]interface{}, header []string) *tablePrinter {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetAutoWrapText(false)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	for _, entry := range data {
		row := []string{}
		for _, fieldName := range header {
			fieldData, ok := entry[fieldName]
			if !ok {
				row = append(row, "")
				continue
			}
			var rowData string
			switch reflect.ValueOf(fieldData).Kind() {
			case reflect.Map, reflect.Struct:
				rowData = "<structured data>"
			case reflect.Slice:
				// Special case: try to render fields that are lists
				// if and only if the table will have only one column.
				//
				// Not sure if this is useful or if it is just "clever" and produces
				// behavior the user is not expecting
				s := reflect.ValueOf(fieldData)
				count := s.Len()
				if count < 1 {
					break
				}
				rowData = "<structured data>"
				if len(header) == 1 {
					// Set the last element of the list in the fieldData
					// as the "normal" row data and then explicitely create additional
					// table rows based on the list in fieldData. the row in rowData
					// will then be added last by the normal row processing below
					rowData = fmt.Sprintf("%+v", s.Index(count-1))
					for i := 0; i < count-2; i++ {
						table.Append([]string{fmt.Sprintf("%+v", s.Index(i))})
					}
				}
			default:
				rowData = fmt.Sprintf("%+v", fieldData)
			}

			row = append(row, rowData)
		}
		table.Append(row)
	}
	return &tablePrinter{
		table: table,
	}
}

func (p *tablePrinter) Print() {
	p.table.Render()
}
