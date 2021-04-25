/*
Copyright © 2021 Bedag Informatik AG

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
	"strings"

	"github.com/bedag/kusible/pkg/inventory"
	invconfig "github.com/bedag/kusible/pkg/inventory/config"
	"github.com/bedag/kusible/pkg/playbook"
	"github.com/bedag/kusible/pkg/target"
	"github.com/bedag/kusible/pkg/wrapper/ejson"
	"github.com/sirupsen/logrus"
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

	c.Log.WithFields(logrus.Fields{
		"path":              inventoryPath,
		"load-kubeconfig":   !skipKubeconfig,
		"cluster-inventory": fmt.Sprintf("%s/%s", clusterInventoryDefaults.Namespace, clusterInventoryDefaults.ConfigMap),
	}).Trace("Loading inventory.")

	inventory, err := inventory.NewInventory(inventoryPath, ejsonSettings, skipKubeconfig, clusterInventoryDefaults)
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to compile inventory.")
		return nil, err
	}

	c.Log.WithFields(logrus.Fields{
		"entries": len(inventory.Entries()),
	}).Trace("Successfully loaded inventory.")

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

	c.Log.WithFields(logrus.Fields{
		"limits":         strings.Join(limits, ","),
		"filter":         filter,
		"group-vars-dir": groupVarsDir,
	}).Trace("Loading targets from inventory.")

	targets, err := target.NewTargets(filter, limits, groupVarsDir, inv, true, &ejsonSettings)
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to compile values for inventory entries.")
		return nil, err
	}

	c.Log.WithFields(logrus.Fields{
		"targets": len(targets.Targets()),
	}).Trace("Successfully loaded targets from inventory.")

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

	c.Log.WithFields(logrus.Fields{
		"playbook-file":          playbookFile,
		"spruce-eval":            !skipEval,
		"load-cluster-inventory": !skipClusterInv,
	}).Trace("Loading playbooks for targets.")

	playbooks, err := playbook.NewSet(playbookFile, targets, skipEval, skipClusterInv)
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to compile playbooks.")
		return nil, err
	}

	if c.Log.IsLevelEnabled(logrus.TraceLevel) {
		plays := 0
		charts := 0
		for _, playbook := range playbooks {
			plays = plays + len(playbook.Config.Plays)
			for _, play := range playbook.Config.Plays {
				charts = charts + len(play.Charts)
			}
		}
		c.Log.WithFields(logrus.Fields{
			"target-playbooks": len(playbooks),
			"plays":            plays,
			"charts":           charts,
		}).Trace("Successfully loaded playbooks for targets.")
	}

	return playbooks, nil
}
