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
	"github.com/bedag/kusible/pkg/printer"
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

	printerQueue := printer.Queue{}
	for name, playbook := range playbookSet {
		playbookMap, err := playbook.Map(skipEval)
		// see https://golang.org/doc/faq#closures_and_goroutines
		name := name
		if err != nil {
			log.WithFields(log.Fields{
				"entry": name,
				"error": err.Error(),
			}).Error("Failed to convert playbook entry to a map.")
			return err
		}

		if len(playbookMap) > 0 {
			job := printer.NewJob(func(fields []string) map[string]interface{} {
				defaultResult := map[string]interface{}{
					"entry":    name,
					"playbook": playbookMap,
				}

				// the output should not be limited to specific fields, just
				// render the whole playbook for each entry
				if len(fields) < 1 {
					return defaultResult
				}

				// Iterate over each play and just render the requested fields
				// for each play. This means the "fields" parameter refers to
				// the fields of each play of the playbook and not the fields of
				// the playbook itself (which would just be one field: "plays")
				if plays, ok := playbookMap["plays"]; ok {
					resultPlays := []map[string]interface{}{}
					for _, p := range plays.([]interface{}) {
						play := p.(map[string]interface{})
						resultPlay := map[string]interface{}{}
						for _, field := range fields {
							if val, ok := play[field]; ok {
								resultPlay[field] = val
							}
						}
						resultPlays = append(resultPlays, resultPlay)
					}
					return map[string]interface{}{
						"entry": name,
						"playbook": map[string]interface{}{
							"plays": resultPlays,
						},
					}
				}

				// the playbook did not have a list of plays for whatever
				// reason, just render the default result (which should be
				// the whole playbook)
				return defaultResult
			})
			printerQueue = append(printerQueue, job)
		}
	}

	return c.output(printerQueue)
}
