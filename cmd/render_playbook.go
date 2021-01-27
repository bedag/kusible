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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newRenderPlaybookCmd(c *Cli) *cobra.Command {
	var cmd = &cobra.Command{
		Use:                   "playbook [playbook]",
		Short:                 "Render the given playbook",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		RunE:                  c.wrap(runRenderPlaybook),
	}
	addRenderFlags(cmd)

	return cmd
}

func runRenderPlaybook(c *Cli, cmd *cobra.Command, args []string) error {
	playbookFile := args[0]
	skipEval := c.viper.GetBool("skip-eval")

	playbookSet, err := loadPlaybooks(c, playbookFile)
	if err != nil {
		return err
	}

	for name, playbook := range playbookSet {
		result, err := playbook.YAML(skipEval)
		if err != nil {
			log.WithFields(log.Fields{
				"entry": name,
				"error": err.Error(),
			}).Error("Failed to convert playbook entry to yaml.")
			return err
		}
		if len(result) > 0 {
			fmt.Printf("======= Plays for %s =======\n", name)
			fmt.Printf("%s", string(result))
		}
	}

	return nil
}
