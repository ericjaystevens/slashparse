package main

import (
	"io/ioutil"
	"log"
	"strconv"
)

func main() {

	outfilename := `.\jsonSchema.go`
	infile := `.\schema.json`

	dat, err := ioutil.ReadFile(infile)
	if err != nil {
		log.Print(err.Error())
	}
	jsonSchemaText := strconv.Quote(string(dat))

	codeStart := `// THIS  FILE IS GENERATED, DO NOT EDIT, INSTEAD UPDATE schema.json and run generate/generateFromStatic.go
	
package slashparse

const jsonSchemaContent =`

	genCode := codeStart + jsonSchemaText
	err = ioutil.WriteFile(outfilename, []byte(genCode), 0644)
	if err != nil {
		log.Print(err.Error())
	}

	outfile := "./templates.go"
	helpTemplateFile := "templates/standardHelp.tpl"

	dat, err = ioutil.ReadFile(helpTemplateFile)
	if err != nil {
		log.Print(err.Error())
	}
	helpTemplateText := strconv.Quote(string(dat))

	genCode = `// THIS  FILE IS GENERATED, DO NOT EDIT, INSTEAD UPDATE templates/standardHelp.tpl and run generate/generateFromStatic.go

package slashparse

const helpTemplateContent =` + helpTemplateText

	err = ioutil.WriteFile(outfile, []byte(genCode), 0644)
	if err != nil {
		log.Print(err.Error())
	}
}
