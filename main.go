//==============================================================================
//
// drone-gdm/main.go: Drone plugin for Google Deployment Manager
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
	"github.com/nytimes/drone-gdm/plugin"
	drone "github.com/drone/drone-plugin-go/plugin"
	"os"
	"io/ioutil"
)

//------------------------------------
// Globals:
//------------------------------------
const gdmTokenPath string = "/tmp/gcloud.json"

var lbl string = "[unknown]"
var rev string = "[unknown]"

const debug bool = true

//------------------------------------
// Main:
//------------------------------------
func errBail(err error) {
	fmt.Printf("ERROR: %s\n", err)
	cleanupToken()
	os.Exit(1)
}

// drone-gdm plugin entry point.
// Actions performed:
// - Parses plugin parameters from environment
// - Validates parsed plugin parameters
// - Executes google deployment manager, via gcloud
func main() {
	fmt.Printf("Drone GDM Plugin %s - built from %s:\n", lbl, rev)

	context := plugin.NewGdmPluginContext()

	// https://godoc.org/github.com/drone/drone-plugin-go/plugin
	if len(os.Args) > 1 {
		workspace := drone.Workspace{}
		drone.Param("workspace", &workspace)
		drone.Param("vargs", &context)
		drone.Parse()
		context.Dir = workspace.Path
	}

	err := context.Parse()
	if err != nil {
		errBail(err)
	}

	err = context.Validate()
	if err != nil {
		errBail(err)
	}

	err = performTokenAuthentication(context)
	if err != nil {
		errBail(err)
	}

	for _,spec := range context.Configurations {
		err = plugin.GdmExecute(context, &spec, gdmTokenPath)
		if err != nil {
			errBail(err)
		}
	}

	cleanupToken()
	os.Exit(0)
}

func performTokenAuthentication(context *plugin.GdmPluginContext) error {
	// Write credentials to tmp file to be picked up by the 'gcloud' command.
	// This is inside the ephemeral plugin container, not on the host:
	err := ioutil.WriteFile(gdmTokenPath, []byte(context.Token), 0600)
	if err != nil {
		return fmt.Errorf("error writing token file: %s\n", err)
	}

	// Ensure the token is cleaned up, no matter exit status:
	err = plugin.ActivateServiceAccount(context, gdmTokenPath)
	return err
}

func cleanupToken() {
	err := os.Remove(gdmTokenPath)
	if err != nil {
		// No need to panic on error, due to likely ephemeral mount
		fmt.Printf("drone-gdm: WARNING: error removing token file: %s\n", err)
	}
}

// EOF
