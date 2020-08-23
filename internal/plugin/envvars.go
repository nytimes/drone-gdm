//==============================================================================
//
// drone-gdm/drone.go: This module import drone environment parameters
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
	"os"
	"strings"
)

type MapStrStr map[string]string

func DroneVars() MapStrStr {
	return varsFromEnvPrefix("DRONE_")
}

func PluginVars() MapStrStr {
	return varsFromEnvPrefix("PLUGIN_")
}

func varsFromEnvPrefix(prefix string) MapStrStr {
	envVars := make(MapStrStr)
	for _, e := range os.Environ() {
		comps := strings.SplitN(e, "=", 2)
		if (len(comps) > 1) && strings.HasPrefix(comps[0], prefix) {
			// TODO: Move this func to another translation unit + add filter
			vname := strings.ToLower(strings.TrimPrefix(comps[0], prefix))
			envVars[vname] = comps[1]
		}
	}
	return envVars
}
