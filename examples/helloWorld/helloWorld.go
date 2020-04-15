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