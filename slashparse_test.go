package slashparse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSlashCommand(t *testing.T) {
	testYamlPath := "./examples/helloWorld/simple.yaml"

	args := []string{"/print"}
	newSlash := NewSlashCommand(args, testYamlPath)

	want := SlashCommand{
		Name:        "Print",
		Description: "Echos back what you type.",
		Arguments: []Argument{
			{
				Name:        "text",
				Description: "text you want to print",
			},
		},
	}
	assert.Equal(t, want, newSlash)
}

func TestGetSlashHelp(t *testing.T) {
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

type getCommandStringTests struct {
	testName string
	args     []string
	want     string
}

func TestGetCommandString(t *testing.T) {
	tests := []getCommandStringTests{
		{
			testName: "valid print example",
			args:     []string{"/print"},
			want:     "print",
		},
		{
			testName: "invalid print example",
			args:     []string{"/print"},
			want:     "",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			testYamlPath := "./examples/helloWorld/simple.yaml"
			newSlash := NewSlashCommand(test.args, testYamlPath)
			got := newSlash.GetCommandString()
			assert.Equal(t, test.want, got)
		})
	}
}
