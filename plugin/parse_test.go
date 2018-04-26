package plugin

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// Basic:
type reqStruct struct {
	Arg string `drone:"env=TEST_ARG,required"`
}

func TestRequiredPresent(t *testing.T) {
	os.Setenv("TEST_ARG", "argval")
	os.Setenv("TEST_OPT", "optval")

	myObj := reqStruct{}

	err := ParsePluginParams(&myObj)
	if assert.Nil(t, err) {
		assert.Equal(t, "argval", myObj.Arg)
	}
}

// Booleans:
type flagStruct struct {
	Flag0 bool `drone:"env=TEST_FLAG0,required"`
	Flag1 bool `drone:"env=TEST_FLAG1,required"`
	Flag2 bool `drone:"env=TEST_FLAG2"`
}

func TestBool(t *testing.T) {
	os.Setenv("TEST_FLAG0", "true")
	os.Setenv("TEST_FLAG1", "false")

	myObj := flagStruct{}

	err := ParsePluginParams(&myObj)
	if assert.Nil(t, err) {
		assert.True(t, myObj.Flag0)
		assert.False(t, myObj.Flag1)
		assert.False(t, myObj.Flag2)
	}
}

// Integers:
type intStruct struct {
	Flag0 int `drone:"env=TEST_FLAG0,required"`
	Flag1 int `drone:"env=TEST_FLAG1,required"`
	Flag2 int `drone:"env=TEST_FLAG2"`
}

func TestInteger(t *testing.T) {
	os.Setenv("TEST_FLAG0", "0")
	os.Setenv("TEST_FLAG1", "1")

	myObj := intStruct{}

	err := ParsePluginParams(&myObj)
	if assert.Nil(t, err) {
		assert.Equal(t, myObj.Flag0, 0)
		assert.Equal(t, myObj.Flag1, 1)
	}
}

// map[string]string:
type mapStruct struct {
	Properties map[string]interface{} `drone:"env=TEST_PROPERTIES,required"`
	Options    map[string]interface{} `drone:"env=TEST_OPTIONS"`
}

func TestMap(t *testing.T) {
	jProp := "{\"Name\":\"Alice\",\"b\":true,\"i\":0}"

	os.Setenv("TEST_PROPERTIES", jProp)
	os.Setenv("TEST_OPTIONS", "")

	myObj := mapStruct{
		Properties: make(map[string]interface{}),
		Options:    make(map[string]interface{}),
	}

	err := ParsePluginParams(&myObj)
	if assert.Nil(t, err) {
		assert.Equal(t, "Alice", myObj.Properties["Name"])

		// Notice, these are strings!
		assert.Equal(t, true, myObj.Properties["b"])
		assert.Equal(t, 0.0, myObj.Properties["i"])
	}

	// This was not populated:
	assert.Empty(t, myObj.Options)
}

// map[string]string:
type gdmStruct struct {
	Configurations []GdmConfigurationSpec `drone:"env=TEST_CONFIGURATIONS"`
}

func TestGdmConfigurationSpec(t *testing.T) {
	os.Setenv("TEST_CONFIGURATIONS", "[{\"description\":\"Test desc\",\"properties\":{\"env\":\"test\"}}]")

	myObj := gdmStruct{}

	err := ParsePluginParams(&myObj)

	assert.Nil(t, err)

	assert.NotNil(t, myObj.Configurations)
	assert.Equal(t, "Test desc", myObj.Configurations[0].Description)
	assert.Equal(t, "test", myObj.Configurations[0].Properties["env"])
}
