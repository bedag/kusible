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
	"os"
	"reflect"

	"github.com/olekukonko/tablewriter"
)

func NewTable(data []map[string]interface{}, header []string) *TablePrinter {
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
			fieldType := reflect.ValueOf(fieldData).Kind()
			switch fieldType {
			case reflect.Map, reflect.Struct:
				rowData = "<Multiline data will not be rendered in table mode>"
			default:
				rowData = fmt.Sprintf("%+v", fieldData)
			}

			row = append(row, rowData)
		}
		table.Append(row)
	}
	return &TablePrinter{
		table: table,
	}
}

func (p *TablePrinter) Print() {
	p.table.Render()
}
