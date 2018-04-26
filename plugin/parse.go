//==============================================================================
//
// drone-gdm/plugin/parse.go: Generic-ish drone plugin env parser for go
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
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// Parse and validate values received from drone.
//
// This function leverages tags in the format:
//     `drone:"[opt1,opt2...optN]"`
//
// Where options have the format:
//     opt[=val]
//
// Available options:
//  * env - indicates this value should be aqcuired from given env var
//  * alt - indicates this value can be acquired from an alternate field
//  * required - indicates the field is required
//
// Unrecognized options are ignored.
// See: http://readme.drone.io/0.5/usage/environment-reference/ for more info.
func ParsePluginParams(context interface{}) error {
	// TODO: I'm certain there's a better way...
	v := reflect.Indirect(reflect.ValueOf(context))

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := v.Type().Field(i).Tag.Get("drone")

		// Skip tag if empty/ignored
		if tag == "" || tag == "-" {
			continue
		}

		spec := newPluginParseSpec(tag)
		val, err := getParamValue(field, spec)
		if err != nil {
			return err
		}

		if val == "" {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(val)

		default:
			jsonIn := reflect.New(field.Type()).Interface()
			err := json.Unmarshal([]byte(val), jsonIn)
			if err != nil {
				return err
			}

			jsonVal := reflect.ValueOf(jsonIn)
			if jsonIn != nil {
				field.Set(reflect.Indirect(jsonVal))
			}
		}
	}

	return nil
}

// Check that the given string parameter falls within a range of expected
// values.
func IsParamInRange(pName string, pVal string, okVals ...string) error {
	validStates := make(map[string]struct{}, len(okVals))

	for _, v := range okVals {
		validStates[v] = struct{}{}
	}

	if _, ok := validStates[pVal]; !ok {
		return fmt.Errorf("%s must be in %v (got \"%s\")", pName, okVals, pVal)
	}
	return nil
}

//------------------------------------
// Utility:
//------------------------------------

// Keep track of the parse status for a given field.
type pluginParseSpec struct {
	envName  string
	required bool
	alt      string
	format   string
}

// New function for pluginParseSpec
func newPluginParseSpec(tag string) *pluginParseSpec {
	p := new(pluginParseSpec)
	values := strings.Split(tag, ",")

	// Set the options:
	for _, val := range values {
		pair := strings.Split(val, "=")
		switch pair[0] {
		case "env":
			p.envName = pair[1]
		case "required":
			p.required = true
		case "alt":
			p.alt = pair[1]
		}
	}
	return p
}

// Get parameter value from environment:
func getParamValue(field reflect.Value, spec *pluginParseSpec) (string, error) {
	var val string
	var err error

	if spec.envName != "" {
		val = os.Getenv(spec.envName)
	}

	// If no value, try alternate - if present:
	if val == "" && spec.alt != "" {
		val = os.Getenv(spec.alt)
	}

	return val, err
}

// EOF
