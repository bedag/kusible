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
	"strings"
)

func New(format Format, fields []string, queue Queue, options Options) (Printer, error) {
	feed := []map[string]interface{}{}
	for _, object := range queue {
		feed = append(feed, object.PrinterData(fields))
	}

	var printer Printer
	switch format {
	case FormatJSON:
		printer = NewJSON(feed, options)
	case FormatYAML:
		printer = NewYAML(feed, options)
	case FormatTable:
		if len(fields) < 1 {
			return nil, fmt.Errorf("'table' printer requires at least one field to print")
		}
		printer = NewTable(feed, fields)
	case FormatSingle:
		if len(fields) != 1 {
			return nil, fmt.Errorf("'single' printer requires exactly one field to print")
		}
		printer = NewSingle(feed, fields[0])
	default:
		return nil, fmt.Errorf("unknown printer format")
	}
	return printer, nil
}

func ParseFormat(format string) (Format, error) {
	switch strings.ToLower(format) {
	case "json":
		return FormatJSON, nil
	case "yaml":
		return FormatYAML, nil
	case "table":
		return FormatTable, nil
	case "single":
		return FormatSingle, nil
	default:
		return InvalidFormat, fmt.Errorf("invalid format: %s", format)
	}
}
