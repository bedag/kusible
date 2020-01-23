/*
Copyright © 2019 Michael Gruener

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

package spruce

import (
	"errors"

	"github.com/geofffranks/spruce"
	"github.com/mitchellh/mapstructure"
	"github.com/pborman/ansi"
)

func stripAnsiError(err error) error {
	if err != nil {
		strippedError, _ := ansi.Strip([]byte(err.Error()))
		return errors.New(string(strippedError))
	}
	return nil
}

// Eval is a wrapper around the Evaluator of https://github.com/geofffranks/spruce
// that handles the necessary type conversion
func Eval(data *map[string]interface{}, skipEval bool, pruneKeys []string) error {
	var spruceMap map[interface{}]interface{}

	err := mapstructure.Decode(data, &spruceMap)
	if err != nil {
		return err
	}

	evaluator := &spruce.Evaluator{Tree: spruceMap, SkipEval: skipEval}
	err = evaluator.Run(pruneKeys, nil)
	if err != nil {
		return stripAnsiError(err)
	}

	decoderConfig := &mapstructure.DecoderConfig{ZeroFields: true, Result: data}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return err
	}

	err = decoder.Decode(evaluator.Tree)
	if err != nil {
		return err
	}
	return nil
}
