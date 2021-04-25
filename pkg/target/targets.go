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

package target

import (
	"fmt"

	inv "github.com/bedag/kusible/pkg/inventory"
	"github.com/bedag/kusible/pkg/wrapper/ejson"
)

func NewTargets(filter string, limits []string, valuesPath string, inventory *inv.Inventory, skipEval bool, ejson *ejson.Settings) (*Targets, error) {
	targetNames, err := inventory.EntryNames(filter, limits)
	if err != nil {
		return nil, fmt.Errorf("failed to get possible entries from inventory: %s", err)
	}

	targets := &Targets{
		limits:     limits,
		filter:     filter,
		valuesPath: valuesPath,
		targets:    make(map[string]*Target, len(targetNames)),
	}
	if len(targetNames) <= 0 {
		return targets, nil
	}

	for _, name := range targetNames {
		entry := inventory.Entries()[name]
		target, err := New(entry, valuesPath, skipEval, ejson)
		if err != nil {
			return nil, fmt.Errorf("failed to create target for inventory entry '%s': %s", name, err)
		}
		targets.targets[name] = target
	}
	return targets, nil
}

func (t *Targets) Names() []string {
	result := []string{}

	for name := range t.targets {
		result = append(result, name)
	}
	return result
}

func (t *Targets) Targets() map[string]*Target {
	return t.targets
}

func (t *Targets) Limits() []string {
	return t.limits
}

func (t *Targets) ValuesPath() string {
	return t.valuesPath
}

func (t *Targets) EJSON() *ejson.Settings {
	return t.ejson
}

func (t *Targets) Filter() string {
	return t.filter
}
