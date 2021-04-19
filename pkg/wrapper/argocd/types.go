/*
Copyright © 2021 Bedag Informatik AG & The ArgoCD Authors

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

package argocd

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

/*
Because of https://github.com/argoproj/argo-cd/issues/4055 it is very cumbersome
(if at all possible in a sensible way) to import the ArgoCD go packages. According
to some comments in the issue it is easier to define the required datastructures
yourself, which we do here.

Prepare for lots of copy & paste from https://github.com/argoproj/argo-cd/blob/master/pkg/apis/application/v1alpha1/types.go

The definitions of the data structures are stripped down to what the kusible code
actually needs to keep the amount of required code maintenance down.
*/

// Application is a definition of an ArgoCD Application resource.
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ApplicationSpec `json:"spec"`
}

// ApplicationSpec represents desired application state. Contains link to repository with application definition and additional parameters link definition revision.
type ApplicationSpec struct {
	// Source is a reference to the location ksonnet application definition
	Source ApplicationSource `json:"source"`
	// Destination overrides the kubernetes server and namespace defined in the environment ksonnet app.yaml
	Destination ApplicationDestination `json:"destination"`
	// Project is a application project name. Empty name means that application belongs to 'default' project.
	Project string `json:"project"`
	// SyncPolicy controls when a sync will be performed
	SyncPolicy *SyncPolicy `json:"syncPolicy,omitempty"`
	// IgnoreDifferences controls resources fields which should be ignored during comparison
	IgnoreDifferences []ResourceIgnoreDifferences `json:"ignoreDifferences,omitempty"`
	// Infos contains a list of useful information (URLs, email addresses, and plain text) that relates to the application
	Info []Info `json:"info,omitempty"`
	// This limits this number of items kept in the apps revision history.
	// This should only be changed in exceptional circumstances.
	// Setting to zero will store no history. This will reduce storage used.
	// Increasing will increase the space used to store the history, so we do not recommend increasing it.
	// Default is 10.
	RevisionHistoryLimit *int64 `json:"revisionHistoryLimit,omitempty"`
}

// ApplicationSource contains information about github repository, path within repository and target application environment.
type ApplicationSource struct {
	// RepoURL is the repository URL of the application manifests
	RepoURL string `json:"repoURL"`
	// TargetRevision defines the commit, tag, or branch in which to sync the application to.
	// If omitted, will sync to HEAD
	TargetRevision string `json:"targetRevision,omitempty"`
	// Helm holds helm specific options
	Helm *ApplicationSourceHelm `json:"helm,omitempty"`
	// Chart is a Helm chart name
	Chart string `json:"chart,omitempty"`
}

// ApplicationDestination contains deployment destination information
type ApplicationDestination struct {
	// Namespace overrides the environment namespace value in the ksonnet app.yaml
	Namespace string `json:"namespace,omitempty"`
	// Name of the destination cluster which can be used instead of server (url) field
	Name string `json:"name,omitempty"`
}

// ResourceIgnoreDifferences contains resource filter and list of json paths which should be ignored during comparison with live state.
type ResourceIgnoreDifferences struct {
	Group        string   `json:"group,omitempty"`
	Kind         string   `json:"kind"`
	Name         string   `json:"name,omitempty"`
	Namespace    string   `json:"namespace,omitempty"`
	JSONPointers []string `json:"jsonPointers"`
}

// SyncPolicy controls when a sync will be performed in response to updates in git
type SyncPolicy struct {
	// Automated will keep an application synced to the target revision
	Automated *SyncPolicyAutomated `json:"automated,omitempty"`
	// Options allow you to specify whole app sync-options
	SyncOptions SyncOptions `json:"syncOptions,omitempty"`
	// Retry controls failed sync retry behavior
	Retry *RetryStrategy `json:"retry,omitempty"`
}

type Info struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ApplicationSourceHelm holds helm specific options
type ApplicationSourceHelm struct {
	// ValuesFiles is a list of Helm value files to use when generating a template
	ValueFiles []string `json:"valueFiles,omitempty"`
	// Parameters are parameters to the helm template
	Parameters []HelmParameter `json:"parameters,omitempty"`
	// The Helm release name. If omitted it will use the application name
	ReleaseName string `json:"releaseName,omitempty"`
	// Values is Helm values, typically defined as a block
	Values string `json:"values,omitempty"`
	// FileParameters are file parameters to the helm template
	FileParameters []HelmFileParameter `json:"fileParameters,omitempty"`
	// Version is the Helm version to use for templating with
	Version string `json:"version,omitempty"`
}

// HelmParameter is a parameter to a helm template
type HelmParameter struct {
	// Name is the name of the helm parameter
	Name string `json:"name,omitempty"`
	// Value is the value for the helm parameter
	Value string `json:"value,omitempty"`
	// ForceString determines whether to tell Helm to interpret booleans and numbers as strings
	ForceString bool `json:"forceString,omitempty"`
}

// HelmFileParameter is a file parameter to a helm template
type HelmFileParameter struct {
	// Name is the name of the helm parameter
	Name string `json:"name,omitempty"`
	// Path is the path value for the helm parameter
	Path string `json:"path,omitempty"`
}

// SyncPolicyAutomated controls the behavior of an automated sync
type SyncPolicyAutomated struct {
	// Prune will prune resources automatically as part of automated sync (default: false)
	Prune bool `json:"prune,omitempty"`
	// SelfHeal enables auto-syncing if  (default: false)
	SelfHeal bool `json:"selfHeal,omitempty"`
	// AllowEmpty allows apps have zero live resources (default: false)
	AllowEmpty bool `json:"allowEmpty,omitempty"`
}

type SyncOptions []string

type RetryStrategy struct {
	// Limit is the maximum number of attempts when retrying a container
	Limit int64 `json:"limit,omitempty"`

	// Backoff is a backoff strategy
	Backoff *Backoff `json:"backoff,omitempty"`
}

// Backoff is a backoff strategy to use within retryStrategy
type Backoff struct {
	// Duration is the amount to back off. Default unit is seconds, but could also be a duration (e.g. "2m", "1h")
	Duration string `json:"duration,omitempty"`
	// Factor is a factor to multiply the base duration after each failed retry
	Factor *int64 `json:"factor,omitempty"`
	// MaxDuration is the maximum amount of time allowed for the backoff strategy
	MaxDuration string `json:"maxDuration,omitempty"`
}
