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
	"github.com/bedag/kusible/pkg/printer"
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

	// as we just want to list the available inventory entries we can (and should)
	// skip kubeconfig retrieval
	inv, err := getInventoryWithoutKubeconfig(c)
	if err != nil {
		return err
	}

	names, err := inv.EntryNames(filter, limits)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to get list of entries")
		return err
	}

	printFn := func(fields []string) map[string]interface{} {
		return map[string]interface{}{
			"entries": names,
		}
	}

	printerQueue := printer.Queue{printer.NewJob(printFn)}

	return c.output(printerQueue)
}
