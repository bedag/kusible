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
	"github.com/bedag/kusible/pkg/target"
	"github.com/bedag/kusible/pkg/wrapper/ejson"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newInventoryValuesCmd(c *Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "values [entry name]",
		Short:                 "Get the values for a given inventory entry",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  c.wrap(runInventoryValues),
	}
	addInventoryFlags(cmd)
	addGroupsFlags(cmd)
	addOutputFormatFlags(cmd)

	return cmd
}

func runInventoryValues(c *Cli, cmd *cobra.Command, args []string) error {
	name := args[0]
	groupVarsDir := c.viper.GetString("group-vars-dir")
	inventoryPath := c.viper.GetString("inventory")
	skipDecrypt := c.viper.GetBool("skip-decrypt")
	skipEval := c.viper.GetBool("skip-eval")
	ejsonPrivKey := c.viper.GetString("ejson-privkey")
	ejsonKeyDir := c.viper.GetString("ejson-key-dir")

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

	ejsonSettings.SkipDecrypt = skipDecrypt
	target, err := target.New(entry, groupVarsDir, skipEval, &ejsonSettings)
	if err != nil {
		log.WithFields(log.Fields{
			"entry": name,
			"error": err.Error(),
		}).Error("Failed to compile values for inventory entry")
		return err
	}
	values := target.Values()

	render := values.YAML
	if c.viper.GetBool("json") {
		render = values.JSON
	}

	result, err := render()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to render compiled group vars.")
		return err
	}

	if !c.viper.GetBool("quiet") {
		fmt.Printf("%s", string(result))
	}
	return nil
}
