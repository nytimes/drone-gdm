//==============================================================================
//
// drone-gdm/plugin/run.go: Utilities for executing gcloud
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
	"os/exec"
	"strings"
)

// Used to keep track of gcloud invocation results
type GcloudResult struct {
	Stdout bytes.Buffer
	Stderr bytes.Buffer
	Okay   bool
}

const escape = "\x1b"

func verbose(context *GdmPluginContext, fspec string, args ...interface{}) {
	if context.Verbose {
		fmt.Printf(fspec, args...)
	}
}

// Function used to run gcloud
func RunGcloud(context *GdmPluginContext, args ...string) *GcloudResult {
	var qualifier string
	if context.DryRun {
		qualifier = " (dry run)"
	}

	verbose(context, "drone-gdm%s:\n\t\"\x1b[34m%s\x1b[0m \x1b[32m%s\x1b[0m\"\n",
		qualifier, context.GcloudPath,
		strings.Join(args, " \\\n\t\t"))

	command := exec.Command(context.GcloudPath, args...)
	result := bindResult(command)
	if !context.DryRun {
		err := command.Run()

		if err == nil {
			result.Okay = true
		}
	} else {
		result.Okay = true
	}

	verbose(context, "\tStatus Okay: \x1b[33m%v\x1b[0m\n", result.Okay)
	return result
}

//------------------------------------
// Utility:
//------------------------------------

// Create a GcloudResult and bind it to the given command.
func bindResult(command *exec.Cmd) *GcloudResult {
	g := new(GcloudResult)

	command.Stdout = &g.Stdout
	command.Stderr = &g.Stderr
	g.Okay = false

	return g
}

// EOF
