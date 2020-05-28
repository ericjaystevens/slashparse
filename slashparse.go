package slashparse

import (
	"errors"
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
func NewSlashCommand(args []string, slashDef []byte) (s SlashCommand, err error) {

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

func (s *SlashCommand) GetValues(args []string) (map[string]string, error) {
	m := make(map[string]string)

	var position int
	for _, slashArg := range s.Arguments {
		if slashArg.ArgType == "quoted text" {
			position = 1 //position indicates  where the value of the argument starts
			// maybee we should turn the args array into a string and no deal with the comlexity of presplit args?
			if len(args) > position {

				//is the first character a quote
				if args[position][0] == '"' {
					// remove quote and start building string
					m["text"] = args[position][1:len(args[position])]

					if args[position][len(args[position])-1:] != `"` {
						position++
						//append strings until end quote is found
						for endQuoteFound := true; endQuoteFound; endQuoteFound = !(args[position][len(args[position])-1:] == `"`) {
							m["text"] += " "
							m["text"] += args[position]
						}
					}
					//remove the ending quote
					m["text"] = m["text"][0 : len(m["text"])-1]
				} else {
					return m, errors.New(slashArg.ErrorMsg)
				}
			}
		}
	}
	return m, nil
}

func (s *SlashCommand) GetCommandString(args []string) (commandString string, err error) {
	if len(args) < 0 {
		return "", err
	}

	command := strings.Replace(args[0], "/", "", 1)

	if strings.EqualFold(command, s.Name) {
		return s.Name, nil
	}

	return "", errors.New(command + " is not a valid command")
}
