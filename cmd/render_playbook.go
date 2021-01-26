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
	"github.com/bedag/kusible/pkg/playbook"
	"github.com/bedag/kusible/pkg/target"
	"github.com/bedag/kusible/pkg/wrapper/ejson"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newRenderPlaybookCmd(c *Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "playbook [playbook]",
		Short:                 "Render the given playbook",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  c.wrap(runRenderPlaybook),
	}
	addRenderFlags(cmd)

	return cmd
}

func runRenderPlaybook(c *Cli, cmd *cobra.Command, args []string) error {
	playbookFile := args[0]
	limits := c.viper.GetStringSlice("limit")
	groupVarsDir := c.viper.GetString("group-vars-dir")
	inventoryPath := c.viper.GetString("inventory")
	skipEval := c.viper.GetBool("skip-eval")
	skipDecrypt := c.viper.GetBool("skip-decrypt")
	ejsonPrivKey := c.viper.GetString("ejson-privkey")
	ejsonKeyDir := c.viper.GetString("ejson-key-dir")
	skipClusterInv := c.viper.GetBool("skip-cluster-inventory")
	clusterInvNamespace := c.viper.GetString("cluster-inventory-namespace")
	clusterInvConfigMap := c.viper.GetString("cluster-inventory-configmap")

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

	clusterInventoryDefaults := invconfig.ClusterInventory{
		Namespace: clusterInvNamespace,
		ConfigMap: clusterInvConfigMap,
	}

	// if we do not retrieve the cluster inventory ConfigMap, we do not need to retrieve
	// the kubeconfig
	inventory, err := inventory.NewInventory(inventoryPath, invEjsonSettings, false, clusterInventoryDefaults)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to compile inventory.")
		return err
	}

	targets, err := target.NewTargets(".*", limits, groupVarsDir, inventory, true, &tgtEjsonSettings)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to compile values for inventory entries.")
		return err
	}

	playbookSet, err := playbook.NewSet(playbookFile, targets, skipEval, skipClusterInv)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to compile playbooks.")
		return err
	}

	for name, playbook := range playbookSet {
		result, err := playbook.YAML(skipEval)
		if err != nil {
			log.WithFields(log.Fields{
				"entry": name,
				"error": err.Error(),
			}).Error("Failed to convert playbook entry to yaml.")
			return err
		}
		if len(result) > 0 {
			fmt.Printf("======= Plays for %s =======\n", name)
			fmt.Printf("%s", string(result))
		}
	}

	return nil
}
