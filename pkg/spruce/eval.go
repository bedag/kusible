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

	"github.com/bedag/kusible/internal/third_party/deinterface"
	"github.com/geofffranks/simpleyaml"
	"github.com/geofffranks/spruce"
	"github.com/mitchellh/mapstructure"
	"github.com/pborman/ansi"
	"sigs.k8s.io/yaml"
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
	// To function, the spruce evaluator expects its data in a very specific
	// structure, which (from what I understand right now), will only be created
	// by https://github.com/geofffranks/simpleyaml/blob/master/simpleyaml.go and
	// the yaml library it uses.
	// To make the evaluator work as expected, we have to convert the input
	// datastructure to the expected format (Marshal with yaml, unmarshal with simpleyaml)
	// eval it and then convert everything back.

	// convert to expected datastructure
	raw, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	y, err := simpleyaml.NewYaml(raw)
	if err != nil {
		return err
	}
	doc, err := y.Map()
	if err != nil {
		return err
	}

	// eval
	evaluator := &spruce.Evaluator{Tree: doc, SkipEval: skipEval}
	err = evaluator.Run(pruneKeys, nil)
	if err != nil {
		return stripAnsiError(err)
	}

	// convert back
	di, err := deinterface.Map(evaluator.Tree, true)

	decoderConfig := &mapstructure.DecoderConfig{ZeroFields: true, Result: data}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return err
	}

	err = decoder.Decode(di)
	if err != nil {
		return err
	}

	return nil
}
