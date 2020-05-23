[![Go Report Card](https://goreportcard.com/badge/github.com/ericjaystevens/slashparse)](https://goreportcard.com/report/github.com/ericjaystevens/slashparse)

# slashparse
Go module for parsing slash commands

This is in Beta and updates will break things, I'm not exactly sure how it will all turn out.


### How to use

#### Define your slash command in Yaml

```
SlashCommand: print 
  Description: "Prints the things you want"
  Arguments: [{
    name: "text",
    type: "quotedText",
    description: "text you want to display",
    Position: 1,
    Required: True,
    shortSwitch: "-t",
    longSwitch: "--text" 
  }]
  ParseErrorMessage: "Invalid Command. Please run ```/print help``` for more information"
```

#### Use slashParse to parse your arguments

```
package main

import com.gitlab.ericjaystevens.slashparse

ExecuteCommand(args *model.CommandArgs) (*model.CommandResponse, Error) {

	printCmd, err := slashparse.NewSlashCommand(PrintCommand.yaml, args)
        if err != nil{
		return err, nil
	}	

	commands = printCmd.GetCommandPath() // returns @("print")

	case commands{
	
	switch printCmd.path:
		return print(input), nil
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
