package slashparse

import (
	"log"

	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Slashdef struct {
	name string
	description string
	B struct {
		RenamedC int   `yaml:"c"`
		D        []int `yaml:",flow"`
	}
}

type SlashCommand struct {
	name string
	description string
}

type Slash interface {
	GetSlashHelp() string
}

//NewSlashCommand define a new slash command to parse
func NewSlashCommand(args []string, pathToYaml string) SlashCommand {

	m := make(map[interface{}]interface{})

	slashDef, yamlerr := ioutil.ReadFile(pathToYaml)
	if yamlerr != nil {
		return SlashCommand{}
	}

	err := yaml.Unmarshal([]byte(slashDef), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	val, _ := m["name"].(string)
	desc, _ := m["description"].(string)

	slashCommand := SlashCommand{
		name: val,
		description: desc,
	}
	return  slashCommand
}

func (s *SlashCommand) GetSlashHelp() string {
	
	header := "## " + s.name + " Help"
	
	description := "* " + s.description + " *"

	arguments := "### Arguments"

	return header + "\n" + description + "\n\n" + arguments + "\n"
}