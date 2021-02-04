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

package argocd

import (
	// "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"fmt"

	"github.com/bedag/kusible/pkg/playbook/config"
	"sigs.k8s.io/yaml"
)

// ApplicationFromPlay renders a a set of ArgoCD Application resources (see https://argoproj.github.io/argo-cd/operator-manual/declarative-setup/)
// for a given play. Each chart of the play results in a separate Application resource containing
// the details of the helm release (release name, chart name, repo of the chart, chart version, values).
//
// The project parameter is the argocd project the application should belong to
// The namespace parameter is the namespace where ArgoCD is expection Application resources
// The server parameter is the server name(!) as configured in ArgoCD where ArgoCD should deploy the rendered resources
func ApplicationsFromPlay(play *config.Play, project string, namespace string, server string) ([]Application, error) {
	// https://github.com/argoproj/argo-cd/blob/master/pkg/apis/application/v1alpha1/types.go
	result := []Application{}
	for _, chart := range play.Charts {
		app := Application{}
		// global Application resource settings
		app.APIVersion = "argoproj.io/v1alpha1"
		app.Kind = "Application"
		app.ObjectMeta.Namespace = namespace
		// TODO: Implement a proper approach to avoid argocd application name collisions.
		//       All argocd application resources for all clusters exist in the same namespace
		//       on the cluster hosting argocd. Aside from a proper name generation this requires
		//       a collision detection because the config structure of kusible allows for a setup
		//       that leads to non-preventable name collisions.
		app.ObjectMeta.Name = fmt.Sprintf("%s.%s.%s", chart.Name, project, server)
		app.Spec.Project = project

		// helm chart settings
		for _, repo := range play.Repos {
			if repo.Name == chart.Repo {
				app.Spec.Source.RepoURL = repo.URL
			}
		}

		if app.Spec.Source.RepoURL == "" {
			return result, fmt.Errorf("no repo '%s' for chart '%s' configured in play", chart.Repo, chart.Name)
		}

		app.Spec.Source.Chart = chart.Chart
		app.Spec.Source.TargetRevision = chart.Version

		// design decision: only support helm 3
		app.Spec.Source.Helm = &ApplicationSourceHelm{}
		app.Spec.Source.Helm.Version = "v3"
		values, err := yaml.Marshal(chart.Values)
		if err != nil {
			return result, fmt.Errorf("failed to convert values of chart '%s' to yaml: %s", chart.Name, err)
		}
		app.Spec.Source.Helm.Values = string(values)

		// target cluster + namespace settings
		app.Spec.Destination.Namespace = chart.Namespace
		app.Spec.Destination.Name = server

		result = append(result, app)
	}
	return result, nil
}
