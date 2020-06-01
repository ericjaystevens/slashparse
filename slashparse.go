package slashparse

import (
	"errors"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type Slashdef struct {
	name        string
	description string
	arguments   struct {
		name        string
		description string
	}
}

type Argument struct {
	Name        string `yaml:"name"`
	ArgType     string `yaml:argtype`
	Description string `yaml:"description"`
	ErrorMsg    string `yaml:"errorMsg"`
}

type SlashCommand struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Arguments   []Argument `yaml:"arguments"`
	Values      map[string]string
}

type Slash interface {
	GetSlashHelp() string
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

func (s *SlashCommand) GetValues(args string) (map[string]string, error) {
	m := make(map[string]string)

	//remove command from string
	command, err := s.GetCommandString(args)
	if err != nil {
		return m, err
	}

	//use regex for case insesitvity
	re := regexp.MustCompile(`(?i)/` + command)
	loc := re.FindStringIndex(args)
	if len(loc) == 0 {
		return m, err //command not included in string?
	}

	parameters := strings.TrimSpace(args[loc[1]:])

	if len(parameters) == 0 {
		return m, nil
	}

	// need to go ordered here?
	positionalParams := GetPositionalParams(parameters)

	for _, slashArg := range s.Arguments {
		if slashArg.ArgType == "quoted text" {
			// remove quote and start building string
			//m["text"] = args[position][1:len(args[position])]
			m["text"] = positionalParams[1]
		}
	}
	return m, nil
}

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

func GetPositionalParams(paramString string) (splitParams []string) {
	var isQuoteText bool
	var previousCharacter rune
	params := make([]string, 0, 20)
	curPosition := 0
	var currentParam string

	for _, character := range paramString {

		switch character {
		case ' ':
			if len(currentParam) > 0 {
				if isQuoteText {
					currentParam += string(character)
				} else {
					params = append(params, currentParam)
					curPosition++
					currentParam = ""
				}
			}
		case '"':
			if isQuoteText {
				//this is and end quote
				//params[curPosition] += string(character)
				isQuoteText = false
				params = append(params, currentParam)
				curPosition++
				currentParam = ""
			} else {
				if previousCharacter != '\\' {
					isQuoteText = true
					//params[curPosition] += string(character)
				} else {
					//remove the escape character from the the value and add the quote
					//params[curPosition][len(params[curPosition])-1:] = string(character)
				}

			}
			previousCharacter = character
			break
		default:
			currentParam += string(character)
		}
		previousCharacter = character

	}
	if len(currentParam) > 0 {
		params = append(params, currentParam)
	}
	splitParams = params
	return
}
