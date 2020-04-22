package slashparse

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
)

type Slashdef struct {
	A string
	B struct {
		RenamedC int   `yaml:"c"`
		D        []int `yaml:",flow"`
	}
}

var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`

//NewSlashCommand define a new slash command to parse
func NewSlashCommand(args []string, pathToYaml string) string {

	m := make(map[interface{}]interface{})

	err := yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	slashCmd := args[0]
	val := m["a"]
	return fmt.Sprintf("works on %s so %v", slashCmd, val)
}
