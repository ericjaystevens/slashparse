package slashparse

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type newSlashCommandTests struct {
	testName      string
	args          []string
	want          SlashCommand
	configPath    string
	expectedError error
}

func TestNewSlashCommand(t *testing.T) {
	tests := []newSlashCommandTests{
		{
			testName:   "simple test",
			args:       []string{"/print"},
			configPath: "./examples/helloWorld/simple.yaml",
			want: SlashCommand{
				Name:        "Print",
				Description: "Echos back what you type.",
				Arguments: []Argument{
					{
						Name:        "text",
						Description: "text you want to print",
					},
				},
			},
		},
		{
			testName:      "invalid test",
			args:          []string{"/pssrint"},
			configPath:    "./examples/helloWorld/simple.yaml",
			expectedError: errors.New("pssrint is not a valid command"),
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			newSlash, err := NewSlashCommand(test.args, test.configPath)
			if err != nil {
				assert.Equal(t, test.expectedError, err)
			}
			assert.Equal(t, test.want, newSlash)
		})
	}
}

func TestGetSlashHelp(t *testing.T) {
	testYamlPath := "./examples/helloWorld/simple.yaml"

	args := []string{"/print"}
	newSlash, _ := NewSlashCommand(args, testYamlPath)

	got := newSlash.GetSlashHelp()

	want := `## Print Help
* Echos back what you type. *

### Arguments

* text: text you want to print
`
	assert.Equal(t, want, got)
}

type getCommandStringTests struct {
	testName    string
	args        []string
	want        string
	expectError bool
}

func TestGetCommandString(t *testing.T) {
	tests := []getCommandStringTests{
		{
			testName: "valid print example",
			args:     []string{"/print"},
			want:     "Print",
		},
		{
			testName:    "invalid print example",
			args:        []string{""},
			want:        "",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			testYamlPath := "./examples/helloWorld/simple.yaml"
			newSlash, _ := NewSlashCommand(test.args, testYamlPath)
			got, err := newSlash.GetCommandString(test.args)
			if err != nil {
				assert.Equal(t, test.expectError, true)
			} else {
				assert.Equal(t, test.want, got)
			}
		})
	}
}
