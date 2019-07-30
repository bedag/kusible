// Copyright © 2019 Michael Gruener
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/mgruener/kusible/pkg/groupvars"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(compileCmd)
}

var compileCmd = &cobra.Command{
	Use:   "compile <group> <group> ...",
	Short: "Compile the values for the given groups",
	Long: `Use the given groups to compile a single yaml file.
	The groups are priorized from least to most specific.
	Values of groups of higher priorities override values
	of groups with lower priorities.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		values, err := groupvars.Compile(args)
	},
}
