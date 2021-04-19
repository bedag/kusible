/*
Copyright © 2019 Copyright © 2021 Bedag Informatik AG

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

package target

import (
	"fmt"

	inv "github.com/bedag/kusible/pkg/inventory"
	"github.com/bedag/kusible/pkg/values"
	"github.com/bedag/kusible/pkg/wrapper/ejson"
)

func New(entry *inv.Entry, valuesPath string, skipEval bool, ejson *ejson.Settings) (*Target, error) {
	target := &Target{
		entry: entry,
	}
	groups := entry.Groups()
	values, err := values.New(valuesPath, groups, skipEval, *ejson)
	if err != nil {
		return nil, fmt.Errorf("failed to compile values for target '%s': %s", entry.Name(), err)
	}
	target.values = values
	return target, nil
}

func (t *Target) Values() values.Values {
	return t.values
}

func (t *Target) Entry() *inv.Entry {
	return t.entry
}
