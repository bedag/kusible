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
func New() (*Helm, error) {
	h := &Helm{
		settings:     cli.New(),
		out:          os.Stdout,
		actionConfig: new(action.Configuration),
		helmDriver:   os.Getenv("HELM_DRIVER"),
	}

	if err := h.actionConfig.Init(h.settings.RESTClientGetter(), h.settings.Namespace(), h.helmDriver, h.debug); err != nil {
		return nil, err
	}
	return h, nil
}

func NewWithGetter(getter genericclioptions.RESTClientGetter) (*Helm, error) {
	h, err := New()
	if err != nil {
		return h, err
	}

	if err := h.actionConfig.Init(getter, h.settings.Namespace(), h.helmDriver, h.debug); err != nil {
		return nil, err
	}
	return h, nil
}

// STUB, just here to be passed to action.Configration.Init()
// TODO: implement proper debug log method
func (h *Helm) debug(format string, v ...interface{}) {
	if h.settings.Debug {
		format = fmt.Sprintf("[debug] %s\n", format)
		fmt.Printf(format, v...)
	}
}
