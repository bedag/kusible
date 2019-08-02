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
	"bytes"
	"fmt"

	"github.com/mgruener/kusible/pkg/values"
	"github.com/pborman/ansi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/geofffranks/spruce"
	log "github.com/sirupsen/logrus"

	// Use geofffranks yaml library instead of go-yaml
	// to ensure compatibility with spruce
	"github.com/geofffranks/yaml"
)

var valuesCmd = &cobra.Command{
	Use:   "values GROUP ...",
	Short: "Compile values for a list of groups",
	Long: `Use the given groups to compile a single values yaml file.
	The groups are priorized from least to most specific.
	Values of groups of higher priorities override values
	of groups with lower priorities.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		groups := args
		groupVarsDir := viper.GetString("group-vars-dir")
		skipEval := viper.GetBool("skip-eval")
		skipDecrypt := viper.GetBool("skip-decrypt")
		ejsonPrivKey := viper.GetString("ejson-privkey")
		ejsonKeyDir := viper.GetString("ejson-key-dir")

		values, err := values.Compile(groupVarsDir, groups, ejsonKeyDir, ejsonPrivKey, skipEval, skipDecrypt)
		if err != nil {
			// spruce error messages can contain ansi colors
			strippedError, _ := ansi.Strip([]byte(err.Error()))
			log.WithFields(log.Fields{
				"error": string(strippedError),
			}).Fatal("Failed to compile group vars.")
			return
		}

		var result string
		// Although we have a --json option, always marshal to yaml
		// first. The reason is that the type returned by the
		// spruce evaluator cannot be easyly converted to json,
		// but the byte slice returned after marshalling to yaml
		// can be (with the help of spruce).
		// As we do this only in memory to the final document and
		// we probably will not have to deal with huge (> a few MB)
		// datasets, the performance penalty of the double convert
		// should be negligible
		merged, err := yaml.Marshal(values)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
				"yaml":  values,
			}).Fatal("Failed to convert compiled group vars to yaml.")
			return
		}

		if viper.GetBool("json") {
			result, _ = spruce.JSONifyIO(bytes.NewReader(merged), false)
		} else {
			result = string(merged)
		}
		if !viper.GetBool("quiet") {
			fmt.Printf("%s", result)
		}
	},
}

func init() {
	valuesCmd.Flags().BoolP("json", "j", false, "Output json instead of yaml")
	valuesCmd.Flags().BoolP("skip-decrypt", "", false, "Skip ejson decryption")
	viper.BindPFlag("json", valuesCmd.Flags().Lookup("json"))
	viper.BindPFlag("skip-decrypt", valuesCmd.Flags().Lookup("skip-decrypt"))

	rootCmd.AddCommand(valuesCmd)
}
