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
	"sort"
	"strings"

	"github.com/bedag/kusible/pkg/groups"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newGroupsCmd(c *Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "groups [regex]",
		Short:                 "List available groups based on given regex",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  c.wrap(runGroups),
	}
	addGroupsFlags(cmd)
	addLimitFlags(cmd)

	return cmd
}

func addGroupsFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("group-vars-dir", "d", "group_vars", "Source directory to read from")
}

func runGroups(c *Cli, cmd *cobra.Command, args []string) error {
	filter := args[0]
	limits := c.viper.GetStringSlice("limit")
	groupVarsDir := c.viper.GetString("group-vars-dir")

	groups, err := groups.Groups(groupVarsDir, filter, limits)
	if err != nil {
		log.WithFields(log.Fields{
			"filter":    filter,
			"limits":    strings.Join(limits[:], " "),
			"directory": groupVarsDir,
		}).Error("Failed to get groups")
		return err
	}

	sort.Strings(groups)
	fmt.Println(strings.Join(groups, "\n"))
	return nil
}
