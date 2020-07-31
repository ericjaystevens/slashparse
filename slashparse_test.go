package slashparse

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

type newSlashCommandTests struct {
	testName      string
	want          SlashCommand
	configPath    string
	expectedError error
}

func getSimpleDef() []byte {
	simpleDef, _ := ioutil.ReadFile("./examples/helloWorld/simple.yaml")
	return simpleDef
}

var SimpleDef = getSimpleDef()

func TestNewSlashCommand(t *testing.T) {
	tests := []newSlashCommandTests{
		{
			testName:   "simple test",
			configPath: "./examples/helloWorld/simple.yaml",
			want: SlashCommand{
				Name:        "Print",
				Description: "Echos back what you type.",
				Arguments: []Argument{
					{
						Name:        "text",
						ArgType:     "quoted text",
						Description: "text you want to print",
						ErrorMsg:    "foo is not a valid value for text. Expected format is quoted text.",
						Position:    1,
					},
				},
				SubCommands: []SubCommand{
					SubCommand{
						Name:        "reverse",
						Description: "reverses back what you type.",
						Arguments: []Argument{
							Argument{
								Name:        "text",
								ArgType:     "quoted text",
								Description: "text you want to print",
								ErrorMsg:    "foo is not a valid value for text. Expected format is quoted text.",
								Position:    0,
							},
						},
					},
					SubCommand{
						Name:        "quote",
						Description: "helps you stand on the shoulders of giants by using words from histories most articulate people",
						Arguments:   []Argument(nil),
						SubCommands: []SubCommand{
							SubCommand{
								Name:        "random",
								Description: "print a random quote from the a random author",
								Arguments:   []Argument(nil),
								SubCommands: []SubCommand(nil),
							},
							SubCommand{
								Name:        "author",
								Description: "prints a quote from the specified author",
								Arguments: []Argument{
									Argument{
										Name:        "authorName",
										ArgType:     "text",
										Description: "Full Name of an author",
										ErrorMsg:    "Please provide a valid author name, try someone famous",
										Position:    0,
									},
								},
								SubCommands: []SubCommand(nil),
							},
						},
					},
				},
			},
		},
		{
			testName:   "quoted text paramater value test",
			configPath: "./examples/helloWorld/simple.yaml",
			want: SlashCommand{
				Name:        "Print",
				Description: "Echos back what you type.",
				Arguments: []Argument{
					{
						Name:        "text",
						ArgType:     "quoted text",
						Description: "text you want to print",
						ErrorMsg:    "foo is not a valid value for text. Expected format is quoted text.",
						Position:    1,
					},
				},
				SubCommands: []SubCommand{
					SubCommand{
						Name:        "reverse",
						Description: "reverses back what you type.",
						Arguments: []Argument{
							Argument{
								Name:        "text",
								ArgType:     "quoted text",
								Description: "text you want to print",
								ErrorMsg:    "foo is not a valid value for text. Expected format is quoted text.",
								Position:    0,
							},
						},
					},
					SubCommand{
						Name:        "quote",
						Description: "helps you stand on the shoulders of giants by using words from histories most articulate people",
						Arguments:   []Argument(nil),
						SubCommands: []SubCommand{
							SubCommand{
								Name:        "random",
								Description: "print a random quote from the a random author",
								Arguments:   []Argument(nil),
								SubCommands: []SubCommand(nil),
							},
							SubCommand{
								Name:        "author",
								Description: "prints a quote from the specified author",
								Arguments: []Argument{
									Argument{
										Name:        "authorName",
										ArgType:     "text",
										Description: "Full Name of an author",
										ErrorMsg:    "Please provide a valid author name, try someone famous",
										Position:    0,
									},
								},
								SubCommands: []SubCommand(nil),
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			slashDef, _ := ioutil.ReadFile(test.configPath)

			newSlash, err := NewSlashCommand(slashDef)
			if err != nil {
				assert.Equal(t, test.expectedError, err)
			}
			assert.Equal(t, test.want, newSlash)
		})
	}
}

type getCommandStringTests struct {
	testName    string
	args        string
	want        string
	expectError bool
}

func TestGetCommandString(t *testing.T) {
	tests := []getCommandStringTests{
		{
			testName: "valid print example",
			args:     "/print",
			want:     "Print",
		},
		{
			testName:    "invalid print example",
			args:        "",
			want:        "",
			expectError: true,
		},
		{
			testName: "sub command example",
			args:     "/print reverse hsals",
			want:     "Print reverse",
		},
		{
			testName: "sub sub command example",
			args:     "/print quote random",
			want:     "Print quote random",
		},
		{
			testName: "sub sub command with value example",
			args:     "/print quote author Ben Franklin",
			want:     "Print quote author",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			testYamlPath := "./examples/helloWorld/simple.yaml"

			slashDef, _ := ioutil.ReadFile(testYamlPath)
			newSlash, _ := NewSlashCommand(slashDef)
			got, err := newSlash.getCommandString(test.args)
			if err != nil {
				assert.Equal(t, test.expectError, true)
			} else {
				assert.Equal(t, test.want, got)
			}
		})
	}
}

func TestGetPositionalArgs(t *testing.T) {
	got := GetPositionalArgs("foo \"man chu\"  \\choo wow")
	want := []string{"foo", "man chu", "\\choo", "wow"}
	assert.Equal(t, want, got)
}

func TestGetValues(t *testing.T) {

	commandAndArgs := "/print foo"
	newSlash, _ := NewSlashCommand(SimpleDef)
	got, _ := newSlash.getValues(commandAndArgs)

	want := map[string]string{"text": "foo"}
	assert.Equal(t, want, got)
}

func TestParse(t *testing.T) {
	slashCommandString := "/print foo"

	wantCommands := "Print"
	wantValues := map[string]string{"text": "foo"}

	newSlash, _ := NewSlashCommand(SimpleDef)
	gotCommands, gotValues, _ := newSlash.Parse(slashCommandString)

	assert.Equal(t, gotCommands, wantCommands)
	assert.Equal(t, gotValues, wantValues)
}

//TODO: move simple2 to a test data folder, create more test yaml some that should validate and some that shouldn't
const testDataDir = "./testData"

type validateSlashDefinitionTests struct {
	testName      string
	yamlName      string
	shouldBeValid bool
}

func TestValidateSlashDefinition(t *testing.T) {
	tests := []validateSlashDefinitionTests{
		{
			testName:      "test simple yaml file",
			yamlName:      "simple2.yaml",
			shouldBeValid: true,
		},
		{
			testName:      "test bad yaml file",
			yamlName:      "badDeffinition1.yaml",
			shouldBeValid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			s := SlashCommand{}

			yamldoc, _ := ioutil.ReadFile(testDataDir + "/" + test.yamlName)
			_ = yaml.Unmarshal([]byte(yamldoc), &s)

			got := validateSlashDefinition(&s)

			if test.shouldBeValid {
				assert.Nil(t, got)
			} else {
				assert.NotNil(t, got)
			}
		})
	}
}
