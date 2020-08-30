// Package slashparse is a parser for slash commands.
package slashparse

import (
	"bytes"
	"errors"
	"fmt"
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
	Default     string `yaml:"default" json:"default"`
	Description string `yaml:"description" json:"description"`
	ErrorMsg    string `yaml:"errorMsg" json:"errorMsg"`
	Position    int    `yaml:"position" json:"position"`
	Required    bool   `yaml:"required" json:"required"`
	ShortName   string `yaml:"shortName" json:"shortName"`
}

//SlashCommand defines the structure of a slash command string
type SlashCommand struct {
	Name               string       `yaml:"name" json:"name,omitempty"`
	Description        string       `yaml:"description" json:"description"`
	Arguments          []Argument   `yaml:"arguments" json:"arguments,omitempty"`
	SubCommands        []SubCommand `yaml:"subcommands" json:"subcommands,omitempty"`
	handler            func(map[string]string) (string, error)
	SubCommandRequired bool `yaml:"subCommandRequired" json:"subCommandRequired"`
}

//SubCommand defines a command that proceeded the slash command
type SubCommand struct {
	Name               string       `yaml:"name" json:"name"`
	Description        string       `yaml:"description" json:"description"`
	Arguments          []Argument   `yaml:"arguments" json:"arguments"`
	SubCommands        []SubCommand `yaml:"subcommands" json:"subcommands"`
	commandPaths       []string
	handler            func(map[string]string) (string, error)
	SubCommandRequired bool `yaml:"subCommandRequired" json:"subCommandRequired"`
}

//implimented by SlashCommand and SubCommand
type command interface {
	getArgsValues() (map[string]string, error)
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

	//range makes a copy so changes are not persistant, so use iterators instead
	for subCommandPosition, subCommand := range s.SubCommands {
		subCommandPath := s.Name + " " + subCommand.Name
		s.SubCommands[subCommandPosition].commandPaths = append(subCommand.commandPaths, subCommandPath)
		for subSubCommandPostion, subSubCommand := range subCommand.SubCommands {
			subSubCommandPath := subCommandPath + " " + subSubCommand.Name
			s.SubCommands[subCommandPosition].SubCommands[subSubCommandPostion].commandPaths = append(subSubCommand.commandPaths, subSubCommandPath)
		}
	}

	//Add built-in help subcommand

	helpSubcommand := SubCommand{
		Name:         "help",
		Description:  "Display help.",
		commandPaths: []string{s.Name + " help"},
		handler:      func(args map[string]string) (string, error) { return s.GetSlashHelp(), nil },
	}
	s.SubCommands = append(s.SubCommands, helpSubcommand)
	return s, nil
}

// getCommandPath gets the full command path that should call a command or sub command
// hardcoded for now
func (s *SubCommand) getCommandPath() string {

	//this will need to be changed to a matching method when muliple command paths are supported (aliases or reversable subcommands)
	return s.commandPaths[0]
}

// SetHandler sets the function that should be called based on the set of slash command and subcommands
func (s *SlashCommand) SetHandler(commandString string, handler func(map[string]string) (string, error)) error {

	if strings.EqualFold(commandString, s.Name) {
		s.handler = handler
	}

	for i, subCommand := range s.SubCommands {
		commandPath := subCommand.getCommandPath()

		if strings.EqualFold(commandString, commandPath) {
			s.SubCommands[i].handler = handler
		}

		for subSubCommandPostion, subSubCommand := range subCommand.SubCommands {
			subSubcommandPath := subSubCommand.getCommandPath()
			if strings.EqualFold(commandString, subSubcommandPath) {
				s.SubCommands[i].SubCommands[subSubCommandPostion].handler = handler
			}
		}
	}
	return nil
}

func (s *SlashCommand) invokeHandler(commandString string, args map[string]string) (string, error) {
	if strings.EqualFold(commandString, s.Name) {
		if s.handler != nil {
			return s.handler(args)
		}
		return "", errors.New("No handler set")
	}

	subCommand, err := s.getSubCommand(commandString)
	if err != nil {
		return "", err
	}

	if subCommand.handler != nil {
		return subCommand.handler(args)
	}
	return "", errors.New("No handler set")
}

//GetSlashHelp returns a markdown formated help for a slash command
func (s *SlashCommand) GetSlashHelp() string {
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	//	helpTemplate, err := template.New("standardHelp.tpl").Funcs(funcMap).ParseFiles("./templates/standardHelp.tpl")
	helpTemplate, err := template.New("standardHelp.tpl").Funcs(funcMap).Parse(helpTemplateContent)
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

	if strings.EqualFold(command, s.Name) {
		return s.getArgsValues(command, CommandAndArgs[loc[1]:], s.Arguments, s.Name)
	}

	subCommand, err := s.getSubCommand(command)
	if err != nil {
		return m, err
	}

	return s.getArgsValues(command, CommandAndArgs[loc[1]:], subCommand.Arguments, s.Name)

}

//getNamedArgValues gets named arguments and switches into a map of strings
func (s *SlashCommand) getNamedArgValues(commandString, argString string) (m map[string]string) {
	m = make(map[string]string)

	splitArgs := GetPositionalArgs(argString)
	var argumentName string
	for _, splitArg := range splitArgs {
		if argumentName != "" {
			m[argumentName] = splitArg
			argumentName = ""
		}
		if strings.HasPrefix(splitArg, "--") {
			isSwitch, arg := s.findMatchingSwitch(commandString, splitArg[2:])
			if isSwitch {
				m[arg.Name] = "on"
			} else {
				argumentName = splitArg[2:]
			}
		} else if strings.HasPrefix(splitArg, "-") {
			argument, _ := s.getArgumentFromShortName(commandString, splitArg[1:])
			argumentName = argument.Name
		}
	}

	return m
}

//pass in potential switch without -- or -
func (s *SlashCommand) findMatchingSwitch(commandString string, potentialSwitch string) (bool, Argument) {

	if strings.EqualFold(commandString, s.Name) {
		for _, arg := range s.Arguments {
			if arg.ArgType == "switch" {
				if strings.EqualFold(arg.Name, potentialSwitch) {
					return true, arg
				}
			}
		}
	}

	subCommand, _ := s.getSubCommand(commandString)
	for _, arg := range subCommand.Arguments {
		if arg.ArgType == "switch" {
			if strings.EqualFold(arg.Name, potentialSwitch) {
				return true, arg
			}
		}
	}

	return false, Argument{}

}

// getArgumentFromShortName returns an argument that matches the shortname
func (s *SlashCommand) getArgumentFromShortName(commandString string, shortName string) (argument Argument, err error) {
	if strings.EqualFold(commandString, s.Name) {
		for _, arg := range s.Arguments {
			if arg.ShortName == shortName {
				return arg, nil
			}
		}
		return argument, fmt.Errorf("Unknown paramater '%s', see /%s help for more details", shortName, commandString)
	}

	subCommand, _ := s.getSubCommand(commandString)

	for _, arg := range subCommand.Arguments {
		if arg.ShortName == shortName {
			return arg, nil
		}
	}

	return argument, fmt.Errorf("Unknown paramater '%s', see /%s help for more details", shortName, commandString)
}

func (s *SlashCommand) getArgsValues(commandString string, argString string, commandArgs []Argument, slashCommandName string) (m map[string]string, err error) {

	m = make(map[string]string)
	missingArgs := make([]string, 0, 8)
	splitArgs := GetPositionalArgs(argString)

	for _, commandArg := range commandArgs {
		if commandArg.Default != "" {
			m[commandArg.Name] = commandArg.Default
		}

		position := commandArg.Position
		if len(splitArgs) > position {
			if strings.HasPrefix(splitArgs[position], "-") {
				break
			}
			switch commandArg.ArgType {
			case "text", "quoted text", "number":
				m[commandArg.Name] = splitArgs[position]
			case "remaining text":
				m[commandArg.Name] = strings.Join(splitArgs[position:], " ")
			}
		} else {
			if commandArg.Required {
				missingArgs = append(missingArgs, commandArg.Name)
			}
		}
	}

	namedMap := s.getNamedArgValues(commandString, argString)

	for k, v := range namedMap {
		m[k] = v
		for i, missingArg := range missingArgs {
			if missingArg == v {
				missingArgs[i] = missingArgs[len(missingArgs)-1]
				missingArgs[len(missingArgs)-1] = ""
				missingArg = missingArg[:len(missingArg)-1]
			}
		}
	}
	if len(missingArgs) > 0 {
		return m, getMissingArgError(missingArgs, slashCommandName)
	}

	return m, nil
}

func getMissingArgError(missingArgs []string, commandName string) error {

	commandName = strings.ToLower(commandName)
	if len(missingArgs) > 2 {
		return fmt.Errorf("required fields %s, and %s are missing, see /%s help for more details", strings.Join(missingArgs[:len(missingArgs)-1], ", "), missingArgs[len(missingArgs)-1], commandName)
	}
	if len(missingArgs) > 1 {
		return fmt.Errorf("required fields %s and %s are missing, see /%s help for more details", missingArgs[0], missingArgs[1], commandName)

	}
	return fmt.Errorf("required field %v is missing, see /%s help for more details", missingArgs[0], commandName)

}

//getCommandString gets and validated the command portion of a command and argument string
func (s *SlashCommand) getCommandString(args string) (commandString string, err error) {
	argsSplit := strings.Fields(args)

	if len(argsSplit) < 1 {
		return "", err
	}

	command := strings.Replace(argsSplit[0], "/", "", 1)
	args = strings.Replace(args, "/", "", 1)

	//check sub subcommand
	for _, subCommand := range s.SubCommands {
		for _, subSubCommand := range subCommand.SubCommands {
			subCommandString := s.Name + " " + subCommand.Name + " " + subSubCommand.Name
			if len(args) >= len(subCommandString) {
				if strings.EqualFold(args[:len(subCommandString)], subCommandString) {
					return subCommandString, nil
				}
			}
		}

		//check each subcommand
		subCommandString := s.Name + " " + subCommand.Name
		if len(args) >= len(subCommandString) {
			if strings.EqualFold(args[:len(subCommandString)], subCommandString) {
				if subCommand.SubCommandRequired {
					var requiredSubCommands []string

					for _, requiredSubCommand := range subCommand.SubCommands {
						requiredSubCommands = append(requiredSubCommands, requiredSubCommand.Name)
					}

					if len(requiredSubCommands) > 2 {
						return "", fmt.Errorf("/%s requires an additional command. Try adding %s, or %s. Please see /%s help for more info", subCommandString, strings.Join(requiredSubCommands[:len(requiredSubCommands)-1], ", "), requiredSubCommands[len(requiredSubCommands)-1], s.Name)
					}
					if len(requiredSubCommands) > 1 {
						return "", fmt.Errorf("/%s requires an additional command. Try adding %s or %s. Please see /%s help for more info", subCommandString, requiredSubCommands[0], requiredSubCommands[1], s.Name)
					}
					return "", fmt.Errorf("/%s requires an additional command. Try adding %s. Please see /%s help for more info", subCommandString, requiredSubCommands[0], s.Name)
				}
				return subCommandString, nil
			}
		}
	}

	if strings.EqualFold(command, s.Name) {
		if s.SubCommandRequired {
			return "", fmt.Errorf("/%s is not a valid command. Please see /%s help", s.Name, s.Name)
		}
		return s.Name, nil
	}

	return "", fmt.Errorf("/%s is not a valid command. Please see /%s help", command, s.Name)
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

//Execute parses and runs the configured handler to process your command.
func (s *SlashCommand) Execute(slashString string) (string, error) {
	commandString, values, err := s.Parse(slashString)
	if err != nil {
		return err.Error(), err
	}

	msg, err := s.invokeHandler(commandString, values)
	return msg, err
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
	schemaLoader := gojsonschema.NewBytesLoader([]byte(jsonSchemaContent))
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
	return errors.New("Slash Command Definition is not valid")
}

func (s *SlashCommand) getSubCommand(commandString string) (SubCommand, error) {

	for _, subCommand := range s.SubCommands {

		for _, path := range subCommand.commandPaths {
			if strings.EqualFold(commandString, path) {
				return subCommand, nil
			}
		}

		for _, subSubCommand := range subCommand.SubCommands {
			for _, path := range subSubCommand.commandPaths {
				if strings.EqualFold(commandString, path) {
					return subSubCommand, nil
				}
			}
		}
	}
	return SubCommand{}, errors.New("Unable to find mathing subcommand")
}
