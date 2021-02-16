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

	"encoding/json"

	"github.com/spf13/cobra"
)

const appName = "kusible"

// Version is the application version. It will be overriden during the build process.
// See https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
var Version = "development"

// BuildTime is the date/time when the binary was build. It should be set by the build process and will not be
// displayed if it is empty.
var BuildTime string

// GitRev is the Git revision of the code the binary is build from. It should be set by the build process and will not be
// displayed if it is empty.
var GitRev string

// GitTreeState indicates if the git tree had uncommited changes when the binary was build. It should be set by the build
// process and will not be displayed if it is empty.
var GitTreeState string

type BuildInfo struct {
	Version      string `json:"Version,omitempty"`
	GitRev       string `json:"GitRev,omitempty"`
	BuildTime    string `json:"BuildTime,omitempty"`
	GitTreeState string `json:"GitTreeState,omitempty"`
}

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: fmt.Sprint("Print the version number of ", appName),
		Run: func(cmd *cobra.Command, args []string) {
			bi := BuildInfo{
				Version:      Version,
				GitRev:       GitRev,
				BuildTime:    BuildTime,
				GitTreeState: GitTreeState,
			}

			version, _ := json.MarshalIndent(bi, "", "  ")
			fmt.Println(string(version))
		},
	}
	return cmd
}
