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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewRootCommand(c *Cli) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:           "kusible",
		Short:         "Render and deploy kubernetes resources",
		Long:          `This is a CLI tool to render and deploy kubernetes resources to multiple clusters`,
		SilenceUsage:  false,
		SilenceErrors: true,
	}

	rootCmd.PersistentFlags().StringP("log-level", "", log.WarnLevel.String(), "log level (trace,debug,info,warn/warning,error,fatal,panic)")
	rootCmd.PersistentFlags().BoolP("json-log", "", false, "log as json")
	rootCmd.PersistentFlags().BoolP("log-functions", "", false, "log function names (performance impact!)")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress all normal output")

	c.bindAllFlags(rootCmd)

	rootCmd.AddCommand(
		newVersionCmd(),
		newValuesCmd(c),
		newGroupsCmd(c),
		newRenderCmd(c),
		newInventoryCmd(c),
	)

	return rootCmd
}
