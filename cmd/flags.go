/*
Copyright © 2021 Michael Gruener & Simon Fuhrer

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

import "github.com/spf13/cobra"

// addOutputFlags adds format and fields flags to a command.
func addOutputFlags(cmd *cobra.Command) {
	cmd.Flags().StringSlice("fields", []string{}, "Fields to print")
	cmd.Flags().String("format", "yaml", "Format to print (single,json,yaml,table)")
	cmd.Flags().Bool("list-wrap-single", false, "Wrap the result in an 'items: []' list even if it has only one element")
}

// addEjsonFlags adds flags to control ejson decryption
func addEjsonFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("ejson-privkey", "k", "", "EJSON private key")
	cmd.Flags().String("ejson-key-dir", "/opt/ejson/keys", "Directory containing EJSON keys")
	cmd.Flags().Bool("skip-decrypt", false, "Skip ejson decryption")
}

// addEvalFlags adds flags that controls spruce eval behavior
func addEvalFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("skip-eval", false, "Skip spruce operator evaluation")
}

func addLimitFlags(cmd *cobra.Command) {
	cmd.Flags().StringSliceP("limit", "l", []string{}, "Limit selected groups")
}

func addSkipClusterInventoryFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("skip-cluster-inventory", false, "Skip downloading the cluster-inventory ConfigMap")
}

func addClusterInventoryDefaultsFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("cluster-inventory-namespace", "c", "kube-system", "Default config namespace for the cluster inventory config map")
	cmd.Flags().String("cluster-inventory-configmap", "cluster-inventory", "Name of the cluster inventory config map in the cluster inventory namespace")
}

func addClusterInventoryFlags(cmd *cobra.Command) {
	addClusterInventoryDefaultsFlags(cmd)
	addSkipClusterInventoryFlags(cmd)
}
