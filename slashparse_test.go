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
//  - [ ] return description and arguments 
// - [ ] newSlashCommand should return errors

func TestNewSlashCommand(t *testing.T) {

	testYamlPath := "./examples/helloWorld/simple.yaml"

	args := []string{"/print"}
	newSlash := NewSlashCommand(args, testYamlPath)

	want := SlashCommand{
		Name: "Print",
		Description: "Echos back what you type.",
		Arguments: []Argument{
			{
			Name: "text", 
			Description: "text you want to print",
			},
		},
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

* text: text you want to print
`

	assert.Equal(t, want, got)
}