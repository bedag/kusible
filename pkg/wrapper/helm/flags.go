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

/*
Kusible needs to "emulate" the helm cli for all commands requireing helm.
The "Globals" struct serves more or less the same role as Cobra in the actual
helm: each helm command (install, template ...) uses the Globals struct
to retrieve cli options which are then used to populate the different
"action" (action.Install...) structs
*/

package helm

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"helm.sh/helm/v3/pkg/cli"
	"k8s.io/client-go/util/homedir"
)

func NewOptions(viper *viper.Viper) Options {
	return Options{
		CreateNamespace:          viper.GetBool("helm-create-namespace"),
		NoHooks:                  viper.GetBool("helm-no-hooks"),
		Replace:                  viper.GetBool("helm-replace"),
		Timeout:                  viper.GetDuration("helm-timeout"),
		Wait:                     viper.GetBool("helm-wait"),
		WaitForJobs:              viper.GetBool("helm-wait-for-jobs"),
		DepdencyUpdate:           viper.GetBool("helm-dependency-update"),
		DisableOpenAPIValidation: viper.GetBool("helm-disable-openapi-validation"),
		Atomic:                   viper.GetBool("helm-atomic"),
		SkipCRDs:                 viper.GetBool("helm-skip-crds"),
		RenderSubChartNotes:      viper.GetBool("helm-render-subchart-notes"),
		Verify:                   viper.GetBool("helm-verify"),
		Keyring:                  viper.GetString("helm-keyring"),
		Validate:                 viper.GetBool("helm-validate"),
		IncludeCRDs:              viper.GetBool("helm-include-crds"),
		ExtraAPIs:                viper.GetStringSlice("helm-api-versions"),
		Force:                    viper.GetBool("helm-force"),
		ResetValues:              viper.GetBool("helm-reset-values"),
		ReuseValues:              viper.GetBool("helm-reuse-values"),
		HistoryMax:               viper.GetInt("helm-history-max"),
		CleanupOnFail:            viper.GetBool("helm-cleanup-on-fail"),
	}
}

func AddHelmInstallFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("helm-create-namespace", true, "create the release namespace if not present")
	cmd.Flags().Bool("helm-no-hooks", false, "prevent hooks from running during install")
	cmd.Flags().Bool("helm-replace", false, "re-use the given name, only if that name is a deleted release which remains in the history. This is unsafe in production")
	cmd.Flags().Duration("helm-timeout", 300*time.Second, "time to wait for any individual Kubernetes operation (like Jobs for hooks)")
	cmd.Flags().Bool("helm-wait", false, "if set, will wait until all Pods, PVCs, Services, and minimum number of Pods of a Deployment, StatefulSet, or ReplicaSet are in a ready state before marking the release as successful. It will wait for as long as --helm-timeout")
	cmd.Flags().Bool("helm-wait-for-jobs", false, "if set and --helm-wait enabled, will wait until all Jobs have been completed before marking the release as successful. It will wait for as long as --helm-timeout")
	cmd.Flags().Bool("helm-dependency-update", false, "run helm dependency update before installing the chart")
	cmd.Flags().Bool("helm-disable-openapi-validation", false, "if set, the installation process will not validate rendered templates against the Kubernetes OpenAPI Schema")
	cmd.Flags().Bool("helm-atomic", false, "if set, the installation process deletes the installation on failure. The --helm-wait flag will be set automatically if --helm-atomic is used")
	cmd.Flags().Bool("helm-skip-crds", false, "if set, no CRDs will be installed. By default, CRDs are installed if not already present")
	cmd.Flags().Bool("helm-render-subchart-notes", false, "if set, render subchart notes along with the parent")
}

func AddHelmUpgradeFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("helm-force", false, "force resource updates through a replacement strategy")
	cmd.Flags().Bool("helm-reset-values", false, "when upgrading, reset the values to the ones built into the chart")
	cmd.Flags().Bool("helm-reuse-values", false, "when upgrading, reuse the last release's values and merge in any overrides from the command line via --set and -f. If '--reset-values' is specified, this is ignored")
	// not sure if instantiating a new EnvSettings here just to get the default is a good idea
	cmd.Flags().Int("helm-history-max", cli.New().MaxHistory, "limit the maximum number of revisions saved per release. Use 0 for no limit")
	cmd.Flags().Bool("helm-cleanup-on-fail", false, "allow deletion of new resources created in this upgrade when upgrade fails")
}

func AddHelmChartPathOptionsFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("helm-verify", false, "verify the package before using it")
	cmd.Flags().String("helm-keyring", defaultKeyring(), "location of public keys used for verification")
}

func AddHelmTemplateFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("helm-validate", false, "validate your manifests against the Kubernetes cluster you are currently pointing at. This is the same validation performed on an install")
	cmd.Flags().Bool("helm-include-crds", false, "include CRDs in the templated output")
	cmd.Flags().StringArrayP("helm-api-versions", "a", []string{}, "Kubernetes api versions used for Capabilities.APIVersions")
}

// defaultKeyring returns the expanded path to the default keyring.
func defaultKeyring() string {
	if v, ok := os.LookupEnv("GNUPGHOME"); ok {
		return filepath.Join(v, "pubring.gpg")
	}
	return filepath.Join(homedir.HomeDir(), ".gnupg", "pubring.gpg")
}
