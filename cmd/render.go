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

package cmd

import (
	"fmt"

	"github.com/bedag/kusible/pkg/inventory"
	"github.com/bedag/kusible/pkg/playbook"
	"github.com/bedag/kusible/pkg/target"
	argocdutil "github.com/bedag/kusible/pkg/wrapper/argocd"
	"github.com/bedag/kusible/pkg/wrapper/ejson"
	helmutil "github.com/bedag/kusible/pkg/wrapper/helm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	helmcli "helm.sh/helm/v3/pkg/cli"
)

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render an application as kubernetes manifests",
}

var renderHelmCmd = &cobra.Command{
	Use:   "helm [playbook]",
	Short: "Use helm to render manifests for an inventory entry",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		playbookFile := args[0]
		limits := viper.GetStringSlice("limit")
		groupVarsDir := viper.GetString("group-vars-dir")
		inventoryPath := viper.GetString("inventory")
		skipEval := viper.GetBool("skip-eval")
		skipClusterInv := viper.GetBool("render-helm-skip-cluster-inventory")
		skipDecrypt := viper.GetBool("skip-decrypt")
		ejsonPrivKey := viper.GetString("ejson-privkey")
		ejsonKeyDir := viper.GetString("ejson-key-dir")

		invEjsonSettings := ejson.Settings{
			PrivKey: ejsonPrivKey,
			KeyDir:  ejsonKeyDir,
			// if we want to retrieve the cluster inventory ConfigMap
			// we need a kubeconfig to retrieve it, so we cannot skip
			// the decryption of the inventory settings
			SkipDecrypt: false,
		}

		tgtEjsonSettings := ejson.Settings{
			PrivKey:     ejsonPrivKey,
			KeyDir:      ejsonKeyDir,
			SkipDecrypt: skipDecrypt,
		}

		// if we do not retrieve the cluster inventory ConfigMap, we do not need to retrieve
		// the kubeconfig
		inventory, err := inventory.NewInventory(inventoryPath, invEjsonSettings, false)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to compile inventory.")
		}

		targets, err := target.NewTargets(".*", limits, groupVarsDir, inventory, true, &tgtEjsonSettings)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to compile values for inventory entries.")
		}

		playbookSet, err := playbook.NewSet(playbookFile, targets, skipEval, skipClusterInv)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to compile playbooks.")
		}

		settings := helmcli.New()

		for name, playbook := range playbookSet {
			for _, play := range playbook.Config.Plays {
				for _, repo := range play.Repos {
					if err := helmutil.RepoAdd(repo.Name, repo.URL, settings); err != nil {
						log.WithFields(log.Fields{
							"play":  play.Name,
							"repo":  repo.Name,
							"entry": name,
							"error": err.Error(),
						}).Fatal("Failed to add helm repo for play.")
					}
				}
				manifests, err := helmutil.TemplatePlay(play, settings)
				if err != nil {
					log.WithFields(log.Fields{
						"play":  play.Name,
						"entry": name,
						"error": err.Error(),
					}).Fatal("Failed to render play manifests with helm.")
				}
				fmt.Printf(manifests)
			}
		}
	},
}

var renderArgoCDCmd = &cobra.Command{
	Use:   "argocd [playbook]",
	Short: "Use render the given playbook into a set of ArgoCD Application resources",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		playbookFile := args[0]
		limits := viper.GetStringSlice("limit")
		groupVarsDir := viper.GetString("group-vars-dir")
		inventoryPath := viper.GetString("inventory")
		skipEval := viper.GetBool("skip-eval")
		skipClusterInv := viper.GetBool("render-argocd-skip-cluster-inventory")
		namespace := viper.GetString("argocd-namespace")
		project := viper.GetString("argocd-project")
		skipDecrypt := viper.GetBool("skip-decrypt")
		ejsonPrivKey := viper.GetString("ejson-privkey")
		ejsonKeyDir := viper.GetString("ejson-key-dir")

		invEjsonSettings := ejson.Settings{
			PrivKey: ejsonPrivKey,
			KeyDir:  ejsonKeyDir,
			// if we want to retrieve the cluster inventory ConfigMap
			// we need a kubeconfig to retrieve it, so we cannot skip
			// the decryption of the inventory settings
			SkipDecrypt: false,
		}

		tgtEjsonSettings := ejson.Settings{
			PrivKey:     ejsonPrivKey,
			KeyDir:      ejsonKeyDir,
			SkipDecrypt: skipDecrypt,
		}

		// if we do not retrieve the cluster inventory ConfigMap, we do not need to retrieve
		// the kubeconfig
		inventory, err := inventory.NewInventory(inventoryPath, invEjsonSettings, false)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to compile inventory.")
		}

		targets, err := target.NewTargets(".*", limits, groupVarsDir, inventory, true, &tgtEjsonSettings)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to compile values for inventory entries.")
		}

		playbookSet, err := playbook.NewSet(playbookFile, targets, skipEval, skipClusterInv)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("Failed to compile playbooks.")
		}

		for name, playbook := range playbookSet {
			for _, play := range playbook.Config.Plays {
				manifests, err := argocdutil.ApplicationFromPlay(play, project, namespace, name)
				if err != nil {
					log.WithFields(log.Fields{
						"play":  play.Name,
						"entry": name,
						"error": err.Error(),
					}).Fatal("Failed to render ArgoCD application manifests.")
				}
				fmt.Printf(manifests)
			}
		}
	},
}

func init() {
	renderHelmCmd.Flags().BoolP("skip-cluster-inventory", "", false, "Skip downloading the cluster-inventory ConfigMap")
	viper.BindPFlag("render-helm-skip-cluster-inventory", renderHelmCmd.Flags().Lookup("skip-cluster-inventory"))

	renderArgoCDCmd.Flags().BoolP("skip-cluster-inventory", "", false, "Skip downloading the cluster-inventory ConfigMap")
	renderArgoCDCmd.Flags().StringP("namespace", "", "argocd", "Namespace where ArgoCD is looking for ArgoCD applications")
	renderArgoCDCmd.Flags().StringP("project", "", "default", "The ArgoCD project to which the applications should be assigned")
	viper.BindPFlag("render-argocd-skip-cluster-inventory", renderArgoCDCmd.Flags().Lookup("skip-cluster-inventory"))
	viper.BindPFlag("argocd-namespace", renderArgoCDCmd.Flags().Lookup("namespace"))
	viper.BindPFlag("argocd-project", renderArgoCDCmd.Flags().Lookup("project"))

	renderCmd.AddCommand(renderHelmCmd)
	renderCmd.AddCommand(renderArgoCDCmd)
	rootCmd.AddCommand(renderCmd)
}
