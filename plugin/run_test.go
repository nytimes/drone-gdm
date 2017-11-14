package plugin

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestBasicRun(t *testing.T) {
	context := NewGdmPluginContext()
	context.GcloudPath = "/bin/echo"

	result := RunGcloud(context, "Open", "the", "pod", "bay", "doors,", "Hal")

	if assert.Equal(t, true, result.Error == nil) {
		cmdOut := strings.TrimSpace(result.Stdout.String())
		assert.Equal(t, "Open the pod bay doors, Hal", cmdOut)
	}
}

func TestMissingCmd(t *testing.T) {
	context := NewGdmPluginContext()
	context.GcloudPath = "/bin/plimbst"

	result := RunGcloud(context, "Open", "the", "pod", "bay", "doors,", "Hal")

	assert.Equal(t, false, result.Error == nil, "Function fails due to missing command")
}
