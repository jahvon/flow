package main

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"slices"

	"github.com/jahvon/flow/tools/docsgen/schema"
)

const (
	schemaDir = "schemas"
	idBase    = "https://raw.githubusercontent.com/jahvon/flow/HEAD"
)

//nolint:gocognit
func generateJSONSchemas() {
	sm := schema.RegisteredSchemaMap()
	for fn, s := range sm {
		if slices.Contains(TopLevelPages, fn.Title()) { //nolint:nestif
			updateFileID(s, fn)
			for key, value := range s.Properties {
				if !value.Ext.IsExported() {
					delete(s.Properties, key)
					continue
				}
				schema.MergeSchemas(s, value, fn, sm)
			}
			for _, value := range s.Definitions {
				schema.MergeSchemas(s, value, fn, sm)
			}

			schemaJSON, err := json.MarshalIndent(s, "", "  ")
			if err != nil {
				panic(err)
			}
			filePath := filepath.Clean(filepath.Join(rootDir(), schemaDir, fn.JSONSchemaFile()))
			file, err := os.Create(filePath)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			if _, err := file.WriteString(string(schemaJSON)); err != nil {
				panic(err)
			}
		}
	}
}

func updateFileID(s *schema.JSONSchema, file schema.FileName) {
	s.ID, _ = url.JoinPath(idBase, schemaDir, file.JSONSchemaFile())
}
