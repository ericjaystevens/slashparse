package slashparse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSlashCommand(t *testing.T) {

	testYamlPath := "./examples/helloWorld/helloWorldCommand.yaml"

	args := []string{"/print"}
	newSlash := NewSlashCommand(args, testYamlPath)

	want := "works on /print so Easy!"
	assert.Equal(t, want, newSlash)
}
