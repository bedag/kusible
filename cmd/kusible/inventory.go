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

	"github.com/bedag/kusible/pkg/inventory"
	"github.com/bedag/kusible/pkg/values"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var inventoryCmd = &cobra.Command{
	Use:   "inventory [regex]",
	Short: "Get inventory information",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filter := args[0]
		inventoryPath := viper.GetString("inventory")
		ejsonPrivKey := viper.GetString("ejson-privkey")
		ejsonKeyDir := viper.GetString("ejson-key-dir")
		skipKubeconfig := viper.GetBool("skip-kubeconfig")

		ejsonSettings := values.EjsonSettings{
			PrivKey:     ejsonPrivKey,
			KeyDir:      ejsonKeyDir,
			SkipDecrypt: false,
		}

		inventory, err := inventory.NewInventory(inventoryPath, ejsonSettings, skipKubeconfig)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to compile inventory.")
			return
		}

		for _, name := range inventory.EntryNames(filter) {
			fmt.Printf("Entry: %s\n", name)
		}
	},
}

func init() {
	inventoryCmd.Flags().BoolP("skip-kubeconfig", "", false, "Skip kubeconfig loading")
	viper.BindPFlag("skip-kubeconfig", inventoryCmd.Flags().Lookup("skip-kubeconfig"))

	rootCmd.AddCommand(inventoryCmd)
}
