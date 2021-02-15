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

	"github.com/bedag/kusible/pkg/playbook/config"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

// DeployPlay renders all charts contained in a given play to a string containing
// kubernetes manifests
func (h *Helm) DeployPlay(play *config.Play) ([]*release.Release, error) {
	releases := make([]*release.Release, len(play.Charts))
	for _, chart := range play.Charts {
		actionConfig, err := h.ActionConfig(chart.Namespace)
		if err != nil {
			return releases, fmt.Errorf("failed initialize helm client: %s", err)
		}
		client := action.NewUpgrade(actionConfig)
		h.getUpgradeOptions(client)

		for _, pr := range play.Repos {
			if pr.Name == chart.Repo {
				client.ChartPathOptions.RepoURL = pr.URL
			}
		}

		if client.ChartPathOptions.RepoURL == "" {
			return releases, fmt.Errorf("no repo '%s' for chart '%s' configured in play", chart.Repo, chart.Name)
		}

		client.Install = true
		client.Version = chart.Version
		client.Namespace = chart.Namespace

		chartName := chart.Chart
		releaseName := chart.Name
		values := chart.Values

		rel, err := h.runUpgrade([]string{releaseName, chartName}, values, client, actionConfig)
		if err != nil {
			return releases, fmt.Errorf("failed to deploy chart '%s' as release '%s': %s", chartName, releaseName, err)
		}
		releases = append(releases, rel)

	}
	return releases, nil
}
