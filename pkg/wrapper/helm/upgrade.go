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
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
)

func (h *Helm) getUpgradeOptions(client *action.Upgrade) {
	client.DryRun = h.options.DryRun
	client.DisableHooks = h.options.NoHooks
	client.Timeout = h.options.Timeout
	client.Wait = h.options.Wait
	client.WaitForJobs = h.options.WaitForJobs
	client.DisableOpenAPIValidation = h.options.DisableOpenAPIValidation
	client.Atomic = h.options.Atomic
	client.SkipCRDs = h.options.SkipCRDs
	client.SubNotes = h.options.RenderSubChartNotes
	client.Force = h.options.Force
	client.ResetValues = h.options.ResetValues
	client.ReuseValues = h.options.ReuseValues
	client.MaxHistory = h.options.HistoryMax
	client.CleanupOnFail = h.options.CleanupOnFail
	h.getChartPathOptions(&client.ChartPathOptions)
}

func (h *Helm) runUpgrade(args []string, vals map[string]interface{}, client *action.Upgrade, cfg *action.Configuration) (*release.Release, error) {
	settings := h.settings
	if client.Namespace == "" {
		client.Namespace = settings.Namespace()
	}

	if client.Install {
		// If a release does not exist, install it.
		histClient := action.NewHistory(cfg)
		histClient.Max = 1
		if _, err := histClient.Run(args[0]); err == driver.ErrReleaseNotFound {
			// Only print this to stdout for table output
			//if outfmt == output.Table {
			//  fmt.Fprintf(h.out, "Release %q does not exist. Installing it now.\n", args[0])
			//}
			instClient := action.NewInstall(cfg)
			h.getInstallOptions(instClient)
			instClient.Version = client.Version
			instClient.Namespace = client.Namespace
			instClient.ChartPathOptions.RepoURL = client.ChartPathOptions.RepoURL

			return h.runInstall(args, vals, instClient)
		} else if err != nil {
			return nil, err
		}
	}

	if client.Version == "" && client.Devel {
		//debug("setting version to >0.0.0-0")
		client.Version = ">0.0.0-0"
	}

	chartPath, err := client.ChartPathOptions.LocateChart(args[1], settings)
	if err != nil {
		return nil, err
	}

	//vals, err := valueOpts.MergeValues(getter.All(settings))
	//if err != nil {
	//	return err
	//}

	// Check chart dependencies to make sure all are present in /charts
	ch, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}
	if req := ch.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(ch, req); err != nil {
			return nil, err
		}
	}

	//if ch.Metadata.Deprecated {
	//	warning("This chart is deprecated")
	//}

	return client.Run(args[0], ch, vals)
}
