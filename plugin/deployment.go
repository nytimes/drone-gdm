//==============================================================================
//
// drone-gdm/plugin/deploy.go: GDM logic for "Deployments"
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
	"strings"
)

type GdmDeploymentCmd struct{}

func NewGdmDeploymentCmd() *GdmDeploymentCmd {
	return &GdmDeploymentCmd{}
}

func (command *GdmDeploymentCmd) Name() string {
	return "deployments"
}

// On success, return a bool indicating whether or not the named deployment
// already exists.
func (command *GdmDeploymentCmd) Exists(context *GdmPluginContext, spec *GdmConfigurationSpec) (bool, error) {
	fmt.Printf("drone-gdm: Checking for existing %s \"%s\"...\n", spec.Group, spec.Name)
	var deployExists bool

	gcmd := getCmdPrelude(command, spec)
	args := append(gcmd, []string{
		"list",
		"--simple-list",
		fmt.Sprintf("--filter=name:%s", spec.Name),
		fmt.Sprintf("--project=%s", context.Project),
	}...)

	result := RunGcloud(context, args...)
	if !result.Okay {
		return deployExists, fmt.Errorf("error listing deployments: %s\n", result.Stderr.String())
	}

	if strings.TrimSpace(result.Stdout.String()) == spec.Name {
		fmt.Printf("drone-gdm: \"%s\" exists\n", spec.Name)
		deployExists = true
	} else {
		fmt.Printf("drone-gdm: \"%s\" does not exist\n", spec.Name)
		deployExists = false
	}
	return deployExists, nil
}

func (command *GdmDeploymentCmd) Action(spec *GdmConfigurationSpec, exists bool) (string, error) {
	if exists {
		switch spec.State {
		case "latest":
			return "update", nil
		case "absent":
			return "delete", nil
		}
	} else {
		switch spec.State {
		case "latest":
			fallthrough
		case "present":
			return "create", nil
		case "absent":
			return "", nil
		}
	}

	// Any other combo results in "no action":
	return "", nil
}

func (command *GdmDeploymentCmd) Options(context *GdmPluginContext, spec *GdmConfigurationSpec, action string) ([]string, error) {
	var options []string
	var properties string

	if action != "delete" {
		noProp := len(spec.Properties)
		if spec.PassAction {
			noProp += 1
		}

		if noProp > 0 {
			i := 0
			var propPairs []string
			for k, v := range spec.Properties {
				propPairs = append(propPairs, fmt.Sprintf("%s:%v", k, v))
				i++
			}

			if spec.PassAction {
				propPairs = append(propPairs, fmt.Sprintf("action:%s", action))
			}

			properties = strings.Join(propPairs, ",")
		}
	}

	configPath := getAdjustedPath(spec.Path, context.Dir)
	switch action {
	case "create":
		options = addOptIfPresent(options, configPath, "--config")
		options = addOptIfPresent(options, spec.Description, "--description")
		labels := mapAsOptions(spec.Labels, "=", ",")
		options = addOptIfPresent(options, labels, "--labels")
		options = addOptIfPresent(options, properties, "--properties")
		if spec.AutoRollback {
			options = append(options, "--automatic-rollback-on-error")
		}

	case "update":
		options = addOptIfPresent(options, configPath, "--configPath")
		options = addOptIfPresent(options, spec.Description, "--description")
		options = addOptIfPresent(options, properties, "--properties")
		options = addOptIfPresent(options, spec.CreatePolicy, "--create-policy")
		options = addOptIfPresent(options, spec.DeletePolicy, "--delete-policy")
		labels := mapAsOptions(spec.Labels, "=", ",")
		options = addOptIfPresent(options, labels, "--update-labels")

	case "delete":
		options = addOptIfPresent(options, spec.DeletePolicy, "--delete-policy")
	}
	return options, nil
}

// EOF
