package slashparse

import (
	"bytes"
	"errors"
	"log"
	"regexp"
	"strings"
	"text/template"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

const (
	space       = ' '
	backspace   = '\\'
	doubleQuote = '"'
)

//Argument defines and argument in a slash command
type Argument struct {
	Name        string `yaml:"name" json:"name"`
	ArgType     string `yaml:"argtype" json:"argtype"`
	Description string `yaml:"description" json:"description"`
	ErrorMsg    string `yaml:"errorMsg" json:"errorMsg"`
	Position    int    `yaml:"position" json:"position"`
	Required    bool   `yaml:"required" json:"required"`
}

//SlashCommand defines the structure of a slash command string
type SlashCommand struct {
	Name        string       `yaml:"name" json:"name,omitempty"`
	Description string       `yaml:"description" json:"description"`
	Arguments   []Argument   `yaml:"arguments" json:"arguments,omitempty"`
	SubCommands []SubCommand `yaml:"subcommands" json:"subcommands"`
}

//SubCommand defines a command that proceded the slash command
type SubCommand struct {
	Name        string       `yaml:"name" json:"name"`
	Description string       `yaml:"description" json:"description"`
	Arguments   []Argument   `yaml:"arguments" json:"arguments"`
	SubCommands []SubCommand `yaml:"subcommands" json:"subcommands"`
}

//NewSlashCommand define a new slash command to parse
func NewSlashCommand(slashDef []byte) (s SlashCommand, err error) {
	unmarshalErr := yaml.Unmarshal([]byte(slashDef), &s)
	if unmarshalErr != nil {
		return s, unmarshalErr
	}

	validationErr := validateSlashDefinition(&s)
	if validationErr != nil {
		return s, validationErr
	}
	return s, nil
}

//GetSlashHelp returns a markdown formated help for a slash command
func (s *SlashCommand) GetSlashHelp() string {
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	helpTemplate, err := template.New("standardHelp.tpl").Funcs(funcMap).ParseFiles("./templates/standardHelp.tpl")
	if err != nil {
		log.Printf("Unable to load help template. %s", err.Error())
		return ""
	}

	var tpl bytes.Buffer
	if err := helpTemplate.Execute(&tpl, s); err != nil {
		log.Printf("Unable to execute help template. %s", err.Error())
		return ""
	}

	result := tpl.String()
	return result
}

//getValues takes a command and arguments and gets a dictionary of values by argument name
func (s *SlashCommand) getValues(CommandAndArgs string) (map[string]string, error) {
	m := make(map[string]string)

	//remove command from string
	command, err := s.getCommandString(CommandAndArgs)
	if err != nil {
		return m, err
	}

	//use regex for case insensitivity
	re := regexp.MustCompile(`(?i)/` + command)
	loc := re.FindStringIndex(CommandAndArgs)
	if len(loc) == 0 {
		return m, err //command not included in string?
	}

	args := strings.TrimSpace(CommandAndArgs[loc[1]:])

	if len(args) == 0 {
		return m, nil
	}

	// need to go ordered here?
	positionalArgs := GetPositionalArgs(args)

	for _, slashArg := range s.Arguments {
		position := slashArg.Position
		if len(positionalArgs) >= position {
			m[slashArg.Name] = positionalArgs[position-1]
		}
	}
	return m, nil
}

//getCommandString gets and validated the command portion of a command and argument string
func (s *SlashCommand) getCommandString(args string) (commandString string, err error) {
	argsSplit := strings.Fields(args)

	if len(argsSplit) < 1 {
		return "", err
	}

	command := strings.Replace(argsSplit[0], "/", "", 1)
	args = strings.Replace(args, "/", "", 1)

	//i hate this, regex might be better
	for _, subCommand := range s.SubCommands {

		for _, subSubCommand := range subCommand.SubCommands {
			subCommandString := s.Name + " " + subCommand.Name + " " + subSubCommand.Name
			if len(args) >= len(subCommandString) {
				if strings.EqualFold(args[:len(subCommandString)], subCommandString) {
					return subCommandString, nil
				}
			}

		}

		subCommandString := s.Name + " " + subCommand.Name
		if len(args) >= len(subCommandString) {
			if strings.EqualFold(args[:len(subCommandString)], subCommandString) {
				return subCommandString, nil
			}
		}
	}

	if strings.EqualFold(command, s.Name) {
		return s.Name, nil
	}

	return "", errors.New(command + " is not a valid command")
}

//Parse parse the command string
func (s *SlashCommand) Parse(slashString string) (string, map[string]string, error) {
	commandString, err := s.getCommandString(slashString)
	if err != nil {
		return "", nil, err
	}

	values, err := s.getValues(slashString)
	if err != nil {
		return "", nil, err
	}

	return commandString, values, nil
}

//GetPositionalArgs takes a string of arguments and splits it up by spaces and double quotes
func GetPositionalArgs(argString string) []string {
	var isQuoteText bool
	var previousCharacter rune
	args := make([]string, 0, 20)
	currentPosition := 0
	var currentArg string

	for _, character := range argString {
		switch character {
		case space:
			if len(currentArg) > 0 {
				if isQuoteText {
					currentArg += string(character)
				} else {
					// ignore duplicate spaces between
					if previousCharacter != space {
						args = append(args, currentArg)
						currentPosition++
						currentArg = ""
					}
				}
			}
		case doubleQuote:
			if isQuoteText {
				//this is and end quote
				isQuoteText = false
				args = append(args, currentArg)
				currentPosition++
				currentArg = ""
			} else {
				if previousCharacter != backspace {
					isQuoteText = true
				} else {
					//remove the escape character from the the value and add the quote
					currentArg = currentArg[:len(currentArg)-1] + string(doubleQuote)
				}
			}
		default:
			currentArg += string(character)
		}
		previousCharacter = character
	}

	if len(currentArg) > 0 {
		args = append(args, currentArg)
	}
	return args
}

func validateSlashDefinition(slashCommandDef *SlashCommand) (err error) {
	schemaLoader := gojsonschema.NewReferenceLoader(`file://C:/Users/eric/code/slashparse/schema.json`)

	documentLoader := gojsonschema.NewGoLoader(&slashCommandDef)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	if result.Valid() {
		return nil
	}
	log.Printf("The document is not valid. see errors :\n")
	for _, desc := range result.Errors() {
		log.Printf("- %s\n", desc)
	}
	return errors.New("Slash Command Deffinitaion is not valid")
}
