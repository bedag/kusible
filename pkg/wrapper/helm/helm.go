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

package helm

import (
	"fmt"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// New returns a new helm instance
func New(options Options, settings *cli.EnvSettings) (*Helm, error) {
	h := &Helm{
		settings:   settings,
		out:        os.Stdout,
		helmDriver: os.Getenv("HELM_DRIVER"),
		options:    options,
	}
	h.restClientGetter = h.settings.RESTClientGetter()

	return h, nil
}

// NewWithGetter returns a new helm instance that uses the provided getter to
// retrieve kubeconfigs
func NewWithGetter(options Options, settings *cli.EnvSettings, getter genericclioptions.RESTClientGetter) (*Helm, error) {
	h, err := New(options, settings)
	if err != nil {
		return nil, err
	}
	h.restClientGetter = getter
	return h, nil
}

func (h *Helm) ActionConfig(namespace string) (*action.Configuration, error) {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(h.restClientGetter, namespace, h.helmDriver, h.debug); err != nil {
		return nil, err
	}
	return actionConfig, nil
}

// STUB, just here to be passed to action.Configration.Init()
// TODO: implement proper debug log method
func (h *Helm) debug(format string, v ...interface{}) {
	if h.settings.Debug {
		format = fmt.Sprintf("[debug] %s\n", format)
		fmt.Printf(format, v...)
	}
}

func (h *Helm) getChartPathOptions(c *action.ChartPathOptions) {
	c.Verify = h.options.Verify
	c.Keyring = h.options.Keyring
}
