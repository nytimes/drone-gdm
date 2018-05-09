//==============================================================================
//
// drone-gdm/plugin/context.go: Plugin context for drone-gdm
//
// Copyright (c) 2017 The New York Times Company
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this library except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//------------------------------------------------------------------------------

package plugin

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

// This is the context object we'll use to store the internal state used
// by the drone-gdm to translate yaml parameters into one or more gcloud
// invocations.
//
// (For information on the "drone" struct tags, see parse.go)
type GdmPluginContext struct {
	// drone:
	Dir         string `drone:"env=DRONE_DIR"`
	Repo        string `drone:"env=DRONE_REPO"`
	Branch      string `drone:"env=DRONE_BRANCH"`
	Commit      string `drone:"env=DRONE_COMMIT"`
	BuildNumber string `drone:"env=DRONE_BUILD_NUMBER"`
	PullRequest string `drone:"env=DRONE_PULL_REQUEST"`
	JobNumber   string `drone:"env=DRONE_JOB_NUMBER"`
	Tag         string `drone:"env=DRONE_TAG"`

	// drone-gdm:
	GcloudPath string            `drone:"env=PLUGIN_GCLOUDPATH"`
	Verbose    bool              `drone:"env=PLUGIN_VERBOSE"`
	DryRun     bool              `drone:"env=PLUGIN_DRYRUN"`
	Vars       map[string]string `drone:"env=PLUGIN_VARS"`

	// gcloud:
	Token   string `drone:"env=TOKEN"`
	Project string `drone:"env=PLUGIN_PROJECT"`

	// deployment-manager:
	Preview        bool                   `drone:"env=PLUGIN_PREVIEW"`
	Async          bool                   `drone:"env=PLUGIN_ASYNC"`
	ConfigFile     string                 `drone:"env=PLUGIN_CONFIGFILE"`
	Configurations []GdmConfigurationSpec `drone:"env=PLUGIN_CONFIGURATIONS"`

	// Internal use only:
	parseOkay bool

	// TODO: other gcloud global parameters
}

// Return a pointer to a new GdmPluginContext
func NewGdmPluginContext() *GdmPluginContext {
	return &GdmPluginContext{
		Preview:    false,
		GcloudPath: "/google-cloud-sdk/bin/gcloud",
		Async:      false,
		parseOkay:  false,
	}
}

// Validate GdmPluginContext, after parsing.
func (context *GdmPluginContext) Validate() error {
	if !context.parseOkay {
		return fmt.Errorf("context invalid (not parsed)")
	}

	if context.Token == "" {
		return fmt.Errorf("\"token\" is a required parameter")
	} else {
		// Be sure to trim any dangling whitespace off the token!
		context.Token = strings.TrimSpace(context.Token)
	}

	if context.Project == "" {
		return fmt.Errorf("\"project\" is a required parameter")
	}

	if len(context.Configurations) == 0 {
		return fmt.Errorf("at least one configuration must be present")
	}

	for _, spec := range context.Configurations {
		err := spec.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (context *GdmPluginContext) loadConfigurations() error {
	if context.ConfigFile == "" {
		return nil
	}

	t, err := template.ParseFiles(context.ConfigFile)
	if err != nil {
		return err
	}

	return t.Execute(os.Stdout, context.Vars)
}

// Parse GdmPluginContext using ParsePluginParams.
// Set parseOkay flag accordingly.
func (context *GdmPluginContext) Parse() error {
	err := ParsePluginParams(context)
	if err != nil {
		return err
	}

	err = context.loadConfigurations()
	if err != nil {
		return err
	}

	context.parseOkay = (err == nil)
	return err
}

// EOF
