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

	"github.com/bedag/kusible/internal/wrapper/ejson"
	"github.com/bedag/kusible/pkg/inventory"
	"github.com/bedag/kusible/pkg/playbook"
	"github.com/bedag/kusible/pkg/target"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var playbookCmd = &cobra.Command{
	Use:   "playbook [file]",
	Short: "Render the given playbook",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		playbookFile := args[0]
		limits := viper.GetStringSlice("limit")
		groupVarsDir := viper.GetString("group-vars-dir")
		inventoryPath := viper.GetString("inventory")
		skipDecrypt := viper.GetBool("skip-decrypt")
		ejsonPrivKey := viper.GetString("ejson-privkey")
		ejsonKeyDir := viper.GetString("ejson-key-dir")

		invEjsonSettings := ejson.Settings{
			PrivKey: ejsonPrivKey,
			KeyDir:  ejsonKeyDir,
			// if we want to retrieve the cluster inventory ConfigMap
			// we need a kubeconfig to retrieve it, so we cannot skip
			// the decryption of the inventory settings
			SkipDecrypt: false,
		}

		tgtEjsonSettings := ejson.Settings{
			PrivKey:     ejsonPrivKey,
			KeyDir:      ejsonKeyDir,
			SkipDecrypt: skipDecrypt,
		}

		// if we do not retrieve the cluster inventory ConfigMap, we do not need to retrieve
		// the kubeconfig
		inventory, err := inventory.NewInventory(inventoryPath, invEjsonSettings, false)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to compile inventory.")
		}

		targets, err := target.NewTargets(".*", limits, groupVarsDir, inventory, true, &tgtEjsonSettings)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to compile values for inventory entries.")
		}

		// TODO: add parameter to skip eval and skip cluster inventory
		playbookSet, err := playbook.NewSet(playbookFile, targets, false, false)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to compile playbooks.")
		}

		for name, playbook := range playbookSet {
			if len(playbook.Config.Plays) > 0 {
				result, err := playbook.Config.YAML()
				if err != nil {
					log.WithFields(log.Fields{
						"entry": name,
						"error": err.Error(),
					}).Fatal("Failed to convert playbook entry to yaml.")
				}
				fmt.Printf("======= Plays for %s =======\n", name)
				fmt.Printf("%s", string(result))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(playbookCmd)
}
