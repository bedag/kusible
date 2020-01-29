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

/*
Package config implements the playbook config format
*/
package config

import (
	"fmt"
	"regexp"
)

// A Pattern is based on the idea of patterns in Ansible:
// https://docs.ansible.com/ansible/latest/user_guide/intro_patterns.html
// Each pattern is regexp that is implicitely enclosed in ^$ and has
// an optional modifier prefix.
type Pattern struct {
	modifier string // &, ! or ""
	regex    *regexp.Regexp
	groups   []string
}

// Validator is used to determin if a list of Patterns added via the
// Add() method is valid.
type Validator struct {
	allOf []*Pattern
	anyOf []*Pattern
}

// NewPattern returns a new pattern based on the given expression string
// and a list of groups. If the first charactor of the expression
// is either ! or &, it will be treated as modifier and the rest is
// treated as a regexp that will be wrapped in ^$. If
// no modifier is present, the whole expression will be treated as a
// ^$ wrapped regexp. The regexp is then matched against the list
// of groups and all matching groups will be stored in its internal
// groups field wich will be used by the Matches() method to
// determinee if the pattern matches the given groups.
func NewPattern(expr string, groups []string) (*Pattern, error) {
	if len(expr) < 1 {
		return nil, fmt.Errorf("empty pattern expression")
	}
	modifier := ""
	value := ""
	if (string(expr[0]) == "&") || (string(expr[0]) == "!") {
		modifier = string(expr[0])
		value = expr[1:]
	} else {
		value = expr
	}

	if len(value) < 1 {
		return nil, fmt.Errorf("only modifier given in expression")
	}

	regex, err := regexp.Compile("^" + value + "$")
	if err != nil {
		return nil, err
	}
	pattern := &Pattern{
		modifier: modifier,
		regex:    regex,
		groups:   []string{},
	}

	for _, group := range groups {
		if pattern.regex.MatchString(group) {
			pattern.groups = append(pattern.groups, group)
		}
	}
	return pattern, nil
}

// Matches returns true if the modifier is !
// and it has no matching groups associated.
// In all other cases Matches is true if the
// list of matching groups is non-empty.
func (p *Pattern) Matches() bool {
	valid := false
	if len(p.groups) > 0 {
		valid = true
	}
	if p.modifier == "!" {
		return !valid
	}
	return valid
}

// Groups returns the list of groups matched by
// the regexp of this pattern
func (p *Pattern) Groups() []string {
	return p.groups
}

// Add adds a given pattern either to the internal
// "all" or "any" list, based on the pattern modifier.
// Patterns with "!" and "&" modifier are added to the
// "all" list and all other Patterns will be added to the
// "any" list. Refer to the Valid() method for details.
func (v *Validator) Add(pattern *Pattern) {
	if pattern.modifier == "&" || pattern.modifier == "!" {
		v.allOf = append(v.allOf, pattern)
	} else {
		v.anyOf = append(v.anyOf, pattern)
	}
}

// Valid returns true if all Patterns in the
// "all" list match and at least one of the
// "any" list
func (v *Validator) Valid() bool {
	// No pattern present. Nothing to match against
	// is equivalent to nothing matches at all
	if (len(v.allOf) < 1) && (len(v.anyOf) < 1) {
		return false
	}

	for _, pattern := range v.allOf {
		if !pattern.Matches() {
			return false
		}
	}

	if len(v.anyOf) < 1 {
		return true
	}

	for _, pattern := range v.anyOf {
		if pattern.Matches() {
			return true
		}
	}

	return false
}
