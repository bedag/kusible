/*
Copyright © 2021 Copyright © 2021 Bedag Informatik AG

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
	"github.com/spf13/cobra"
)

func newRenderCmd(c *Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "render",
		Short:                 "Render an application as kubernetes manifests",
		Args:                  cobra.NoArgs,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
	}

	cmd.AddCommand(
		newRenderPlaybookCmd(c),
		newRenderHelmCmd(c),
		newRenderArgoCDCmd(c),
	)

	return cmd
}

func addRenderFlags(cmd *cobra.Command) {
	addGroupsFlags(cmd)
	addInventoryFlags(cmd)
	addSkipClusterInventoryFlags(cmd)
}
