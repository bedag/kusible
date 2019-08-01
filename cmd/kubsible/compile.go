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
	"fmt"

	"github.com/mgruener/kusible/pkg/groupvars"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"

	// Use geofffranks yaml library instead of go-yaml
	// to ensure compatibility with spruce
	"github.com/geofffranks/yaml"
)

var optGroupVarsDir string
var optQuiet bool

func init() {
	compileCmd.Flags().StringVarP(&optGroupVarsDir, "dir", "d", "group_vars", "Source directory to read from")
	compileCmd.Flags().BoolVarP(&optQuiet, "quiet", "q", false, "Suppress all normal output")
	rootCmd.AddCommand(compileCmd)
}

var compileCmd = &cobra.Command{
	Use:   "compile GROUP ...",
	Short: "Compile the values for the given groups",
	Long: `Use the given groups to compile a single yaml file.
	The groups are priorized from least to most specific.
	Values of groups of higher priorities override values
	of groups with lower priorities.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		values, err := groupvars.Compile(optGroupVarsDir, args)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to compile group vars.")
			return
		}

		merged, err := yaml.Marshal(values)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
				"yaml":  values,
			}).Fatal("Failed to convert compiled group vars to yaml.")
			return
		}

		if !optQuiet {
			fmt.Printf("%s", string(merged))
		}
	},
}
