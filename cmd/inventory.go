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
	"github.com/spf13/cobra"
)

func newInventoryCmd(c *Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "inventory",
		Short:                 "Get inventory information",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(
		newInventoryListCmd(c),
		newInventoryKubeconfigCmd(c),
		newInventoryValuesCmd(c),
		newInventoryLoaderCmd(c),
	)
	return cmd
}

func addInventoryFlags(cmd *cobra.Command) {
	addEjsonFlags(cmd)
	addEvalFlags(cmd)
	addLimitFlags(cmd)
	cmd.Flags().StringP("inventory", "i", "inventory.yml", "Path to the inventory")
}
