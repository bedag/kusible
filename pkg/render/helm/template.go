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

package helm

import (
	"fmt"
	"os"
	"strings"

	"github.com/bedag/kusible/pkg/playbook/config"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	helmcli "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
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

		client.ReleaseName = chart.Name
		client.Version = chart.Version
		client.Namespace = chart.Namespace

		name := chart.Chart
		values := chart.Values

		manifests, err := Template(name, values, client, settings)
		if err != nil {
			return result, err
		}
		result = fmt.Sprintf("%s%s\n", result, strings.TrimSpace(manifests))
	}
	return result, nil
}

// Template renders a given chart + values to a string containing
// kubernetes manifests
func Template(chart string, values map[string]interface{}, client *action.Install, settings *helmcli.EnvSettings) (string, error) {

	cp, err := client.ChartPathOptions.LocateChart(chart, settings)
	if err != nil {
		return "", err
	}

	p := getter.All(settings)

	chartRequested, err := loader.Load(cp)
	if err != nil {
		return "", err
	}

	if chartRequested.Metadata.Type != "" && chartRequested.Metadata.Type != "application" {
		return "", fmt.Errorf("%s charts are not installable", chartRequested.Metadata.Type)
	}

	// TODO: warning?
	//if chartRequested.Metadata.Deprecated {
	//	return "", fmt.Errorf("this chart is deprecated")
	//}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					Out:              os.Stdout,
					ChartPath:        cp,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: settings.RepositoryConfig,
					RepositoryCache:  settings.RepositoryCache,
					Debug:            settings.Debug,
				}
				if err := man.Update(); err != nil {
					return "", err
				}
				// Reload the chart with the updated Chart.lock file.
				if chartRequested, err = loader.Load(cp); err != nil {
					return "", errors.Wrap(err, "failed reloading chart after repo update")
				}
			} else {
				return "", err
			}
		}
	}

	rel, err := client.Run(chartRequested, values)

	if err != nil && !settings.Debug {
		if rel != nil {
			return "", fmt.Errorf("%w\n\nUse the HELM_DEBUG env var to render out invalid YAML", err)
		}
		return "", err
	}

	return rel.Manifest, nil
}
