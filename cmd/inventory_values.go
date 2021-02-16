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
	"github.com/bedag/kusible/internal/third_party/deepcopy"
	"github.com/bedag/kusible/internal/wrapper/spruce"
	"github.com/bedag/kusible/pkg/printer"
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
)

func newInventoryValuesCmd(c *Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "values [filter]",
		Short:                 "Get the values all inventory entries matched by filter",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  c.wrap(runInventoryValues),
	}
	addInventoryFlags(cmd)
	addGroupsFlags(cmd)
	addSkipClusterInventoryFlags(cmd)

	return cmd
}

func runInventoryValues(c *Cli, cmd *cobra.Command, args []string) error {
	filter := args[0]
	skipClusterInv := c.viper.GetBool("skip-cluster-inventory")
	skipEval := c.viper.GetBool("skip-eval")

	targets, err := loadTargets(c, filter)
	if err != nil {
		return err
	}

	printerQueue := printer.Queue{}
	for name, target := range targets.Targets() {
		values := target.Values().Map()
		clusterInventory := map[string]interface{}{}

		if !skipClusterInv {
			ci, err := target.Entry().ClusterInventory()
			if err != nil {
				return err
			}
			clusterInventory = *ci
		}

		// see https://golang.org/doc/faq#closures_and_goroutines
		name := name

		job := printer.NewJob(func(fields []string) map[string]interface{} {
			// TODO error handling
			mergeResult, _ := deepcopy.Map(clusterInventory)
			mergo.Merge(&mergeResult, values, mergo.WithOverride)
			spruce.Eval(&mergeResult, skipEval, []string{})

			defaultResult := map[string]interface{}{
				"entry":  name,
				"values": mergeResult,
			}

			if len(fields) < 1 {
				return defaultResult
			}

			resultValues := map[string]interface{}{}
			for _, field := range fields {
				if val, ok := values[field]; ok {
					resultValues[field] = val
				}
			}

			return map[string]interface{}{
				"entry":  name,
				"values": resultValues,
			}
		})
		printerQueue = append(printerQueue, job)
	}
	return c.output(printerQueue)
}
