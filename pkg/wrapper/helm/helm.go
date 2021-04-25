/*
Copyright © 2021 Bedag Informatik AG & The Helm Authors

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
	"fmt"
	"os"

	"github.com/bedag/kusible/pkg/inventory"

	"github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
)

// New returns a new helm instance
func New(options Options, settings *cli.EnvSettings, log logrus.FieldLogger) (*Helm, error) {
	h := &Helm{
		settings:   settings,
		out:        os.Stdout,
		helmDriver: os.Getenv("HELM_DRIVER"),
		options:    options,
		log:        log,
	}

	return h, nil
}

// NewWithGetter returns a new helm instance that uses the provided getter to
// retrieve kubeconfigs
func NewWithGetter(options Options, settings *cli.EnvSettings, getter *inventory.Kubeconfig, log logrus.FieldLogger) (*Helm, error) {
	h, err := New(options, settings, log)
	if err != nil {
		return nil, err
	}
	h.restClientGetter = getter
	return h, nil
}

func (h *Helm) ActionConfig(namespace string) (*action.Configuration, error) {
	// Set namespace acording to chart namespace in kubeconfig
	h.restClientGetter.SetNamespace(namespace)

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(h.restClientGetter, namespace, h.helmDriver, h.debug); err != nil {
		return nil, err
	}
	return actionConfig, nil
}

func (h *Helm) debug(format string, v ...interface{}) {
	if h.settings.Debug {
		format = fmt.Sprintf("[helm-debug] %s", format)
		h.log.Debugf(format, v...)
	}
}

func (h *Helm) getChartPathOptions(c *action.ChartPathOptions) {
	c.Verify = h.options.Verify
	c.Keyring = h.options.Keyring
}
