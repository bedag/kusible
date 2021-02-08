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

	"github.com/bedag/kusible/pkg/printer"
	helmutil "github.com/bedag/kusible/pkg/wrapper/helm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/release"
	"sigs.k8s.io/yaml"
)

func newDeployHelmCmd(c *Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "helm [playbook]",
		Short:                 "Use helm to deploy an application",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  c.wrap(runDeployHelm),
	}
	addDeployFlags(cmd)
	helmutil.AddHelmInstallFlags(cmd)
	helmutil.AddHelmChartPathOptionsFlags(cmd)
	helmutil.AddHelmTemplateFlags(cmd)

	return cmd
}

func runDeployHelm(c *Cli, cmd *cobra.Command, args []string) error {
	playbookFile := args[0]

	inv, err := getInventoryWithKubeconfig(c)
	if err != nil {
		return err
	}

	playbookSet, err := loadPlaybooks(c, playbookFile)
	if err != nil {
		return err
	}

	helmGlobals := helmutil.GlobalsFromViper(c.viper)

	releases := []*release.Release{}
	for name, playbook := range playbookSet {
		entry := inv.Entries()[name]
		for _, play := range playbook.Config.Plays {
			helm, err := helmutil.NewWithGetter(helmGlobals, entry.Kubeconfig())
			if err != nil {
				return fmt.Errorf("failed to create helm client instance: %s", err)
			}
			for _, repo := range play.Repos {
				if err := helm.RepoAdd(repo.Name, repo.URL); err != nil {
					log.WithFields(log.Fields{
						"play":  play.Name,
						"repo":  repo.Name,
						"entry": name,
						"error": err.Error(),
					}).Error("Failed to add helm repo for play.")
					outErr := c.output(deployHelmStatusQueue(releases))
					if outErr != nil {
						return fmt.Errorf("%s + %s", err, outErr)
					}
					return err
				}
			}
			rels, err := helm.InstallPlay(play)
			if err != nil {
				log.WithFields(log.Fields{
					"play":  play.Name,
					"entry": name,
					"error": err.Error(),
				}).Error("Failed to render play manifests with helm.")
				outErr := c.output(deployHelmStatusQueue(releases))
				if outErr != nil {
					return fmt.Errorf("%s + %s", err, outErr)
				}
				return err
			}
			releases = append(releases, rels...)
		}
	}

	return c.output(deployHelmStatusQueue(releases))
}

func deployHelmStatusQueue(releases []*release.Release) printer.Queue {
	printerQueue := printer.Queue{}
	for _, rel := range releases {
		// see https://golang.org/doc/faq#closures_and_goroutines
		rel := rel
		job := printer.NewJob(func(fields []string) map[string]interface{} {
			var defaultResult map[string]interface{}
			status, _ := yaml.Marshal(rel)
			yaml.Unmarshal([]byte(status), &defaultResult)

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
