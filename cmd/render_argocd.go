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

	"github.com/bedag/kusible/pkg/inventory"
	invconfig "github.com/bedag/kusible/pkg/inventory/config"
	"github.com/bedag/kusible/pkg/playbook"
	"github.com/bedag/kusible/pkg/target"
	argocdutil "github.com/bedag/kusible/pkg/wrapper/argocd"
	"github.com/bedag/kusible/pkg/wrapper/ejson"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
	cmd.Flags().StringP("argocd-namespace", "", "argocd", "Namespace where ArgoCD is looking for ArgoCD applications")
	cmd.Flags().StringP("argocd-project", "", "default", "The ArgoCD project to which the applications should be assigned")

	return cmd
}

func runRenderArgoCD(c *Cli, cmd *cobra.Command, args []string) error {
	playbookFile := args[0]
	limits := c.viper.GetStringSlice("limit")
	groupVarsDir := c.viper.GetString("group-vars-dir")
	inventoryPath := c.viper.GetString("inventory")
	skipEval := c.viper.GetBool("skip-eval")
	namespace := c.viper.GetString("argocd-namespace")
	project := c.viper.GetString("argocd-project")
	skipDecrypt := c.viper.GetBool("skip-decrypt")
	ejsonPrivKey := c.viper.GetString("ejson-privkey")
	ejsonKeyDir := c.viper.GetString("ejson-key-dir")
	skipClusterInv := c.viper.GetBool("skip-cluster-inventory")
	clusterInvNamespace := c.viper.GetString("cluster-inventory-namespace")
	clusterInvConfigMap := c.viper.GetString("cluster-inventory-configmap")

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

	clusterInventoryDefaults := invconfig.ClusterInventory{
		Namespace: clusterInvNamespace,
		ConfigMap: clusterInvConfigMap,
	}

	// if we do not retrieve the cluster inventory ConfigMap, we do not need to retrieve
	// the kubeconfig
	inventory, err := inventory.NewInventory(inventoryPath, invEjsonSettings, false, clusterInventoryDefaults)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to compile inventory.")
		return err
	}

	targets, err := target.NewTargets(".*", limits, groupVarsDir, inventory, true, &tgtEjsonSettings)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to compile values for inventory entries.")
		return err
	}

	playbookSet, err := playbook.NewSet(playbookFile, targets, skipEval, skipClusterInv)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to compile playbooks.")
		return err
	}

	for name, playbook := range playbookSet {
		for _, play := range playbook.Config.Plays {
			manifests, err := argocdutil.ApplicationFromPlay(play, project, namespace, name)
			if err != nil {
				log.WithFields(log.Fields{
					"play":  play.Name,
					"entry": name,
					"error": err.Error(),
				}).Error("Failed to render ArgoCD application manifests.")
				return err
			}
			fmt.Printf(manifests)
		}
	}
	return nil
}
