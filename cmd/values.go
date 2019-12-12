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

package cmd

import (
	"fmt"

	"github.com/bedag/kusible/pkg/values"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
	// Use geofffranks yaml library instead of go-yaml
	// to ensure compatibility with spruce
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

		ejsonSettings := values.EjsonSettings{
			PrivKey:     ejsonPrivKey,
			KeyDir:      ejsonKeyDir,
			SkipDecrypt: skipDecrypt,
		}

		values, err := values.New(groupVarsDir, groups, skipEval, ejsonSettings)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to compile group vars.")
			return
		}

		var result []byte

		if viper.GetBool("json") {
			result, err = values.JSON()
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Fatal("Failed to convert compiled group vars to json.")
			}
		} else {
			result, err = values.YAML()
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Fatal("Failed to convert compiled group vars to yaml.")
			}
		}
		if !viper.GetBool("quiet") {
			fmt.Printf("%s", string(result))
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
