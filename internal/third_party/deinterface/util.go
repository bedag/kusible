/*
The MIT License (MIT)

Copyright (c) 2016 Geoff Franks

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

Based on https://github.com/geofffranks/spruce/blob/a17d090c935bb2252af1bf84c1220f24e791e9c3/json.go
*/

package deinterface

import "fmt"

func Map(o interface{}, strict bool) (interface{}, error) {
	switch o.(type) {
	case map[interface{}]interface{}:
		return deinterfaceMap(o.(map[interface{}]interface{}), strict)
	case []interface{}:
		return deinterfaceList(o.([]interface{}), strict)
	default:
		return o, nil
	}
}

func addKeyToMap(m map[string]interface{}, k interface{}, v interface{}, strict bool) error {
	vs := fmt.Sprintf("%v", k)
	_, exists := m[vs]
	if exists {
		//NewWarningError(eContextAll, "@Y{Duplicate key detected: %s}", vs).Warn()
		return nil
	}
	dv, err := Map(v, strict)
	if err != nil {
		return err
	}
	m[vs] = dv
	return nil
}

func deinterfaceMap(o map[interface{}]interface{}, strict bool) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	for k, v := range o {

		switch k.(type) {
		case string:
			err := addKeyToMap(m, k, v, strict)
			if err != nil {
				return nil, err
			}
		default:
			if strict {
				return nil, fmt.Errorf("Non-string keys found during strict JSON conversion")
			} else {
				addKeyToMap(m, k, v, strict)
			}
		}

	}
	return m, nil
}

func deinterfaceList(o []interface{}, strict bool) ([]interface{}, error) {
	l := make([]interface{}, len(o))
	for i, v := range o {
		v_, err := Map(v, strict)
		if err != nil {
			return nil, err
		}
		l[i] = v_
	}
	return l, nil
}
