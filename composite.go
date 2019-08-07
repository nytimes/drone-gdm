//==============================================================================
//
// drone-gdm/composite.go: GDM logic for "composite types"
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
	"fmt"
	"strings"
)

type GdmCompositeCmd struct{}

func NewGdmCompositeCmd() *GdmCompositeCmd {
	return &GdmCompositeCmd{}
}

func (command *GdmCompositeCmd) Name() string {
	return "types"
}

// On success, return a bool indicating whether or not the named deployment
// already exists.
func (command *GdmCompositeCmd) Exists(context *GdmPluginContext, spec *GdmConfigurationSpec) (bool, error) {
	fmt.Printf("drone-gdm: Checking for existing %s \"%s\"...\n", spec.Group, spec.Name)
	var deployExists bool

	gcmd := getCmdPrelude(command, spec)
	args := append(gcmd, []string{
		"list",
		fmt.Sprintf("--project=%s", context.Project),
		fmt.Sprintf("--provider-project=%s", context.Project),
		fmt.Sprintf("--filter=types.name=%s", spec.Name),
	}...)

	result := RunGcloud(context, args...)
	if result.Error != nil {
		return deployExists, fmt.Errorf("error listing types: %s\n", result.Stderr.String())
	}

	if strings.TrimSpace(result.Stdout.String()) != "" {
		fmt.Printf("drone-gdm: \"%s\" exists\n", spec.Name)
		deployExists = true
	} else {
		fmt.Printf("drone-gdm: \"%s\" does not exist\n", spec.Name)
		deployExists = false
	}
	return deployExists, nil
}

func (command *GdmCompositeCmd) Action(spec *GdmConfigurationSpec, exists bool) (string, error) {
	if spec.State == "latest" {
		return "", fmt.Errorf("\"latest\" is not a valid state for composite types")
	}

	if exists && spec.State == "absent" {
		return "delete", nil
	}

	if !exists && spec.State == "present" {
		return "create", nil
	}

	// Any other combo results in "no action":
	return "", nil
}

func (command *GdmCompositeCmd) Options(context *GdmPluginContext, spec *GdmConfigurationSpec, action string) ([]string, error) {
	var options []string

	templatePath := getAdjustedPath(spec.Path, context.Dir)
	if action != "delete" && (spec.Path == "") && (spec.Template == "") {
		return options, fmt.Errorf("At least one of \"path\", \"config\", or \"template\" is required for \"%s\"", spec.Group)
	}

	switch action {
	case "create":
		addOptIfPresent(&options, "--template", templatePath)
		addOptIfPresent(&options, "--description", spec.Description)
		labels := mapAsOptions(spec.Labels, "=", ",")
		addOptIfPresent(&options, "--labels", labels)
		addOptIfPresent(&options, "--status", spec.Status)
	case "update":
		addOptIfPresent(&options, "--template", templatePath)
		addOptIfPresent(&options, "--description", spec.Description)
		labels := mapAsOptions(spec.Labels, "=", ",")
		addOptIfPresent(&options, "--update-labels", labels)
		addOptIfPresent(&options, "--status", spec.Status)
	}
	return options, nil
}

// EOF
