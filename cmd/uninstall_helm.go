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
	"fmt"

	"github.com/bedag/kusible/pkg/printer"
	helmutil "github.com/bedag/kusible/pkg/wrapper/helm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newUninstallHelmCmd(c *Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "helm [playbook]",
		Short:                 "Uninstall an application deployed with helm",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  c.wrap(runUninstallHelm),
	}

	addUninstallFlags(cmd)
	helmutil.AddHelmUninstallFlags(cmd)

	return cmd
}

func runUninstallHelm(c *Cli, cmd *cobra.Command, args []string) error {
	playbookFile := args[0]

	inv, err := getInventoryWithKubeconfig(c)
	if err != nil {
		return err
	}

	playbookSet, err := loadPlaybooks(c, playbookFile)
	if err != nil {
		return err
	}

	helmOptions := helmutil.NewOptions(c.viper)

	status := map[string][]string{}
	for name, playbook := range playbookSet {
		entry := inv.Entries()[name]
		entryStatus := []string{}
		for _, play := range playbook.Config.Plays {
			helm, err := helmutil.NewWithGetter(helmOptions, c.HelmEnv, entry.Kubeconfig())
			if err != nil {
				return fmt.Errorf("failed to create helm client instance: %s", err)
			}
			playStatus, err := helm.UninstallPlay(play)
			entryStatus = append(entryStatus, playStatus...)
			if err != nil {
				log.WithFields(log.Fields{
					"play":  play.Name,
					"entry": name,
					"error": err.Error(),
				}).Error("Failed to uninstall application with helm.")
				status[name] = entryStatus
				outErr := c.output(uninstallHelmStatusQueue(status))
				if outErr != nil {
					return fmt.Errorf("%s + %s", err, outErr)
				}
				return err
			}
		}
		status[name] = entryStatus
	}

	return c.output(uninstallHelmStatusQueue(status))
}

func uninstallHelmStatusQueue(releases map[string][]string) printer.Queue {
	printerQueue := printer.Queue{}

	for name, entryStatus := range releases {
		name := name
		entryStatus := entryStatus
		job := printer.NewJob(func(fields []string) map[string]interface{} {
			defaultResult := map[string]interface{}{
				"entry":  name,
				"status": entryStatus,
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
	return printerQueue
}
