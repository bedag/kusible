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
	argocdutil "github.com/bedag/kusible/pkg/wrapper/argocd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

func newRenderArgoCDCmd(c *Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "argocd [playbook]",
		Short:                 "Use render the given playbook into a set of ArgoCD Application resources",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  c.wrap(runRenderArgoCD),
	}
	addRenderFlags(cmd)
	cmd.Flags().String("argocd-namespace", "argocd", "Namespace where ArgoCD is looking for ArgoCD applications")
	cmd.Flags().String("argocd-project", "default", "The ArgoCD project to which the applications should be assigned")

	return cmd
}

func runRenderArgoCD(c *Cli, cmd *cobra.Command, args []string) error {
	playbookFile := args[0]
	namespace := c.viper.GetString("argocd-namespace")
	project := c.viper.GetString("argocd-project")

	playbookSet, err := loadPlaybooks(c, playbookFile)
	if err != nil {
		return err
	}

	allApps := []argocdutil.Application{}
	for name, playbook := range playbookSet {
		for _, play := range playbook.Config.Plays {
			c.Log.WithFields(logrus.Fields{
				"play":  play.Name,
				"entry": name,
			}).Debug("Rendering play.")

			apps, err := argocdutil.ApplicationsFromPlay(play, project, namespace, name)
			if err != nil {
				c.Log.WithFields(logrus.Fields{
					"play":  play.Name,
					"entry": name,
					"error": err.Error(),
				}).Error("Failed to render ArgoCD application manifests.")

				return err
			}
			allApps = append(allApps, apps...)
		}
	}

	printerQueue := printer.Queue{}
	for _, app := range allApps {
		// see https://golang.org/doc/faq#closures_and_goroutines
		app := app
		job := printer.NewJob(func(fields []string) map[string]interface{} {
			manifest, _ := yaml.Marshal(app)
			var defaultResult map[string]interface{}
			yaml.Unmarshal(manifest, &defaultResult)

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
