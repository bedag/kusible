/*
Copyright © 2019 Copyright © 2021 Bedag Informatik AG

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

func newInventoryLoaderCmd(c *Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "loader [regex]",
		Short:                 "Get kubeconfig loader information for entries matched by the regex",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  c.wrap(runInventoryLoader),
	}
	addInventoryFlags(cmd)
	cmd.Flags().Bool("unsafe", false, "Show confidential loader info")

	return cmd
}

func runInventoryLoader(c *Cli, cmd *cobra.Command, args []string) error {
	filter := args[0]
	limits := c.viper.GetStringSlice("limit")
	unsafe := c.viper.GetBool("unsafe")

	// we just need the values for the given entry, skip the kubeconfig retrieval
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

	printerQueue := printer.Queue{}
	for _, name := range names {
		entry := inv.Entries()[name]
		// see https://golang.org/doc/faq#closures_and_goroutines
		name := name

		loaderType := entry.Kubeconfig().Loader().Type()
		loaderConfig := entry.Kubeconfig().Loader().Config().Sanitize()
		if unsafe {
			loaderConfig = entry.Kubeconfig().Loader().Config()
		}

		job := printer.NewJob(func(fields []string) map[string]interface{} {
			defaultResult := map[string]interface{}{
				"entry":  name,
				"type":   loaderType,
				"config": loaderConfig,
			}

			if len(fields) < 1 {
				return defaultResult
			}

			result := map[string]interface{}{}
			for _, field := range fields {
				if val, ok := defaultResult[field]; ok {
					result[field] = val
				}
			}
			return result
		})
		printerQueue = append(printerQueue, job)
	}

	return c.output(printerQueue)
}
