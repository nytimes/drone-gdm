//==============================================================================
//
// drone-gdm/plugin/yaml2json.go: JSON marshalling overrides to facilitate easy
//                                conversion of yaml data to json.
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
	"encoding/json"
	"fmt"
)

// Encode an object into JSON which may have initially been non-compliant due
// to non-string keys in the top-level or nested objects.
func Y2JMarshal(i interface{}) ([]byte, error) {
	return json.Marshal(y2jConvert(i))
}

// Take an input interface which may have been populated by YAML decoder (i.e.
// a structure which may contain non-string keys or have nested objects with
// non-string keys) and return a clone of the object with the key types updated
// to be JSON compliant.
func y2jConvert(i interface{}) interface{} {
	switch i.(type) {
	case map[interface{}]interface{}:
		return y2jMap(i)
	case []interface{}:
		return y2jList(i)
	default:
		return i
	}
}

// Convert YAML maps to JSON maps
func y2jMap(yMap interface{}) interface{} {
	jMap := make(map[string](interface{}))
	for k, v := range yMap.(map[interface{}]interface{}) {
		jMap[fmt.Sprintf("%v", k)] = y2jConvert(v)
	}
	return jMap
}

// Convert YAML lists to JSON lists
func y2jList(yList interface{}) interface{} {
	var jList []interface{}
	for _, i := range yList.([]interface{}) {
		jList = append(jList, y2jConvert(i))
	}
	return jList
}
