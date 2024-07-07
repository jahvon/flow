package main

import (
	"encoding/json"
	"os"

	"github.com/xeipuuv/gojsonschema"

	"github.com/jahvon/flow/tools/builder"
)

func GenerateExecutableDefinitionSchema() error {
	schemaLoader := gojsonschema.NewGoLoader(builder.ExamplesExecExecutableDefinition())
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return err
	}

	file, err := os.Create("executables_schema.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(schema)
}
