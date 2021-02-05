/*
Copyright © 2021 Michael Gruener & The Helm Authors

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

/* Lots of code straight from github.com/helm/helm and adapted to be used here */

package helm

import (
	"fmt"
	"strings"

	"github.com/bedag/kusible/pkg/playbook/config"
	"helm.sh/helm/v3/pkg/action"
	helmcli "helm.sh/helm/v3/pkg/cli"
)

// TemplatePlay renders all charts contained in a given play to a string containing
// kubernetes manifests
func TemplatePlay(play *config.Play, settings *helmcli.EnvSettings) (string, error) {
	result := ""
	for _, chart := range play.Charts {
		actionConfig := new(action.Configuration)
		client := action.NewInstall(actionConfig)
		client.DryRun = true
		client.Replace = true
		client.ClientOnly = true
		client.IncludeCRDs = true

		for _, pr := range play.Repos {
			if pr.Name == chart.Repo {
				client.ChartPathOptions.RepoURL = pr.URL
			}
		}

		if client.ChartPathOptions.RepoURL == "" {
			return result, fmt.Errorf("no repo '%s' for chart '%s' configured in play", chart.Repo, chart.Name)
		}

		client.ReleaseName = chart.Name
		client.Version = chart.Version
		client.Namespace = chart.Namespace

		name := chart.Chart
		values := chart.Values

		manifests, err := template(chart.Name, name, values, client, settings)
		if err != nil {
			return result, err
		}
		result = fmt.Sprintf("%s%s\n", result, strings.TrimSpace(manifests))
	}
	return result, nil
}

// Template renders a given chart + values to a string containing
// kubernetes manifests
func template(release string, chart string, values map[string]interface{}, client *action.Install, settings *helmcli.EnvSettings) (string, error) {
	args := []string{release, chart}
	rel, err := runInstall(args, values, client, settings)

	if err != nil && !settings.Debug {
		if rel != nil {
			return "", fmt.Errorf("%w\n\nUse the HELM_DEBUG env var to render out invalid YAML", err)
		}
		return "", err
	}

	return rel.Manifest, nil
}
