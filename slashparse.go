package slashparse

import (
	"errors"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	space       = ' '
	backspace   = '\\'
	doubleQuote = '"'
)

//Argument defines and argument in a slash command
type Argument struct {
	Name        string `yaml:"name"`
	ArgType     string `yaml:"argtype"`
	Description string `yaml:"description"`
	ErrorMsg    string `yaml:"errorMsg"`
	Position    int    `yaml:"position"`
}

//SlashCommand defines the structure of a slash command string
type SlashCommand struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Arguments   []Argument `yaml:"arguments"`
	Values      map[string]string
}

//NewSlashCommand define a new slash command to parse
func NewSlashCommand(args string, slashDef []byte) (s SlashCommand, err error) {
	unmarshalErr := yaml.Unmarshal([]byte(slashDef), &s)
	if unmarshalErr != nil {
		return s, unmarshalErr
	}

	_, commandErr := s.GetCommandString(args)
	if commandErr != nil {
		return SlashCommand{}, commandErr
	}

	var argErr error
	s.Values, argErr = s.GetValues(args)
	if argErr != nil {
		return SlashCommand{}, argErr
	}

	return s, nil
}

//GetSlashHelp returns a markdown formated help for a slash command
func (s *SlashCommand) GetSlashHelp() string {
	header := "## " + s.Name + " Help"

	description := "* " + s.Description + " *"

	arguments := "### Arguments"

	//for each argument in arguments print name.
	for _, argument := range s.Arguments {
		arguments += "\n\n* " + argument.Name + ": " + argument.Description
	}
	return header + "\n" + description + "\n\n" + arguments + "\n"
}

//GetValues takes a command and arguments and gets a dictionary of values by argument name
func (s *SlashCommand) GetValues(CommandAndArgs string) (map[string]string, error) {
	m := make(map[string]string)

	//remove command from string
	command, err := s.GetCommandString(CommandAndArgs)
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

//GetCommandString gets and validated the command portion of a command and argument string
func (s *SlashCommand) GetCommandString(args string) (commandString string, err error) {
	argsSplit := strings.Fields(args)

	if len(argsSplit) < 1 {
		return "", err
	}

	command := strings.Replace(argsSplit[0], "/", "", 1)

	if strings.EqualFold(command, s.Name) {
		return s.Name, nil
	}

	return "", errors.New(command + " is not a valid command")
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
