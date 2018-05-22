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
	"bytes"
	"fmt"
	"gopkg.in/Masterminds/sprig.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"strings"
	"text/template"
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
		fmt.Sprintf("--filter=name=%s", spec.Name),
		fmt.Sprintf("--project=%s", context.Project),
	}...)

	result := RunGcloud(context, args...)
	if result.Error != nil {
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
	var err error
	var options []string
	var properties string

	if action != "delete" {
		if err != nil {
			return options, err
		}
	}

	fileOption, configPath, err := command.getFileOptions(context, spec, action)
	if err != nil {
		return options, err
	}

	switch action {
	case "create":
		addOptIfPresent(&options, fileOption, configPath)
		addOptIfPresent(&options, "--description", spec.Description)
		labels := mapAsOptions(spec.Labels, "=", ",")
		addOptIfPresent(&options, "--labels", labels)
		addOptIfPresent(&options, "--properties", properties)
		if spec.AutoRollback {
			options = append(options, "--automatic-rollback-on-error")
		}

	case "update":
		addOptIfPresent(&options, fileOption, configPath)
		addOptIfPresent(&options, "--description", spec.Description)
		addOptIfPresent(&options, "--properties", properties)
		addOptIfPresent(&options, "--create-policy", spec.CreatePolicy)
		addOptIfPresent(&options, "--delete-policy", spec.DeletePolicy)
		labels := mapAsOptions(spec.Labels, "=", ",")
		addOptIfPresent(&options, "--update-labels", labels)

	case "delete":
		addOptIfPresent(&options, "--delete-policy", spec.DeletePolicy)
	}
	return options, nil
}

func (command *GdmDeploymentCmd) getFileOptions(context *GdmPluginContext, spec *GdmConfigurationSpec, action string) (string, string, error) {
	if action == "delete" {
		return "", "", nil
	}

	pathSpecs := []struct {
		param string
		val   string
	}{
		{"path", spec.Path},
		{"config", spec.Config},
		{"template", spec.Template},
	}

	var configParam string
	var configPath string
	for _, pathSpec := range pathSpecs {
		if pathSpec.val != "" {
			if configPath == "" {
				configParam = pathSpec.param
				configPath = getAdjustedPath(pathSpec.val, context.Dir)
			} else {
				return "", "", fmt.Errorf(
					"Exactly one of \"path\", \"config\", or \"template\" is required for \"%s\". Got: \"%s: %s\" but already had \"%s: %s\"",
					spec.Group, pathSpec.param, pathSpec.val, configParam, configPath)
			}
		}
	}

	if configPath == "" {
		return "", "", fmt.Errorf(
			"Exactly one of \"path\", \"config\", or \"template\" is required for \"%s\"",
			spec.Group)
	}

	var err error
	var configOption string
	switch configParam {
	// Compatibility: for "path" parameter, determine option by extension
	case "path":
		if strings.HasSuffix(configPath, ".yml") || strings.HasSuffix(configPath, ".yaml") {
			configOption = "--config"
		} else {
			configOption = "--template"
		}
	case "config":
		configPath, err = command.getConfigFile(context, spec, configPath)
		configOption = "--config"
	case "template":
		configOption = "--template"
	}
	return configOption, configPath, err
}

func (command *GdmDeploymentCmd) getConfigFile(context *GdmPluginContext, spec *GdmConfigurationSpec, configPath string) (string, error) {
	t := template.New(path.Base(configPath))
	t.Funcs(template.FuncMap{
		"yaml": func(i interface{}) (string, error) {
			data, err := yaml.Marshal(i)
			return string(data), err
		},
	}).Funcs(sprig.GenericFuncMap())

	t, err := t.ParseFiles(configPath)
	if err != nil {
		return configPath, fmt.Errorf("Failed to parse configuration yaml", err)
	}

	var buff bytes.Buffer
	configVars := make(map[string]interface{})
	configVars["drone"] = DroneVars()
	configVars["vars"] = context.Vars
	configVars["properties"] = spec.Properties
	err = t.Execute(&buff, configVars)
	if err != nil {
		return configPath, err
	}

	tmpFile, err := ioutil.TempFile(context.TempDir(), "gdm-config")
	if err == nil {
		_, err = tmpFile.Write(buff.Bytes())
		if err == nil {
			tmpFile.Close()
		}
	}
	return tmpFile.Name(), err
}

func (command *GdmDeploymentCmd) getProperties(spec *GdmConfigurationSpec, action string) (string, error) {
	var properties string
	noProp := len(spec.Properties)
	if spec.PassAction {
		noProp += 1
	}

	if noProp > 0 {
		var propPairs []string
		for k, v := range spec.Properties {
			propVal, err := Y2JMarshal(v)
			if err != nil {
				return properties, err
			}
			propPairs = append(propPairs, fmt.Sprintf("%s:%s", k, propVal))
		}

		if spec.PassAction {
			propPairs = append(propPairs, fmt.Sprintf("action:%s", action))
		}

		properties = strings.Join(propPairs, ",")
	}
	return properties, nil
}

// EOF
