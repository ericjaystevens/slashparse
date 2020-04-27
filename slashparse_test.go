package slashparse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO next
// - [x] parse yaml from file
// - [ ] Print Help
// 	- [x] Print command name
//  - [x] Give Yaml usefull names 
// - [ ] newSlashCommand should return errors

func TestNewSlashCommand(t *testing.T) {

	testYamlPath := "./examples/helloWorld/simple.yaml"

	args := []string{"/print"}
	newSlash := NewSlashCommand(args, testYamlPath)

	want := SlashCommand{
		name: "Print",
		description:"Echos back what you type.",
	}
	assert.Equal(t, want, newSlash)
}

func TestGetSlashHelp(t *testing.T){

	testYamlPath := "./examples/helloWorld/simple.yaml"

	args := []string{"/print"}
	newSlash := NewSlashCommand(args, testYamlPath)

	got := newSlash.GetSlashHelp()

	want := `## Print Help
* Echos back what you type. *

### Arguments
`

	assert.Equal(t, want, got)
}