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

package main

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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
	Dir string `drone:"env=DRONE_DIR"`

	// drone-gdm:
	Debug      bool                   `drone:"env=PLUGIN_DEBUG"`
	GcloudPath string                 `drone:"env=PLUGIN_GCLOUDPATH"`
	Verbose    bool                   `drone:"env=PLUGIN_VERBOSE"`
	DryRun     bool                   `drone:"env=PLUGIN_DRYRUN"`
	Vars       map[string]interface{} `drone:"env=PLUGIN_VARS"`

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
	tmpDir    string

	// TODO: other gcloud global parameters
}

// Return a pointer to a new GdmPluginContext
func NewGdmPluginContext() (*GdmPluginContext, error) {
	tmpDir, err := ioutil.TempDir("", "drone-gdm")
	if err != nil {
		return nil, fmt.Errorf("Unable to create gdmPluginContext: %s", err)
	}

	return &GdmPluginContext{
		// drone-gdm:
		Debug:      false,
		GcloudPath: "/google-cloud-sdk/bin/gcloud",
		Verbose:    false,
		DryRun:     false,

		// deployment-manager:
		Preview: false,
		Async:   false,

		// internal:
		parseOkay: false,
		tmpDir:    tmpDir,
	}, nil
}

func (context *GdmPluginContext) TempDir() string {
	return context.tmpDir
}

// Parse GdmPluginContext using ParsePluginParams.
// Set parseOkay flag accordingly.
func (context *GdmPluginContext) Parse() error {
	err := ParsePluginParams(context)

	if err != nil {
		// Ensure the debug flag is set appropriatley, regardless of parse
		// status so that main can leverage this for error reporting:
		pluginVars := PluginVars()
		dbg, ok := pluginVars["debug"]
		if ok && dbg == "true" {
			context.Debug = true
		}
		return err
	}

	err = context.loadConfigurations()
	if err != nil {
		return err
	}

	context.parseOkay = (err == nil)
	return err
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

func (context *GdmPluginContext) Authenticate() error {
	tokenFile, err := ioutil.TempFile(context.tmpDir, "gdm-token.json")
	if err != nil {
		return fmt.Errorf("error creating token file: %s", err)
	}

	// Write credentials to tmp file to be picked up by the 'gcloud' command.
	// This is inside the ephemeral plugin container, not on the host:
	_, err = tokenFile.Write([]byte(context.Token))
	if err != nil {
		return fmt.Errorf("error writing token file: %s\n", err)
	}

	return context.ActivateServiceAccount(tokenFile.Name())
}

func (context *GdmPluginContext) ActivateServiceAccount(gdmTokenPath string) error {
	fmt.Println("drone-gdm: Authenticating")
	args := []string{
		"auth",
		"activate-service-account",
		"--key-file",
		gdmTokenPath,
	}

	result := RunGcloud(context, args...)
	if result.Error != nil {
		return fmt.Errorf("error activating service account: %s\n", result.Stderr.String())
	}
	return nil
}

func (context *GdmPluginContext) Cleanup() error {
	if context == nil {
		return nil
	}

	if context.Debug {
		fmt.Println("\x1b[01;33mEnvironment Variables Defined for this Build:\x1b[00m")

		for k, v := range PluginVars() {
			fmt.Printf("\t\t\x1b[00;33m%s: \x1b[00;34m%s\x1b[00m\n", k, v)
		}
	}

	var err error
	if (context.tmpDir != "") && (!context.Debug) {
		err = os.RemoveAll(context.tmpDir)
	}
	return err
}

func (context *GdmPluginContext) loadConfigurations() error {
	if context.ConfigFile == "" {
		return nil
	}

	t, err := template.ParseFiles(getAdjustedPath(context.ConfigFile, context.Dir))
	if err != nil {
		return err
	}

	tmplVars := make(map[string]interface{})
	tmplVars["drone"] = DroneVars()
	tmplVars["plugin"] = PluginVars()
	tmplVars["context"] = context.Vars

	var buff bytes.Buffer
	err = t.Execute(&buff, tmplVars)
	if err != nil {
		return err
	}

	var configurations []GdmConfigurationSpec
	err = yaml.Unmarshal(buff.Bytes(), &configurations)
	if err != nil {
		return err
	}

	if len(configurations) != 0 {
		context.Configurations = append(context.Configurations, configurations...)
	}
	return nil
}

// EOF
