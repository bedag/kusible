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

// NewJob returns a struct that implements the Printable interface.
// It expects a function that generates the expected printer data
// given a list of fields to be printed
func NewJob(fn DataFn) Printable {
	return job{
		dataFn: fn,
	}
}

func (j job) PrinterData(fields []string) map[string]interface{} {
	return j.dataFn(fields)
}
