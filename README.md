[![Go Report Card](https://goreportcard.com/badge/github.com/ericjaystevens/slashparse)](https://goreportcard.com/report/github.com/ericjaystevens/slashparse)

# slashparse
Go module for parsing slash commands

This is in the proof of concept stages of development, I'm not exactly sure how it will all turn out or if it has a lot of value.

## Design Goals

1. Support many Go implementation of Slash commands,  
1. Provide useful tooling to generate help docs and autocompletion in popular formats
1. Be opinionated about standardization and conventions in Slash commands.
1. Use text like YAML or JSON to define a command in an effort to make definitions multi-platform and easy to translate to other written languages. 


### How to use

#### Define your slash command in Yaml

```
---
name: Print
description: Echos back what you type.
arguments:
  - name: text
    argtype: quoted text
    description: text you want to print
    errorMsg: foo is not a valid value for text. Expected format is quoted text.
    position: 1
subcommands:
  - name: quote
    description: helps you stand on the shoulders of giants by using words from histories most articulate people
    subcommands:
      - name: random
        description: print a random quote from the a random author
      - name: author
        description: prints a quote from the specified author
        arguments:
          - name: authorName
            argtype: text
            description: Full Name of an author
            errorMsg: Please provide a valid author name, try someone famous "
```

#### Use slashParse to parse your arguments

```
package main

import com.gitlab.ericjaystevens.slashparse

const pathToYaml = "path/to/yaml"

ExecuteCommand(args string) (responce string, Error) {

  //define the slash command
  slashDef, _ := ioutil.ReadFile(pathToYaml)
  slashCommand, _ = slashparse.NewSlashCommand(slashDef)
	
  //parse it to get a command string and a map with all your arguments and their values
  //This should provie a helpfull error if a require argument is missing or the command is not valid
  command, values, err := p.slashCommand.Parse(args.Command)
	if err != nil {
		text := "bad command see help"
		return text, err
	}


	switch command {
	case "Print":
		return executePrint(values["text"])
	case "help":
		markdownHelp := p.slashCommand.GetSlashHelp()
		return markdownHelp, nil
  case "quote random":
    return executeQuoteRandom()
  case "quote author":
    return executeQuoteByAuthor(values["authorName"])
	default:
		text := "Unknown unknown"
	}

}

func print(input string) string{
	return input
}
```

### What your users will see

#### argument parsing

```
> /print "Hello World!"
```

will provide the expected output

```
Hello World!
```

Alternatively if they could use the flags you set instead of position

#### Help

If your user can access generated help by running your slash command and help. 

```
/print help
```

Slash parse give nice help output. 

```
Print: "Prints the things you want"

/print "text"

  Text: -t,--tex
    text you want to display
```


#### Invalid commands

This examples requires quotes for the string so if a user runs

```
/print foo
```

They will receive the error message you defined.

```
Invalid Command. Please run /print help for more information
```
