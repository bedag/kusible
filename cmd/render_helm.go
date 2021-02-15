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
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

func newRenderHelmCmd(c *Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "helm [playbook]",
		Short:                 "Use helm to render manifests for an inventory entry",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  c.wrap(runRenderHelm),
	}
	addRenderFlags(cmd)
	helmutil.AddHelmTemplateFlags(cmd)

	return cmd
}

func runRenderHelm(c *Cli, cmd *cobra.Command, args []string) error {
	playbookFile := args[0]

	playbookSet, err := loadPlaybooks(c, playbookFile)
	if err != nil {
		return err
	}

	helmOptions := helmutil.NewOptions(c.viper)
	helm, err := helmutil.New(helmOptions, c.HelmEnv, c.Log)
	if err != nil {
		return fmt.Errorf("failed to create helm client instance: %s", err)
	}

	bigManifest := ""
	for name, playbook := range playbookSet {
		for _, play := range playbook.Config.Plays {
			for _, repo := range play.Repos {
				c.Log.WithFields(logrus.Fields{
					"play":  play.Name,
					"repo":  repo.Name,
					"entry": name,
				}).Debug("Adding helm repository.")
				if err := helm.RepoAdd(repo.Name, repo.URL); err != nil {
					c.Log.WithFields(logrus.Fields{
						"play":  play.Name,
						"repo":  repo.Name,
						"entry": name,
						"error": err.Error(),
					}).Error("Failed to add helm repo for play.")
					return err
				}
			}
			c.Log.WithFields(logrus.Fields{
				"play":  play.Name,
				"entry": name,
			}).Debug("Rendering play charts.")
			manifest, err := helm.TemplatePlay(play)
			if err != nil {
				c.Log.WithFields(logrus.Fields{
					"play":  play.Name,
					"entry": name,
					"error": err.Error(),
				}).Error("Failed to render play manifests with helm.")
				return err
			}
			bigManifest = fmt.Sprintf("%s%s\n---", bigManifest, manifest)
		}
	}

	manifests, err := helmutil.SplitSortManifest(bigManifest)
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to split combinded playbook manifest into separate yaml resources.")
		return err
	}

	printerQueue := printer.Queue{}
	for _, manifest := range manifests {
		// see https://golang.org/doc/faq#closures_and_goroutines
		manifest := manifest
		job := printer.NewJob(func(fields []string) map[string]interface{} {
			var defaultResult map[string]interface{}
			yaml.Unmarshal([]byte(manifest), &defaultResult)

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
