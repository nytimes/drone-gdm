//==============================================================================
//
// drone-gdm/config.go: Configuration file context for drone-gdm
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
)

type GdmConfigurationSpec struct {
	Vars          map[string]interface{} `json:"vars" yaml:"vars"`
	GdmVersion    string                 `json:"version" yaml:"version" yaml:"version" yaml:"version"`
	Group         string                 `json:"group" yaml:"group"`
	State         string                 `json:"state" yaml:"state"`
	Name          string                 `json:"name" yaml:"name"`
	Path          string                 `json:"path" yaml:"path"`
	Config        string                 `json:"config" yaml:"config"`
	Template      string                 `json:"template" yaml:"template"`
	Description   string                 `json:"description" yaml:"description"`
	DescriptorURL string                 `json:"descriptorUrl" yaml:"descriptorUrl"`
	APIOptions    string                 `json:"apiOptions" yaml:"apiOptions"`
	Labels        map[string]string      `json:"labels" yaml:"labels"`
	Properties    map[string]interface{} `json:"properties" yaml:"properties"`
	AutoRollback  bool                   `json:"automaticRollbackOnError" yaml:"automaticRollbackOnError"`
	CreatePolicy  string                 `json:"createPolicy" yaml:"createPolicy"`
	DeletePolicy  string                 `json:"deletePolicy" yaml:"deletePolicy"`
	Status        string                 `json:"status" yaml:"status"`
	PassAction    bool                   `json:"passAction" yaml:"passAction"`
}

func (spec *GdmConfigurationSpec) Validate() error {
	var err error

	// version
	err = IsParamInRange("version", spec.GdmVersion, "alpha", "beta")
	if spec.GdmVersion != "" && err != nil {
		return fmt.Errorf("configuration error: %v", err)
	}

	err = IsParamInRange("group", spec.Group, "deployment", "composite", "typeprovider")
	if err != nil {
		return fmt.Errorf("configuration error: %v", err)
	}

	// (drone-gdm) state
	err = IsParamInRange("state", spec.State, "latest", "present", "absent")
	if err != nil {
		return fmt.Errorf("configuration error: %v", err)
	}

	// name
	if spec.Name == "" {
		return fmt.Errorf("\"name\" is a required parameter")
	}

	// input file
	// path --> --config OR --template ; optional for delete-only
	// --description: (optional)
	// --labels: (optional)
	// --properties: (optional)
	// --autorollback: (optional)

	// --create-policy (optional)
	if spec.CreatePolicy != "" {
		err = IsParamInRange("createPolicy", spec.CreatePolicy, "ACQUIRE", "CREATE_OR_ACQUIRE")
		if spec.CreatePolicy != "" && err != nil {
			return fmt.Errorf("configuration error: %v", err)
		}
	}

	// --delete-policy (optional)
	if spec.DeletePolicy != "" {
		err = IsParamInRange("deletePolicy", spec.DeletePolicy, "ABANDON", "DELETE")
		if spec.DeletePolicy != "" && err != nil {
			return fmt.Errorf("configuration error: %v", err)
		}
	}

	// --status (optional)
	if spec.Status != "" {
		err = IsParamInRange("status", spec.Status, "DEPRECATED", "EXPERIMENTAL", "SUPPORTED")
		if spec.Status != "" && err != nil {
			return fmt.Errorf("configuration error: %v", err)
		}
	}

	// --update-labels: (optional)
	return nil
}

// EOF
