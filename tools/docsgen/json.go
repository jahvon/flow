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
	schemaDir = "docs/schemas"
	idBase    = "https://flowexec.io/schemas"
)

//nolint:gocognit
func generateJSONSchemas() {
	sm := schema.RegisteredSchemaMap()
	for fn, s := range sm {
		if slices.Contains(TopLevelPages, fn.Title()) {
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
	s.ID, _ = url.JoinPath(idBase, file.JSONSchemaFile())
}
