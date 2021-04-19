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
	"fmt"

	"github.com/bedag/kusible/pkg/printer"
	"github.com/imdario/mergo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/google/uuid"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func newInventoryKubeconfigCmd(c *Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "kubeconfig [filter]",
		Short:                 "Get the kubeconfig for one or more inventory entries",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  c.wrap(runInventoryKubeconfig),
	}
	addInventoryFlags(cmd)

	return cmd
}

func runInventoryKubeconfig(c *Cli, cmd *cobra.Command, args []string) error {
	filter := args[0]
	limits := c.viper.GetStringSlice("limit")

	inv, err := getInventoryWithKubeconfig(c)
	if err != nil {
		return err
	}

	names, err := inv.EntryNames(filter, limits)
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to get list of entries")
		return err
	}

	kubeconfigs := []*clientcmdapi.Config{}
	for _, name := range names {
		entry := inv.Entries()[name]
		clientConfig, err := entry.Kubeconfig().Config()
		if err != nil {
			c.Log.WithFields(logrus.Fields{
				"entry": name,
				"error": err.Error(),
			}).Error("Failed to get kubeconfig")
			return err
		}

		config, err := clientConfig.RawConfig()
		if err != nil {
			c.Log.WithFields(logrus.Fields{
				"entry": name,
				"error": err.Error(),
			}).Error("Failed to get kubeconfig")
			return err
		}

		// make all cluster/user/context names unique to
		// prevent collisions when merging with other entry
		// kubeconfigs
		makeUnique(&config)
		kubeconfigs = append(kubeconfigs, &config)
	}

	kubeconfig := mergeKubeconfigs(kubeconfigs)
	data, err := clientcmd.Write(*kubeconfig)
	if err != nil {
		c.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to render merged kubeconfigs as yaml")
		return err
	}

	printFn := func(fields []string) map[string]interface{} {
		var defaultResult map[string]interface{}

		yaml.Unmarshal(data, &defaultResult)
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
	}

	printerQueue := printer.Queue{printer.NewJob(printFn)}
	return c.output(printerQueue)
}

func mergeKubeconfigs(kubeconfigs []*clientcmdapi.Config) *clientcmdapi.Config {
	// copy of https://github.com/kubernetes/client-go/blob/ab82d40f6e857a3162e22ac8a5888b6314f9b0eb/tools/clientcmd/loader.go#L225
	mapConfig := clientcmdapi.NewConfig()

	for _, kubeconfig := range kubeconfigs {
		mergo.Merge(mapConfig, kubeconfig, mergo.WithOverride)
	}

	// merge all of the struct values in the reverse order so that priority is given correctly
	// errors are not added to the list the second time
	nonMapConfig := clientcmdapi.NewConfig()
	for i := len(kubeconfigs) - 1; i >= 0; i-- {
		kubeconfig := kubeconfigs[i]
		mergo.Merge(nonMapConfig, kubeconfig, mergo.WithOverride)
	}

	// since values are overwritten, but maps values are not, we can merge the non-map config on top of the map config and
	// get the values we expect.
	config := clientcmdapi.NewConfig()
	mergo.Merge(config, mapConfig, mergo.WithOverride)
	mergo.Merge(config, nonMapConfig, mergo.WithOverride)

	for name := range config.Contexts {
		config.CurrentContext = name
		break
	}

	return config
}

func makeUnique(config *clientcmdapi.Config) {
	uuidSuffix := uuid.New().String()
	authInfos := make(map[string]*clientcmdapi.AuthInfo, len(config.AuthInfos))
	clusters := make(map[string]*clientcmdapi.Cluster, len(config.Clusters))
	contexts := make(map[string]*clientcmdapi.Context, len(config.Contexts))

	for n, authInfo := range config.AuthInfos {
		name := fmt.Sprintf("%s-%s", n, uuidSuffix)
		authInfos[name] = authInfo
	}
	config.AuthInfos = authInfos

	for n, cluster := range config.Clusters {
		name := fmt.Sprintf("%s-%s", n, uuidSuffix)
		clusters[name] = cluster
	}
	config.Clusters = clusters

	for n, context := range config.Contexts {
		name := fmt.Sprintf("%s-%s", n, uuidSuffix)
		context.AuthInfo = fmt.Sprintf("%s-%s", context.AuthInfo, uuidSuffix)
		context.Cluster = fmt.Sprintf("%s-%s", context.Cluster, uuidSuffix)
		contexts[name] = context
	}
	config.Contexts = contexts
}
