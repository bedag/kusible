/*
Copyright © 2021 Michael Gruener

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

package helm

import (
	"os"
	"time"

	"github.com/bedag/kusible/pkg/inventory"
	"github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/cli"
)

// Helm encabsules a single helm instance
type Helm struct {
	settings         *cli.EnvSettings
	out              *os.File
	helmDriver       string
	restClientGetter *inventory.Kubeconfig
	options          Options
	log              logrus.FieldLogger
}

// Options holds all (relevant) helm cli options
type Options struct {
	CreateNamespace          bool
	DryRun                   bool
	NoHooks                  bool
	Replace                  bool
	Timeout                  time.Duration
	Wait                     bool
	WaitForJobs              bool
	DepdencyUpdate           bool
	DisableOpenAPIValidation bool
	Atomic                   bool
	SkipCRDs                 bool
	RenderSubChartNotes      bool
	Verify                   bool
	Keyring                  string
	Validate                 bool
	IncludeCRDs              bool
	ExtraAPIs                []string
	Force                    bool
	ResetValues              bool
	ReuseValues              bool
	HistoryMax               int
	CleanupOnFail            bool
	KeepHistory              bool
}
