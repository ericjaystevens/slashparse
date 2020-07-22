#### /{{.Name}} Help
-- *{{.Description}}*

`/{{ .Name | ToLower }}{{range $arg := .Arguments}} {{if not $arg.Required}}[{{end}}{{$arg.Name}}{{if not $arg.Required}}]{{end}}{{end}}`

#### Arguments
{{range $arg := .Arguments}}
* **{{$arg.Name}}**: {{if not $arg.Required}}(optional){{end}} _{{$arg.Description}}_
{{end}}