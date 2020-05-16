package slashparse

import (
	"errors"
	"strings"

	"io/ioutil"

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
	Description string `yaml:"description"`
}

type SlashCommand struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Arguments   []Argument `yaml:"arguments"`
	Values map[string]string
}

type Slash interface {
	GetSlashHelp() string
}

//NewSlashCommand define a new slash command to parse
func NewSlashCommand(args []string, pathToYaml string) (s SlashCommand, err error) {

	slashDef, yamlErr := ioutil.ReadFile(pathToYaml)
	if yamlErr != nil {
		return s, yamlErr
	}

	unmarshalErr := yaml.Unmarshal([]byte(slashDef), &s)
	if unmarshalErr != nil {
		return s, unmarshalErr
	}

	_, commandErr := s.GetCommandString(args)
	if commandErr != nil {
		return SlashCommand{}, commandErr
	}

	s.Values = s.GetValues(args)

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

func (s *SlashCommand) GetValues (args []string) (map[string]string) {

m := make(map[string]string)
	m["text"] = "foo bar"
	return m
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
