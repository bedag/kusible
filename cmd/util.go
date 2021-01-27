/*
Copyright © 2021 Michael Gruener

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
	"github.com/bedag/kusible/pkg/inventory"
	invconfig "github.com/bedag/kusible/pkg/inventory/config"
	"github.com/bedag/kusible/pkg/playbook"
	"github.com/bedag/kusible/pkg/target"
	"github.com/bedag/kusible/pkg/wrapper/ejson"
	log "github.com/sirupsen/logrus"
)

func getEjsonSettings(c *Cli) ejson.Settings {
	return ejson.Settings{
		PrivKey:     c.viper.GetString("ejson-privkey"),
		KeyDir:      c.viper.GetString("ejson-key-dir"),
		SkipDecrypt: c.viper.GetBool("skip-decrypt"),
	}
}

func loadInventory(c *Cli, skipKubeconfig bool) (*inventory.Inventory, error) {
	ejsonSettings := getEjsonSettings(c)
	inventoryPath := c.viper.GetString("inventory")

	clusterInventoryDefaults := invconfig.ClusterInventory{
		Namespace: c.viper.GetString("cluster-inventory-namespace"),
		ConfigMap: c.viper.GetString("cluster-inventory-configmap"),
	}

	inventory, err := inventory.NewInventory(inventoryPath, ejsonSettings, skipKubeconfig, clusterInventoryDefaults)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to compile inventory.")
		return nil, err
	}

	return inventory, nil
}

// wrapper around loadInventory() to make it more obvious what is intended
func getInventoryWithKubeconfig(c *Cli) (*inventory.Inventory, error) {
	return loadInventory(c, false)
}

// wrapper around loadInventory() to make it more obvious what is intended
func getInventoryWithoutKubeconfig(c *Cli) (*inventory.Inventory, error) {
	return loadInventory(c, true)
}

func loadTargets(c *Cli, filter string) (*target.Targets, error) {
	getInventory := getInventoryWithKubeconfig
	if c.viper.GetBool("skip-cluster-inventory") {
		getInventory = getInventoryWithoutKubeconfig
	}

	inv, err := getInventory(c)
	if err != nil {
		return nil, err
	}
	return loadTargetsWithInventory(c, filter, inv)
}

func loadTargetsWithInventory(c *Cli, filter string, inv *inventory.Inventory) (*target.Targets, error) {
	limits := c.viper.GetStringSlice("limit")
	groupVarsDir := c.viper.GetString("group-vars-dir")

	ejsonSettings := getEjsonSettings(c)

	targets, err := target.NewTargets(filter, limits, groupVarsDir, inv, true, &ejsonSettings)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to compile values for inventory entries.")
		return nil, err
	}

	return targets, nil
}

func loadPlaybooks(c *Cli, playbookFile string) (playbook.Set, error) {
	targets, err := loadTargets(c, ".*")
	if err != nil {
		return nil, err
	}
	return loadPlaybooksWithTargets(c, playbookFile, targets)
}

func loadPlaybooksWithTargets(c *Cli, playbookFile string, targets *target.Targets) (playbook.Set, error) {
	skipEval := c.viper.GetBool("skip-eval")
	skipClusterInv := c.viper.GetBool("skip-cluster-inventory")

	playbooks, err := playbook.NewSet(playbookFile, targets, skipEval, skipClusterInv)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to compile playbooks.")
		return nil, err
	}

	return playbooks, nil
}
