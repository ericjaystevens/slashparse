#### /{{.Name}} Help
-- *{{.Description}}*

`/{{ .Name | ToLower }}{{range $arg := .Arguments}} {{if not $arg.Required}}[{{end}}{{$arg.Name}}{{if not $arg.Required}}]{{end}}{{end}}`

#### Arguments
{{range $arg := .Arguments}}
* **{{$arg.Name}}**: {{if not $arg.Required}}(optional){{end}} _{{$arg.Description}}_
{{end}}

#### Available Commands
{{range $subCommand := .SubCommands }}
* **{{$subCommand.Name}}**: __{{$subCommand.Description}}__
  `/{{$.Name | ToLower}} {{$subCommand.Name}} {{range $arg := $subCommand.Arguments}}{{if not $arg.Required}}[{{end}}{{$arg.Name}}{{if not $arg.Required}}]{{end}}{{end}}
  {{range $subSubCommand := .SubCommands}}  *  **{{$subSubCommand.Name}}**: {{$subSubCommand.Description}}
    `/{{$.Name}} {{$subCommand.Name}} {{$subSubCommand.Name}} {{range $arg := $subSubCommand.Arguments}}{{if not $arg.Required}}[{{end}}{{$arg.Name}}{{if not $arg.Required}}[{{end}}{{end}}`{{end}}
{{end}}