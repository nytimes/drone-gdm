//==============================================================================
//
// drone-gdm/plugin/composite.go: GDM logic for "composite types"
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

type GdmTypeProviderCmd struct{}

func NewTypeProviderCmd() *GdmTypeProviderCmd {
	return &GdmTypeProviderCmd{}
}

func (command *GdmTypeProviderCmd) Name() string {
	return "type-providers"
}

// On success, return a bool indicating whether or not the named deployment
// already exists.
func (command *GdmTypeProviderCmd) Exists(context *GdmPluginContext, spec *GdmConfigurationSpec) (bool, error) {
	fmt.Printf("drone-gdm: Checking for existing %s \"%s\"...\n", spec.Group, spec.Name)
	var deployExists bool

	gcmd := getCmdPrelude(command, spec)
	args := append(gcmd, []string{
		"list",
		fmt.Sprintf("--project=%s", context.Project),
		fmt.Sprintf("--filter=name=%s", spec.Name),
	}...)

	result := RunGcloud(context, args...)
	if result.Error != nil {
		return deployExists, fmt.Errorf("error listing type providers: %s\n", result.Stderr.String())
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

func (command *GdmTypeProviderCmd) Action(spec *GdmConfigurationSpec, exists bool) (string, error) {
	switch {
	case spec.State == "latest":
		return "", fmt.Errorf("\"latest\" is not a valid state for type providers")
	case exists && spec.State == "absent":
		return "delete", nil
	case !exists && spec.State == "present":
		return "create", nil
	case exists && spec.State == "present":
		return "update", nil
	}
	return "", nil
}

func (command *GdmTypeProviderCmd) Options(context *GdmPluginContext, spec *GdmConfigurationSpec, action string) ([]string, error) {
	var options []string

	if action != "delete" && spec.DescriptorURL == "" {
		return options, fmt.Errorf("\"descriptorURL\" is required for \"%s\"", spec.Group)
	}
	addOptIfPresent(&options, "--descriptor-url", spec.DescriptorURL)
	if (spec.APIOptions != "") {
		templatePath := getAdjustedPath(spec.APIOptions, context.Dir)
		addOptIfPresent(&options, "--api-options-file", templatePath)
	}
	return options, nil
}

// EOF
