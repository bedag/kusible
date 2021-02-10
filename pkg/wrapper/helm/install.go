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
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
)

func (h *Helm) getInstallOptions(client *action.Install) {
	client.CreateNamespace = h.options.CreateNamespace
	client.DryRun = h.options.DryRun
	client.DisableHooks = h.options.NoHooks
	client.Replace = h.options.Replace
	client.Timeout = h.options.Timeout
	client.Wait = h.options.Wait
	client.WaitForJobs = h.options.WaitForJobs
	client.DependencyUpdate = h.options.DepdencyUpdate
	client.DisableOpenAPIValidation = h.options.DisableOpenAPIValidation
	client.Atomic = h.options.Atomic
	client.SkipCRDs = h.options.SkipCRDs
	client.SubNotes = h.options.RenderSubChartNotes
	h.getChartPathOptions(&client.ChartPathOptions)
}

func (h *Helm) runInstall(args []string, vals map[string]interface{}, client *action.Install) (*release.Release, error) {
	out := h.out
	settings := h.settings
	//debug("Original chart version: %q", client.Version)
	if client.Version == "" && client.Devel {
		//debug("setting version to >0.0.0-0")
		client.Version = ">0.0.0-0"
	}

	name, chart, err := client.NameAndChart(args)
	if err != nil {
		return nil, err
	}
	client.ReleaseName = name

	cp, err := client.ChartPathOptions.LocateChart(chart, settings)
	if err != nil {
		return nil, err
	}

	//debug("CHART PATH: %s\n", cp)

	p := getter.All(settings)
	//vals, err := valueOpts.MergeValues(p)
	//if err != nil {
	//	return nil, err
	//}

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	if err := checkIfInstallable(chartRequested); err != nil {
		return nil, err
	}

	//if chartRequested.Metadata.Deprecated {
	//	warning("This chart is deprecated")
	//}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					Out:              out,
					ChartPath:        cp,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: settings.RepositoryConfig,
					RepositoryCache:  settings.RepositoryCache,
					Debug:            settings.Debug,
				}
				if err := man.Update(); err != nil {
					return nil, err
				}
				// Reload the chart with the updated Chart.lock file.
				if chartRequested, err = loader.Load(cp); err != nil {
					return nil, errors.Wrap(err, "failed reloading chart after repo update")
				}
			} else {
				return nil, err
			}
		}
	}

	if client.Namespace == "" {
		client.Namespace = settings.Namespace()
	}
	return client.Run(chartRequested, vals)
}

// Application chart type is only installable
func checkIfInstallable(ch *chart.Chart) error {
	switch ch.Metadata.Type {
	case "", "application":
		return nil
	}
	return errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}
