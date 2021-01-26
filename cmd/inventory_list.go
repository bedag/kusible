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

func newInventoryListCmd(c *Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "list [regex]",
		Short:                 "List inventory entries",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  c.wrap(runInventoryList),
	}
	addInventoryFlags(cmd)

	return cmd
}

func runInventoryList(c *Cli, cmd *cobra.Command, args []string) error {
	filter := args[0]
	limits := c.viper.GetStringSlice("limit")
	inventoryPath := c.viper.GetString("inventory")
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

	// as we just want to list the available inventory entries we can (and should)
	// skip kubeconfig retrieval
	skipKubeconfig := true
	inventory, err := inventory.NewInventory(inventoryPath, ejsonSettings, skipKubeconfig, clusterInventoryDefaults)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to compile inventory.")
		return err
	}

	names, err := inventory.EntryNames(filter, limits)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to get list of entries")
		return err
	}
	for _, name := range names {
		fmt.Printf("%s\n", name)
	}
	return nil
}
