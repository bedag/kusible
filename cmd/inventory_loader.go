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

	"github.com/bedag/kusible/pkg/inventory"
	invconfig "github.com/bedag/kusible/pkg/inventory/config"
	"github.com/bedag/kusible/pkg/wrapper/ejson"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newInventoryLoaderCmd(c *Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "loader [entry name]",
		Short:                 "Get kubeconfig loader information for the given entry",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  c.wrap(runInventoryLoader),
	}
	addInventoryFlags(cmd)
	cmd.Flags().BoolP("unsafe", "", false, "Show confidential loader info")

	return cmd
}

func runInventoryLoader(c *Cli, cmd *cobra.Command, args []string) error {
	name := args[0]
	inventoryPath := c.viper.GetString("inventory")
	unsafe := c.viper.GetBool("unsafe")
	ejsonPrivKey := c.viper.GetString("ejson-privkey")
	ejsonKeyDir := c.viper.GetString("ejson-key-dir")
	skipClusterInv := c.viper.GetBool("skip-cluster-inventory")
	clusterInvNamespace := c.viper.GetString("cluster-inventory-namespace")
	clusterInvConfigMap := c.viper.GetString("cluster-inventory-configmap")

	ejsonSettings := ejson.Settings{
		PrivKey:     ejsonPrivKey,
		KeyDir:      ejsonKeyDir,
		SkipDecrypt: false,
	}

	if skipClusterInv {
		log.Info("--skip-cluster-inventory has no effect when retrieving kubeconfig files")
	}

	clusterInventoryDefaults := invconfig.ClusterInventory{
		Namespace: clusterInvNamespace,
		ConfigMap: clusterInvConfigMap,
	}

	// we just need the values for the given entry, skip the kubeconfig retrieval
	skipKubeconfig := true
	inv, err := inventory.NewInventory(inventoryPath, ejsonSettings, skipKubeconfig, clusterInventoryDefaults)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to compile inventory.")
		return err
	}

	entry, ok := inv.Entries()[name]
	if !ok {
		log.WithFields(log.Fields{
			"entry": name,
		}).Error("Entry does not exist")
		return err
	}
	loaderConfig, err := entry.Kubeconfig().Loader().Config().Yaml(unsafe)
	if err != nil {
		log.WithFields(log.Fields{
			"entry": name,
			"error": err.Error(),
		}).Error("Failed to get loader config")
		return err
	}
	fmt.Printf("Loader type: %s\n", entry.Kubeconfig().Loader().Type())
	fmt.Println("Loader config: ")
	fmt.Printf("%4s", string(loaderConfig))
	return nil
}
