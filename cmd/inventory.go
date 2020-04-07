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

	"github.com/bedag/kusible/pkg/wrapper/ejson"
	"github.com/bedag/kusible/pkg/inventory"
	"github.com/bedag/kusible/pkg/target"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var inventoryCmd = &cobra.Command{
	Use:   "inventory",
	Short: "Get inventory information",
}

var inventoryListCmd = &cobra.Command{
	Use:   "list [regex]",
	Short: "List inventory entries",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filter := args[0]
		limits := viper.GetStringSlice("limit")
		inventoryPath := viper.GetString("inventory")
		ejsonPrivKey := viper.GetString("ejson-privkey")
		ejsonKeyDir := viper.GetString("ejson-key-dir")
		ejsonSettings := ejson.Settings{
			PrivKey:     ejsonPrivKey,
			KeyDir:      ejsonKeyDir,
			SkipDecrypt: false,
		}

		// as we just want to list the available inventory entries we can (and should)
		// skip kubeconfig retrieval
		skipKubeconfig := true
		inventory, err := inventory.NewInventory(inventoryPath, ejsonSettings, skipKubeconfig)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to compile inventory.")
			return
		}

		names, err := inventory.EntryNames(filter, limits)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to get list of entries")
		}
		for _, name := range names {
			fmt.Printf("%s\n", name)
		}
	},
}

var inventoryKubeconfigCmd = &cobra.Command{
	Use:   "kubeconfig [entry name]",
	Short: "Get the kubeconfig for a given inventory entry",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		inventoryPath := viper.GetString("inventory")
		ejsonPrivKey := viper.GetString("ejson-privkey")
		ejsonKeyDir := viper.GetString("ejson-key-dir")

		ejsonSettings := ejson.Settings{
			PrivKey:     ejsonPrivKey,
			KeyDir:      ejsonKeyDir,
			SkipDecrypt: false,
		}

		// as we want to retrieve the kubeconfig, it makes no sense to
		// skip kubeconfig retrieval
		skipKubeconfig := false
		inv, err := inventory.NewInventory(inventoryPath, ejsonSettings, skipKubeconfig)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to compile inventory.")
			return
		}

		entry, ok := inv.Entries()[name]
		if !ok {
			log.WithFields(log.Fields{
				"entry": name,
			}).Fatal("Entry does not exist")
		}

		kubeconfig, err := entry.Kubeconfig().Yaml()
		if err != nil {
			log.WithFields(log.Fields{
				"entry": name,
				"error": err.Error(),
			}).Fatal("Failed to get kubeconfig")
		}
		fmt.Println(string(kubeconfig))
	},
}

var inventoryValuesCmd = &cobra.Command{
	Use:   "values [entry name]",
	Short: "Get the values for a given inventory entry",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		groupVarsDir := viper.GetString("group-vars-dir")
		inventoryPath := viper.GetString("inventory")
		skipDecrypt := viper.GetBool("skip-decrypt")
		skipEval := viper.GetBool("skip-eval")
		ejsonPrivKey := viper.GetString("ejson-privkey")
		ejsonKeyDir := viper.GetString("ejson-key-dir")

		ejsonSettings := ejson.Settings{
			PrivKey:     ejsonPrivKey,
			KeyDir:      ejsonKeyDir,
			SkipDecrypt: false,
		}

		// we just need the values for the given entry, skip the kubeconfig retrieval
		skipKubeconfig := true
		inv, err := inventory.NewInventory(inventoryPath, ejsonSettings, skipKubeconfig)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to compile inventory.")
			return
		}

		entry, ok := inv.Entries()[name]
		if !ok {
			log.WithFields(log.Fields{
				"entry": name,
			}).Fatal("Entry does not exist")
		}

		ejsonSettings.SkipDecrypt = skipDecrypt
		target, err := target.New(entry, groupVarsDir, skipEval, &ejsonSettings)
		if err != nil {
			log.WithFields(log.Fields{
				"entry": name,
				"error": err.Error(),
			}).Fatal("Failed to compile values for inventory entry")
		}
		values := target.Values()

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
		fmt.Printf("%s", string(result))
	},
}

var inventoryLoaderCmd = &cobra.Command{
	Use:   "loader [entry name]",
	Short: "Get kubeconfig loader information for the given entry",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		inventoryPath := viper.GetString("inventory")
		unsafe := viper.GetBool("unsafe")
		ejsonPrivKey := viper.GetString("ejson-privkey")
		ejsonKeyDir := viper.GetString("ejson-key-dir")

		ejsonSettings := ejson.Settings{
			PrivKey:     ejsonPrivKey,
			KeyDir:      ejsonKeyDir,
			SkipDecrypt: false,
		}

		// we just need the values for the given entry, skip the kubeconfig retrieval
		skipKubeconfig := true
		inv, err := inventory.NewInventory(inventoryPath, ejsonSettings, skipKubeconfig)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to compile inventory.")
		}

		entry, ok := inv.Entries()[name]
		if !ok {
			log.WithFields(log.Fields{
				"entry": name,
			}).Fatal("Entry does not exist")
		}
		loaderConfig, err := entry.Kubeconfig().Loader().Config().Yaml(unsafe)
		if err != nil {
			log.WithFields(log.Fields{
				"entry": name,
				"error": err.Error(),
			}).Fatal("Failed to get loader config")
		}
		fmt.Printf("Loader type: %s\n", entry.Kubeconfig().Loader().Type())
		fmt.Println("Loader config: ")
		fmt.Printf("%4s", string(loaderConfig))
	},
}

func init() {
	inventoryValuesCmd.Flags().BoolP("json", "j", false, "Output json instead of yaml")
	inventoryLoaderCmd.Flags().BoolP("unsafe", "", false, "Show confidential loader info")
	viper.BindPFlag("json", inventoryValuesCmd.Flags().Lookup("json"))
	viper.BindPFlag("unsafe", inventoryLoaderCmd.Flags().Lookup("unsafe"))

	inventoryCmd.AddCommand(inventoryListCmd)
	inventoryCmd.AddCommand(inventoryKubeconfigCmd)
	inventoryCmd.AddCommand(inventoryValuesCmd)
	inventoryCmd.AddCommand(inventoryLoaderCmd)

	rootCmd.AddCommand(inventoryCmd)
}
