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
	drone "github.com/drone/drone-plugin-go/plugin"
	"github.com/nytimes/drone-gdm/internal/plugin"
	"os"
)

var context *plugin.GdmPluginContext
var lbl string = "[label unset]"
var rev string = "[revision unset]"

// drone-gdm plugin entry point.
func main() {
	fmt.Printf("Drone GDM Plugin: %s (build revision: %s)\n", lbl, rev)

	var err error
	context, err = plugin.NewGdmPluginContext()
	if err != nil {
		errBail(err)
	}

	// https://godoc.org/github.com/drone/drone-plugin-go/plugin
	if len(os.Args) > 1 {
		workspace := drone.Workspace{}
		drone.Param("workspace", &workspace)
		drone.Param("vargs", &context)
		drone.Parse()
		context.Dir = workspace.Path
	}

	err = context.Parse()
	if err != nil {
		errBail(err)
	}

	err = context.Validate()
	if err != nil {
		errBail(err)
	}

	err = context.Authenticate()
	if err != nil {
		errBail(err)
	}

	for _, spec := range context.Configurations {
		err = plugin.GdmExecute(context, &spec)
		if err != nil {
			errBail(err)
		}
	}

	doCleanup()
	os.Exit(0)
}

func errBail(err error) {
	fmt.Printf("\x1b[00;31mERROR: %s\n\x1b[00m", err)
	doCleanup()
	os.Exit(1)
}

func doCleanup() {
	err := context.Cleanup()
	if err != nil {
		// No need to panic on error; (likely ephemeral mount disappeared)
		fmt.Printf("drone-gdm: WARNING: cleanup failed with: %s\n", err)
	}
}

// EOF
