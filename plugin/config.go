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
)

type GdmConfigurationSpec struct {
	GdmVersion   string `json:"version"`
	Group        string `json:"group"`
	State        string
	Name         string
	Path         string
	Description  string
	Labels       map[string]string
	Properties   map[string]interface{}
	AutoRollback bool `json:"automaticRollbackOnError"`
	CreatePolicy string
	DeletePolicy string
	Status       string
	PassAction   bool `json:"passAction"`
}

func (spec *GdmConfigurationSpec) Validate() error {
	var err error

	// version
	err = IsParamInRange("version", spec.GdmVersion, "alpha", "beta")
	if spec.GdmVersion != "" && err != nil {
		return fmt.Errorf("configuration error: %v", err)
	}

	err = IsParamInRange("group", spec.Group, "deployment", "composite")
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
		if spec.DeletePolicy != "" && err != nil {
			return fmt.Errorf("configuration error: %v", err)
		}
	}

	// --update-labels: (optional)
	return nil
}

// EOF
