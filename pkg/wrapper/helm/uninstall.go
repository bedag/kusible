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

	"github.com/bedag/kusible/pkg/playbook/config"
	"helm.sh/helm/v3/pkg/action"
)

func (h *Helm) UninstallPlay(play *config.Play) ([]string, error) {
	result := []string{}
	for _, chart := range play.Charts {
		actionConfig, err := h.ActionConfig(chart.Namespace)
		if err != nil {
			return result, fmt.Errorf("failed initialize helm client: %s", err)
		}
		client := action.NewUninstall(actionConfig)
		h.getUninstallOptions(client)

		releaseName := chart.Name
		status, err := h.runUninstall(releaseName, client)
		if err != nil {
			return result, fmt.Errorf("failed to uninstall release '%s': %s", releaseName, err)
		}
		result = append(result, status)
	}
	return result, nil
}

func (h *Helm) runUninstall(name string, client *action.Uninstall) (string, error) {
	res, err := client.Run(name)
	if err != nil {
		return "", err
	}
	if res != nil && res.Info != "" {
		return res.Info, nil
	}
	return fmt.Sprintf("release '%s' uninstalled", name), nil
}

func (h *Helm) getUninstallOptions(client *action.Uninstall) {
	client.DryRun = h.options.DryRun
	client.DisableHooks = h.options.NoHooks
	client.Timeout = h.options.Timeout
	client.KeepHistory = h.options.KeepHistory
}
