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

package config

/*
Each run-relevant inventory entry has its own "view" on the given playbook containting
only the relevant plays for its groups.

Given a list of groups, the playbook loader

* loads the playbook (without evaluation)
* filters the plays based on the given groups
* loads the values relevant for the given groups (without evaluation)
* merges the filtered playbook and values
* evaluates the result
* unmarshalls the merged/evaluated playbook/value map into a valid playbook config structure
*/
