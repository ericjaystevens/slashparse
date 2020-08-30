package slashparse

import (
	"errors"
	"io/ioutil"
	"strings"
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
var lotsOfArgsDef, _ = ioutil.ReadFile("./examples/helloWorld/lotsofargs.yaml")
var simpleDef2, _ = ioutil.ReadFile("./testData/simple2.yaml")
var doroDef, _ = ioutil.ReadFile("./testData/doro.yaml")
var wranglerDef, _ = ioutil.ReadFile("./testData/wrangler.yaml")

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
						Position:    0,
						ShortName:   "t",
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
								Required:    true,
								ShortName:   "t",
							},
						},
						commandPaths: []string{"Print reverse"},
					},
					SubCommand{
						Name:        "quote",
						Description: "helps you stand on the shoulders of giants by using words from histories most articulate people",
						Arguments:   []Argument(nil),
						SubCommands: []SubCommand{
							SubCommand{
								Name:         "random",
								Description:  "print a random quote from the a random author",
								Arguments:    []Argument(nil),
								SubCommands:  []SubCommand(nil),
								commandPaths: []string{"Print quote random"},
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
								SubCommands:  []SubCommand(nil),
								commandPaths: []string{"Print quote author"},
							},
						},
						commandPaths: []string{"Print quote"},
					},
					SubCommand{
						Name:         "help",
						Description:  "Display help.",
						commandPaths: []string{"Print help"},
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
						Position:    0,
						ShortName:   "t",
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
								Required:    true,
								ShortName:   "t",
							},
						},
						commandPaths: []string{"Print reverse"},
					},
					SubCommand{
						Name:        "quote",
						Description: "helps you stand on the shoulders of giants by using words from histories most articulate people",
						Arguments:   []Argument(nil),
						SubCommands: []SubCommand{
							SubCommand{
								Name:         "random",
								Description:  "print a random quote from the a random author",
								Arguments:    []Argument(nil),
								SubCommands:  []SubCommand(nil),
								commandPaths: []string{"Print quote random"},
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
								SubCommands:  []SubCommand(nil),
								commandPaths: []string{"Print quote author"},
							},
						},
						commandPaths: []string{"Print quote"},
					},
					SubCommand{
						Name:         "help",
						Description:  "Display help.",
						commandPaths: []string{"Print help"},
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
			assert.Equal(t, test.want.Name, newSlash.Name)
			assert.Equal(t, test.want.Arguments, newSlash.Arguments)
			assert.Equal(t, test.want.SubCommands[0], newSlash.SubCommands[0])
			assert.Equal(t, test.want.SubCommands[1], newSlash.SubCommands[1])
			assert.Equal(t, test.want.SubCommands[2].Name, newSlash.SubCommands[2].Name) //testing help with the handler set is tricky, so I tested around it.
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

type getValuesTests struct {
	testName       string
	commandAndArgs string
	slashDef       []byte
	want           map[string]string
}

func TestGetValues(t *testing.T) {
	tests := []getValuesTests{
		{
			testName:       "Slash command single arg",
			commandAndArgs: "/print foo",
			slashDef:       SimpleDef,
			want:           map[string]string{"text": "foo"},
		},
		{
			testName:       "Slash sub command single arg",
			commandAndArgs: "/print reverse foo",
			slashDef:       SimpleDef,
			want:           map[string]string{"text": "foo"},
		},
		{
			testName:       "Slash sub sub command single arg",
			commandAndArgs: "/print quote author foo",
			slashDef:       SimpleDef,
			want:           map[string]string{"authorName": "foo"},
		},
		{
			testName:       "slash command single named argument",
			commandAndArgs: "/print --text foo",
			slashDef:       SimpleDef,
			want:           map[string]string{"text": "foo"},
		},
		{
			testName:       "positional followed by named arguments",
			commandAndArgs: `/search "moon river" --search river --replace rising`,
			slashDef:       lotsOfArgsDef,
			want:           map[string]string{"text": "moon river", "search": "river", "replace": "rising"},
		},
		{
			testName:       "slash command single short named paramater",
			commandAndArgs: `/print -t foo`,
			slashDef:       SimpleDef,
			want:           map[string]string{"text": "foo"},
		},
		{
			testName:       "slash sub command single short named paramater",
			commandAndArgs: `/print reverse -t foo`,
			slashDef:       SimpleDef,
			want:           map[string]string{"text": "foo"},
		},
		{
			testName:       "mixed short and long named paramater",
			commandAndArgs: `/search -t "this land is your land" --search land -r hand`,
			slashDef:       lotsOfArgsDef,
			want:           map[string]string{"text": "this land is your land", "search": "land", "replace": "hand"},
		},
		{
			testName:       "wrangler list channels",
			commandAndArgs: `/wrangler list channels --channel-filter foo`,
			slashDef:       wranglerDef,
			want:           map[string]string{"channel-filter": "foo"},
		},
		{
			testName:       "wrangler list messages",
			commandAndArgs: `/wrangler list messages --count 10 --trim-length 11`,
			slashDef:       wranglerDef,
			want:           map[string]string{"count": "10", "trim-length": "11"},
		},
		{
			testName:       "wrangler list messages with defaults",
			commandAndArgs: `/wrangler list messages`,
			slashDef:       wranglerDef,
			want:           map[string]string{"count": "20", "trim-length": "50"},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			newSlash, err := NewSlashCommand(test.slashDef)
			if err != nil {
				assert.Fail(t, "slash command not processed")
			}
			got, _ := newSlash.getValues(test.commandAndArgs)

			assert.Equal(t, test.want, got)
		})
	}
}

type testParseTest struct {
	testName          string
	commandString     string
	slashDef          []byte
	wantCommandString string
	wantValues        map[string]string
	wantError         error
}

func TestParse(t *testing.T) {

	tests := []testParseTest{
		{
			testName:          "simple test",
			commandString:     "/print foo",
			slashDef:          SimpleDef,
			wantCommandString: "Print",
			wantValues:        map[string]string{"text": "foo"},
		},
		{
			testName:          "valid wrangler command",
			commandString:     "/wrangler list channels",
			slashDef:          wranglerDef,
			wantCommandString: "wrangler list channels",
			wantValues:        map[string]string{},
		},
		{
			testName:      "missing required subsubcommand",
			commandString: "/wrangler list",
			slashDef:      wranglerDef,
			wantError:     errors.New("/wrangler list requires an additional command. Try adding channels or messages. Please see /wrangler help for more info"),
		},
		{
			testName:      "missing required subcommand",
			commandString: "/wrangler",
			slashDef:      wranglerDef,
			wantError:     errors.New("/wrangler is not a valid command. Please see /wrangler help"),
		},
		{
			testName:      "invalid required subcommand",
			commandString: "/wrangler move invalid",
			slashDef:      wranglerDef,
			wantError:     errors.New("/wrangler move requires an additional command. Try adding thread. Please see /wrangler help for more info"),
		},
		{
			testName:          "switch on",
			commandString:     "/wrangler move thread 123 321 --show-root-message-in-summary",
			slashDef:          wranglerDef,
			wantCommandString: "wrangler move thread",
			wantValues:        map[string]string{"messageID": "123", "channelID": "321", "show-root-message-in-summary": "on"},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			newSlash, _ := NewSlashCommand(test.slashDef)
			gotCommands, gotValues, gotErr := newSlash.Parse(test.commandString)

			assert.Equal(t, test.wantCommandString, gotCommands)
			assert.Equal(t, test.wantValues, gotValues)
			if test.wantError != nil {
				assert.EqualError(t, gotErr, test.wantError.Error())
			}
		})
	}
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

func TestSetHandler(t *testing.T) {
	newSlash, _ := NewSlashCommand(SimpleDef)
	commandString, _, _ := newSlash.Parse("/print reverse pickle")
	myHandler := func(args map[string]string) (string, error) {
		return "reverseHandler called with text set as " + args["text"], nil
	}

	got := newSlash.SetHandler(commandString, myHandler)

	assert.Nil(t, got)
}

type invokeHandlerTests struct {
	testName      string
	commandString string
	slashDef      []byte
	want          string
	handler       func(map[string]string) (string, error)
}

func TestInvokeHandler(t *testing.T) {

	tests := []invokeHandlerTests{
		{
			testName:      "simple subCommand",
			commandString: "/print quote",
			slashDef:      SimpleDef,
			want:          "quoteHandler called",
			handler: func(args map[string]string) (string, error) {
				return "quoteHandler called", nil
			},
		},
		{
			testName:      "subCommand with argument",
			commandString: "/print reverse pickle",
			slashDef:      SimpleDef,
			want:          "reverseHandler called with text set as pickle",
			handler: func(args map[string]string) (string, error) {
				return "reverseHandler called with text set as " + args["text"], nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			newSlash, _ := NewSlashCommand(test.slashDef)
			commandString, values, _ := newSlash.Parse(test.commandString)

			_ = newSlash.SetHandler(commandString, test.handler)
			got, _ := newSlash.invokeHandler(commandString, values)

			assert.Equal(t, test.want, got)
		})
	}

	t.Run("invoke without setting handler", func(t *testing.T) {
		newSlash, _ := NewSlashCommand(SimpleDef)
		commandString, values, _ := newSlash.Parse("/print reverse pick")
		_, err := newSlash.invokeHandler(commandString, values)

		assert.EqualError(t, err, "No handler set")
	})

}

type executeTests struct {
	name          string
	commandString string
	want          string
	slashDef      []byte
}

func TestExecute(t *testing.T) {

	tests := []executeTests{
		{
			name:          "slashCommand Test",
			commandString: "/print echo",
			want:          "print called with argument echo",
			slashDef:      SimpleDef,
		},
		{
			name:          "subcommand Test",
			commandString: "/print reverse deep",
			want:          "reverseHandler called with text set as deep",
			slashDef:      SimpleDef,
		},
		{
			name:          "sub sub command Test",
			commandString: "/print quote author Shakespeare",
			want:          "quoteAuthorHandler called with authorName set as Shakespeare",
			slashDef:      SimpleDef,
		},
		{
			name:          "missing required argument",
			commandString: `/search "I once had a dream, it was a good dream to dream"`,
			want:          "required field search is missing, see /search help for more details",
			slashDef:      lotsOfArgsDef,
		},
		{
			name:          "multiple missing required arguments",
			commandString: `/search`,
			want:          "required fields text and search are missing, see /search help for more details",
			slashDef:      lotsOfArgsDef,
		},
		{
			name:          "missing required argument in sub command",
			commandString: `/print reverse`,
			want:          "required field text is missing, see /print help for more details",
			slashDef:      SimpleDef,
		},
		{
			name:          "missing 3 required args",
			commandString: `/print reverse`,
			want:          "required field text is missing, see /print help for more details",
			slashDef:      SimpleDef,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			newSlash, _ := NewSlashCommand(test.slashDef)

			newSlash.SetHandler("print quote author", func(args map[string]string) (string, error) {
				return "quoteAuthorHandler called with authorName set as " + args["authorName"], nil
			})

			newSlash.SetHandler("print reverse", func(args map[string]string) (string, error) {
				return "reverseHandler called with text set as " + args["text"], nil
			})

			newSlash.SetHandler("print", func(args map[string]string) (string, error) {
				return "print called with argument " + args["text"], nil
			})

			got, _ := newSlash.Execute(test.commandString)
			assert.Equal(t, test.want, got)

		})
	}

}

func TestGetSubCommand(t *testing.T) {
	newSlash, _ := NewSlashCommand(SimpleDef)
	commandString := "print quote random"
	got, _ := newSlash.getSubCommand(commandString)
	want := SubCommand{
		Name:         "random",
		Description:  "print a random quote from the a random author",
		Arguments:    []Argument(nil),
		SubCommands:  []SubCommand(nil),
		commandPaths: []string{"Print quote random"},
	}
	assert.Equal(t, want, got)

}

func TestGetSlashHelp(t *testing.T) {
	newSlash, _ := NewSlashCommand(SimpleDef)
	commandString := "print help"
	got, _ := newSlash.Execute(commandString)

	//just test the first line, to avoid so this doesn't have to be maintained while features are changeing so rapidly
	firstLine := strings.Split(got, "\n")[0]
	assert.Equal(t, firstLine, "#### /Print Help")
}

type argumentTypesTests struct {
	name          string
	commandString string
	want          string
	argName       string
	slashDef      []byte
}

func TestArgumentTypes(t *testing.T) {

	tests := []argumentTypesTests{
		{
			name:          "single arg remaining text",
			commandString: "/print reverse bunch of wisdom",
			argName:       "text",
			want:          "bunch of wisdom",
			slashDef:      simpleDef2,
		},
		{
			name:          "text arg then remaining text arg",
			commandString: "/doro start 45 getting after it",
			argName:       "log",
			want:          "getting after it",
			slashDef:      doroDef,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			newSlash, _ := NewSlashCommand(test.slashDef)
			got, _ := newSlash.getValues(test.commandString)

			assert.Equal(t, test.want, got[test.argName])
		})

	}
}
