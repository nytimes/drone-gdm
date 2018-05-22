//==============================================================================
//
// drone-gdm/plugin/deploy.go: Central drone-gdm deployment logic
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
	"path/filepath"
	"strings"
)

type GdmCommand interface {
	Name() string
	Exists(context *GdmPluginContext, spec *GdmConfigurationSpec) (bool, error)
	Action(spec *GdmConfigurationSpec, exists bool) (string, error)
	Options(context *GdmPluginContext, spec *GdmConfigurationSpec, action string) ([]string, error)
}

// Run a complte GDM deployment, according to the parameters passed to the
// drone-gdm plugin:
//  - Store passed in credentials in a temp file
//  - Activate service account using temp file
//  - Check the present state of the deployment
//  - Execute a gdm command to transform present --> desired state
func GdmExecute(context *GdmPluginContext, spec *GdmConfigurationSpec) error {
	command := getGdmCommand(spec)
	if command == nil {
		return fmt.Errorf("\"%s\" is not a supported command", spec.Group)
	}

	exists, err := command.Exists(context, spec)
	if err != nil {
		return err
	}

	action, err := command.Action(spec, exists)
	if err != nil {
		return err
	}
	if action == "" {
		fmt.Println("drone-gdm: No action required")
		return nil
	}

	return executeDeploymentAction(context, spec, action, command)
}

//------------------------------------
// Utility:
//------------------------------------
func getGdmCommand(spec *GdmConfigurationSpec) GdmCommand {
	switch spec.Group {
	case "deployment":
		return NewGdmDeploymentCmd()
	case "composite":
		return NewGdmCompositeCmd()
	}
	return nil
}

// Execute the depoloyment manager action to transform present-->desired state.
func executeDeploymentAction(context *GdmPluginContext, spec *GdmConfigurationSpec, action string, command GdmCommand) error {
	fmt.Printf("drone-gdm: Performing \"%s\" action for deployment \"%s\"...\n", action, spec.Name)

	// Mandatory arguments:
	gcmd := getCmdPrelude(command, spec)
	args := append(gcmd, []string{
		action,
		spec.Name,
		fmt.Sprintf("--project=%s", context.Project),
	}...)

	// Command options:
	options, err := command.Options(context, spec, action)
	if err != nil {
		return err
	}

	// Context options:
	args = append(args, options...)
	if context.Async {
		args = append(args, "--async")
	}

	// Don't request user input for delete actions.
	if action == "delete" {
		args = append(args, "-q")
	}

	// Engage!
	result := RunGcloud(context, args...)
	if result.Error != nil {
		return fmt.Errorf("error performing \"%s\" action on \"%s\": %s\n", action, spec.Name, result.Stderr.String())
	}

	verbose(context, "Results: %s\n", result.Stdout.String())
	return nil
}

func getCmdPrelude(command GdmCommand, spec *GdmConfigurationSpec) []string {
	var cmd []string
	switch spec.Group {
	case "deployment":
		cmd = []string{"deployment-manager", command.Name()}
	case "composite":
		cmd = []string{"deployment-manager", command.Name()}
	}

	if cmd != nil && spec.GdmVersion != "" {
		cmd = append([]string{spec.GdmVersion}, cmd...)
	}
	return cmd
}

func mapAsOptions(optMap map[string]string, op string, sep string) string {
	var optArg string

	if optMap != nil {
		var opts []string
		for k, v := range optMap {
			opts = append(opts, fmt.Sprintf("%s%s%s", k, op, v))
		}
		optArg = strings.Join(opts, sep)
	}
	return optArg
}

func addOptIfPresent(options *[]string, optName string, val string) {
	if val != "" {
		*options = append(*options, fmt.Sprintf("%s=%s", optName, val))
	}
}

func getAdjustedPath(fileArg string, cwd string) string {
	var adjPath string
	if filepath.IsAbs(fileArg) || cwd == "" {
		adjPath = fileArg
	} else {
		adjPath = filepath.Join(cwd, fileArg)
	}
	return adjPath
}

// EOF
